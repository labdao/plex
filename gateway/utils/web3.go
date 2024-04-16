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

	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/internal/ipfs"
	"gorm.io/gorm"
)

var autotaskWebhook = os.Getenv("AUTOTASK_WEBHOOK")

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

		toolPinataHash, err := pinJSONToIPFSEndpoint(json.RawMessage(tool.ToolJson), tool.Name)
		if err != nil {
			return "", fmt.Errorf("failed to pin tool JSON to Pinata: %v", err)
		}
		ioObject["tool"] = map[string]interface{}{
			"cid": toolPinataHash,
		}

		for _, inputFile := range inputFiles {
			inputTempFilePath, err := ipfs.DownloadFileToTemp(inputFile.CID, inputFile.Filename)
			if err != nil {
				return "", fmt.Errorf("failed to download input file: %v", err)
			}
			defer os.Remove(inputTempFilePath)

			inputPinataHash, err := pinFileToIPFSEndpoint(inputTempFilePath)
			if err != nil {
				return "", fmt.Errorf("failed to pin input file to Pinata: %v", err)
			}
			ioObject["inputs"] = append(ioObject["inputs"].([]map[string]interface{}), map[string]interface{}{
				"cid":      inputPinataHash,
				"filename": inputFile.Filename,
			})
		}

		for _, outputFile := range outputFiles {
			outputTempFilePath, err := ipfs.DownloadFileToTemp(outputFile.CID, outputFile.Filename)
			if err != nil {
				return "", fmt.Errorf("failed to download output file: %v", err)
			}
			defer os.Remove(outputTempFilePath)

			outputPinataHash, err := pinFileToIPFSEndpoint(outputTempFilePath)
			if err != nil {
				return "", fmt.Errorf("failed to pin output file to Pinata: %v", err)
			}
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

func pinJSONToIPFSEndpoint(jsonData json.RawMessage, name string) (string, error) {
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

func pinFileToIPFSEndpoint(filePath string) (string, error) {
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
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

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

func GenerateAndStoreRecordCID(db *gorm.DB, flow *models.Flow) (string, error) {
	metadataJSON, err := BuildTokenMetadata(db, flow)
	if err != nil {
		return "", fmt.Errorf("failed to build token metadata: %v", err)
	}

	metadataCID, err := pinJSONToIPFSEndpoint(json.RawMessage(metadataJSON), flow.Name+"_record_metadata.json")
	if err != nil {
		return "", fmt.Errorf("failed to pin token metadata to Pinata: %v", err)
	}

	flow.RecordCID = metadataCID
	if err := db.Save(flow).Error; err != nil {
		return "", fmt.Errorf("failed to update Flow's RecordCID: %v", err)
	}

	return metadataCID, nil
}

// func GenerateAndStoreRecordCID(db *gorm.DB, flow *models.Flow) (string, error) {
// 	metadataJSON, err := BuildTokenMetadata(db, flow)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to build token metadata: %v", err)
// 	}

// 	payload := map[string]interface{}{
// 		"pinataOptions": map[string]int{
// 			"cidVersion": 0,
// 		},
// 		"pinataMetadata": map[string]string{
// 			"name": flow.Name + "_record_metadata.json",
// 		},
// 		"pinataContent": json.RawMessage(metadataJSON),
// 	}

// 	jsonPayload, err := json.Marshal(payload)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal JSON payload: %v", err)
// 	}

// 	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinJSONToIPFS", bytes.NewBuffer(jsonPayload))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create HTTP request: %v", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+os.Getenv("PINATA_API_TOKEN"))

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to send HTTP request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	responseBody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read response body: %v", err)
// 	}

// 	var pinataResponse struct {
// 		IpfsHash string `json:"IpfsHash"`
// 	}
// 	err = json.Unmarshal(responseBody, &pinataResponse)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to unmarshal response JSON: %v", err)
// 	}

// 	flow.RecordCID = pinataResponse.IpfsHash
// 	if err := db.Save(flow).Error; err != nil {
// 		return "", fmt.Errorf("failed to update Flow's RecordCID: %v", err)
// 	}

// 	return pinataResponse.IpfsHash, nil
// }

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
