package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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
			"tool": map[string]interface{}{
				"cid":          tool.CID,
				"name":         tool.Name,
				"container":    tool.Container,
				"memory":       tool.Memory,
				"cpu":          tool.Cpu,
				"gpu":          tool.Gpu,
				"network":      tool.Network,
				"display":      tool.Display,
				"taskCategory": tool.TaskCategory,
			},
			"inputs":  []map[string]interface{}{},
			"outputs": []map[string]interface{}{},
			"state":   job.State,
			"errMsg":  job.Error,
		}

		for _, inputFile := range inputFiles {
			ioObject["inputs"] = append(ioObject["inputs"].([]map[string]interface{}), map[string]interface{}{
				"cid":      inputFile.CID,
				"filename": inputFile.Filename,
			})
		}

		for _, outputFile := range outputFiles {
			ioObject["outputs"] = append(ioObject["outputs"].([]map[string]interface{}), map[string]interface{}{
				"cid":      outputFile.CID,
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

func GenerateAndStoreRecordCID(db *gorm.DB, flow *models.Flow) (string, error) {
	metadataJSON, err := BuildTokenMetadata(db, flow)
	if err != nil {
		return "", fmt.Errorf("failed to build token metadata: %v", err)
	}

	tempFile, err := ioutil.TempFile("", "metadata-*.json")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(metadataJSON); err != nil {
		return "", fmt.Errorf("failed to write metadata to file: %v", err)
	}
	tempFile.Close()

	metadataCID, err := ipfs.PinFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to pin metadata file: %v", err)
	}

	flow.RecordCID = metadataCID
	if err := db.Save(flow).Error; err != nil {
		return "", fmt.Errorf("failed to update Flow's RecordCID: %v", err)
	}

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
