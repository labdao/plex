package ipwl

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/labdao/plex/gateway/models"
	s3client "github.com/labdao/plex/internal/s3"
	"gorm.io/gorm"
)

type ModelInput struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Array       bool        `json:"array"`
	Glob        []string    `json:"glob"`
	Default     interface{} `json:"default"`
	Min         string      `json:"min"`
	Max         string      `json:"max"`
	Example     string      `json:"example"`
	Grouping    string      `json:"grouping"`
	Position    string      `json:"position"`
	Required    bool        `json:"required"`
}

type ModelOutput struct {
	Type string   `json:"type"`
	Item string   `json:"item"`
	Glob []string `json:"glob"`
}

type ModelType string

const (
	ModelTypeBacalhau ModelType = "bacalhau"
	ModelTypeRay      ModelType = "ray"
)

type Model struct {
	Name                 string                 `json:"name"`
	Description          string                 `json:"description"`
	Guide                string                 `json:"guide"`
	Author               string                 `json:"author"`
	GitHub               string                 `json:"github"`
	Paper                string                 `json:"paper"`
	Task                 string                 `json:"task"`
	CheckpointCompatible bool                   `json:"checkpointCompatible"`
	BaseCommand          []string               `json:"baseCommand"`
	Arguments            []string               `json:"arguments"`
	DockerPull           string                 `json:"dockerPull"`
	GpuBool              bool                   `json:"gpuBool"`
	MemoryGB             *int                   `json:"memoryGB"`
	Cpu                  *float64               `json:"cpu"`
	NetworkBool          bool                   `json:"networkBool"`
	Inputs               map[string]ModelInput  `json:"inputs"`
	Outputs              map[string]ModelOutput `json:"outputs"`
	TaskCategory         string                 `json:"taskCategory"`
	MaxRunningTime       int                    `json:"maxRunningTime"`
	ComputeCost          int                    `json:"computeCost"`
	ModelType            ModelType              `json:"modelType"`
	RayServiceEndpoint   string                 `json:"rayServiceEndpoint"`
	XAxis                string                 `json:"xAxis"`
	YAxis                string                 `json:"yAxis"`
}

func ReadModelConfig(modelPath string, db *gorm.DB) (Model, ModelInfo, error) {
	var ipwlmodel Model
	var dbModel models.Model
	var modelInfo ModelInfo
	var err error

	s3client, err := s3client.NewS3Client()
	if err != nil {
		return ipwlmodel, modelInfo, err
	}

	err = db.Where("cid = ?", modelPath).First(&dbModel).Error
	if err != nil {
		return ipwlmodel, modelInfo, fmt.Errorf("failed to get model from database: %w", err)
	}
	modelInfo.S3 = modelPath
	modelInfo.Name = dbModel.Name
	bucket, key, err := s3client.GetBucketAndKeyFromURI(dbModel.S3URI)
	if err != nil {
		return ipwlmodel, modelInfo, fmt.Errorf("failed to get bucket and key from URI: %w", err)
	}
	fileName := filepath.Base(key)
	err = s3client.DownloadFile(bucket, key, fileName)
	if err != nil {
		return ipwlmodel, modelInfo, fmt.Errorf("failed to download file: %w", err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
	}
	defer file.Close()
	defer os.Remove(fileName)

	fmt.Println("File opened successfully")
	bytes, err := io.ReadAll(file)
	if err != nil {
		return ipwlmodel, modelInfo, fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(bytes, &ipwlmodel)
	if err != nil {
		return ipwlmodel, modelInfo, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return ipwlmodel, modelInfo, nil
}
