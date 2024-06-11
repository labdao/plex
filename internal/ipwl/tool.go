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

type ToolInput struct {
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

type ToolOutput struct {
	Type string   `json:"type"`
	Item string   `json:"item"`
	Glob []string `json:"glob"`
}

type ToolType string

const (
	ToolTypeBacalhau ToolType = "bacalhau"
	ToolTypeRay      ToolType = "ray"
)

type Tool struct {
	Name                 string                `json:"name"`
	Description          string                `json:"description"`
	Guide                string                `json:"guide"`
	Author               string                `json:"author"`
	GitHub               string                `json:"github"`
	Paper                string                `json:"paper"`
	Task                 string                `json:"task"`
	CheckpointCompatible bool                  `json:"checkpointCompatible"`
	BaseCommand          []string              `json:"baseCommand"`
	Arguments            []string              `json:"arguments"`
	DockerPull           string                `json:"dockerPull"`
	GpuBool              bool                  `json:"gpuBool"`
	MemoryGB             *int                  `json:"memoryGB"`
	Cpu                  *float64              `json:"cpu"`
	NetworkBool          bool                  `json:"networkBool"`
	Inputs               map[string]ToolInput  `json:"inputs"`
	Outputs              map[string]ToolOutput `json:"outputs"`
	TaskCategory         string                `json:"taskCategory"`
	MaxRunningTime       int                   `json:"maxRunningTime"`
	ToolType             ToolType              `json:"toolType"`
	RayServiceEndpoint   string                `json:"rayServiceEndpoint"`
	XAxis                string                `json:"xAxis"`
	YAxis                string                `json:"yAxis"`
}

func ReadToolConfig(toolPath string, db *gorm.DB) (Tool, ToolInfo, error) {
	var ipwltool Tool
	var dbTool models.Tool
	var toolInfo ToolInfo
	var err error

	s3client, err := s3client.NewS3Client()
	if err != nil {
		return ipwltool, toolInfo, err
	}

	err = db.Where("cid = ?", toolPath).First(&dbTool).Error
	if err != nil {
		return ipwltool, toolInfo, fmt.Errorf("failed to get tool from database: %w", err)
	}
	toolInfo.S3 = toolPath
	toolInfo.Name = dbTool.Name
	bucket, key, err := s3client.GetBucketAndKeyFromURI(dbTool.S3URI)
	if err != nil {
		return ipwltool, toolInfo, fmt.Errorf("failed to get bucket and key from URI: %w", err)
	}
	fileName := filepath.Base(key)
	err = s3client.DownloadFile(bucket, key, fileName)
	if err != nil {
		return ipwltool, toolInfo, fmt.Errorf("failed to download file: %w", err)
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
		return ipwltool, toolInfo, fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(bytes, &ipwltool)
	if err != nil {
		return ipwltool, toolInfo, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return ipwltool, toolInfo, nil
}
