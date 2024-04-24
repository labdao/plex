package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/internal/ipfs"
	"gorm.io/gorm"
)

var autotaskWebhook = os.Getenv("AUTOTASK_WEBHOOK")

var rateLimiter = NewTokenBucketRateLimiter(1, 1)

type TokenBucketRateLimiter struct {
	tokenBucket chan struct{}
	fillRate    time.Duration
}

func NewTokenBucketRateLimiter(fillRate time.Duration, capacity int) *TokenBucketRateLimiter {
	limiter := &TokenBucketRateLimiter{
		tokenBucket: make(chan struct{}, capacity),
		fillRate:    fillRate,
	}
	go limiter.fillTokenBucket()
	return limiter
}

func (l *TokenBucketRateLimiter) Wait() {
	<-l.tokenBucket
}

func (l *TokenBucketRateLimiter) fillTokenBucket() {
	ticker := time.NewTicker(l.fillRate)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case l.tokenBucket <- struct{}{}:
		default:
		}
	}
}

type postData struct {
	RecipientAddress string `json:"recipientAddress"`
	Cid              string `json:"cid"`
}

type Response struct {
	Status string `json:"status"`
}

func BuildTokenMetadata(db *gorm.DB, flow *models.Flow) (string, error) {
	var jobs []models.Job
	if err := db.Where("flow_id = ?", flow.ID).Find(&jobs).Error; err != nil {
		return "", fmt.Errorf("failed to retrieve jobs: %v", err)
	}

	metadata := map[string]interface{}{
		"name":        flow.Name,
		"description": "Research, Reimagined. All Scientists Welcome.",
		"image":       "ipfs://QmQZLrUPxh4WMmzpQGhUYRsMwU2BXfmFa3YAFhFKkRgHTZ", // Default image is glitchy LabDAO logo gif
		"flow":        []map[string]interface{}{},
	}

	for _, job := range jobs {
		var tool models.Tool
		if err := db.Where("cid = ?", job.ToolID).First(&tool).Error; err != nil {
			return "", fmt.Errorf("failed to retrieve tool: %v", err)
		}

		var inputFiles []models.DataFile
		if err := db.Model(&job).Association("InputFiles").Find(&inputFiles); err != nil {
			return "", fmt.Errorf("failed to retrieve input files: %v", err)
		}

		var outputFiles []models.DataFile
		if err := db.Model(&job).Association("OutputFiles").Find(&outputFiles); err != nil {
			return "", fmt.Errorf("failed to retrieve output files: %v", err)
		}

		ioObject := map[string]interface{}{
			"tool":    map[string]interface{}{},
			"inputs":  []map[string]interface{}{},
			"outputs": []map[string]interface{}{},
			"state":   job.State,
			"errMsg":  job.Error,
		}

		log.Printf("Pinning tool JSON to IPFS: %s", tool.Name)
		toolPinataHash, err := pinJSONToPublicIPFS(json.RawMessage(tool.ToolJson), tool.Name)
		if err != nil {
			return "", fmt.Errorf("failed to pin tool JSON to Pinata: %v", err)
		}
		log.Printf("Pinned tool JSON to public IPFS with CID: %s", toolPinataHash)
		ioObject["tool"] = map[string]interface{}{
			"cid": toolPinataHash,
		}

		for _, inputFile := range inputFiles {
			log.Printf("Downloading input file: %s", inputFile.Filename)
			inputTempFilePath, err := ipfs.DownloadFileToTemp(inputFile.CID, inputFile.Filename)
			if err != nil {
				return "", fmt.Errorf("failed to download input file: %v", err)
			}
			defer os.Remove(inputTempFilePath)

			log.Printf("Pinning input file to IPFS: %s", inputFile.Filename)
			inputPinataHash, err := pinFileToPublicIPFS(inputTempFilePath, inputFile.Filename)
			if err != nil {
				return "", fmt.Errorf("failed to pin input file to Pinata: %v", err)
			}
			log.Printf("Pinned input file to public IPFS with CID: %s", inputPinataHash)
			ioObject["inputs"] = append(ioObject["inputs"].([]map[string]interface{}), map[string]interface{}{
				"cid":      inputPinataHash,
				"filename": inputFile.Filename,
			})
		}

		for _, outputFile := range outputFiles {
			log.Printf("Downloading output file: %s", outputFile.Filename)
			outputTempFilePath, err := ipfs.DownloadFileToTemp(outputFile.CID, outputFile.Filename)
			if err != nil {
				return "", fmt.Errorf("failed to download output file: %v", err)
			}
			defer os.Remove(outputTempFilePath)

			log.Printf("Pinning output file to IPFS: %s", outputFile.Filename)
			outputPinataHash, err := pinFileToPublicIPFS(outputTempFilePath, outputFile.Filename)
			if err != nil {
				return "", fmt.Errorf("failed to pin output file to Pinata: %v", err)
			}
			log.Printf("Pinned output file to public IPFS with CID: %s", outputPinataHash)
			ioObject["outputs"] = append(ioObject["outputs"].([]map[string]interface{}), map[string]interface{}{
				"cid":      outputPinataHash,
				"filename": outputFile.Filename,
			})
		}

		metadata["flow"] = append(metadata["flow"].([]map[string]interface{}), ioObject)
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %v", err)
	}

	log.Printf("Token Metadata: %s", metadataJSON)

	return string(metadataJSON), nil
}

func pinJSONToPublicIPFS(jsonData json.RawMessage, name string) (string, error) {
	url := "https://api.pinata.cloud/pinning/pinJSONToIPFS"

	payload := map[string]interface{}{
		"pinataContent": jsonData,
		"pinataMetadata": map[string]string{
			"name": name,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("PINATA_API_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result struct {
		IpfsHash string `json:"IpfsHash"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response JSON: %v", err)
	}

	return result.IpfsHash, nil
}

func pinFileToPublicIPFS(filePath, name string) (string, error) {
	url := "https://api.pinata.cloud/pinning/pinFileToIPFS"

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}

	err = writer.WriteField("pinataMetadata", fmt.Sprintf(`{"name":"%s"}`, name))
	if err != nil {
		return "", fmt.Errorf("failed to write pinataMetadata field: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %v", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+os.Getenv("PINATA_API_TOKEN"))

	client := &http.Client{}

	maxRetries := 3
	retryDelay := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		rateLimiter.Wait()

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending request to Pinata API: %v", err)
			time.Sleep(retryDelay)
			retryDelay *= 2
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			log.Println("Rate limit exceeded. Retrying after a delay...")
			time.Sleep(retryDelay)
			retryDelay *= 2
			continue
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %v", err)
		}

		var result struct {
			IpfsHash string `json:"IpfsHash"`
		}
		err = json.Unmarshal(respBody, &result)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal response JSON: %v", err)
		}

		return result.IpfsHash, nil
	}

	return "", fmt.Errorf("failed to pin file to Pinata after %d retries", maxRetries)
}

func GenerateAndStoreRecordCID(db *gorm.DB, flow *models.Flow) (string, error) {
	log.Println("Generating token metadata...")
	metadataJSON, err := BuildTokenMetadata(db, flow)
	if err != nil {
		return "", fmt.Errorf("failed to build token metadata: %v", err)
	}
	log.Println("Generated token metadata.")

	log.Println("Pinning token metadata to IPFS...")
	metadataCID, err := pinJSONToPublicIPFS(json.RawMessage(metadataJSON), flow.Name+"_record_metadata.json")
	if err != nil {
		return "", fmt.Errorf("failed to pin token metadata to Pinata: %v", err)
	}
	log.Printf("Pinned token metadata to public IPFS with CID: %s", metadataCID)

	log.Println("Updating Flow's RecordCID...")
	flow.RecordCID = metadataCID
	if err := db.Save(flow).Error; err != nil {
		return "", fmt.Errorf("failed to update Flow's RecordCID: %v", err)
	}
	log.Println("Updated Flow's RecordCID.")

	return metadataCID, nil
}

func MintNFT(db *gorm.DB, flow *models.Flow, metadataCID string) error {
	if autotaskWebhook == "" {
		return fmt.Errorf("AUTOTASK_WEBHOOK must be set")
	}

	log.Println("Triggering minting process via Defender Autotask...")

	data := postData{
		RecipientAddress: flow.WalletAddress,
		Cid:              metadataCID,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON data: %v", err)
	}

	req, err := http.NewRequest("POST", autotaskWebhook, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("error unmarshaling response JSON: %v", err)
	}

	if result.Status != "success" {
		return fmt.Errorf("minting process failed: %s", string(body))
	}

	log.Println("Minting process successful.")

	return nil
}
