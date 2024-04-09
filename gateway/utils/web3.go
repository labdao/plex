package utils

import (
	"encoding/json"
	"fmt"

	"github.com/labdao/plex/gateway/models"
	"gorm.io/gorm"
)

func BuildTokenMetadata(db *gorm.DB, flow *models.Flow) (string, error) {
	var jobs []models.Job
	if err := db.Where("flow_id = ?", flow.ID).Find(&jobs).Error; err != nil {
		return "", fmt.Errorf("failed to retrieve jobs: %v", err)
	}

	metadata := map[string]interface{}{
		"name":        flow.Name,
		"description": "Research, Reimagined. All Scientists Welcome.",
		"image":       "", // Set the appropriate image URL or IPFS hash
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

	return string(metadataJSON), nil
}
