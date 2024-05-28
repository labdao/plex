package ipwl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/labdao/plex/internal/ipfs"
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
	RayServiceURL        string                `json:"rayServiceURL"`
	XAxis                string                `json:"xAxis"`
	YAxis                string                `json:"yAxis"`
}

func ReadToolConfig(toolPath string) (Tool, ToolInfo, error) {
	var tool Tool
	var toolInfo ToolInfo
	var toolFilePath string
	var cid string
	var err error

	if ipfs.IsValidCID(toolPath) {
		toolInfo.IPFS = toolPath
		toolFilePath, err = ipfs.DownloadToTempDir(toolPath)
		if err != nil {
			return tool, toolInfo, err
		}

		fileInfo, err := os.Stat(toolFilePath)
		if err != nil {
			return tool, toolInfo, err
		}

		// If the downloaded content is a directory, search for a .json file in it
		if fileInfo.IsDir() {
			files, err := ioutil.ReadDir(toolFilePath)
			if err != nil {
				return tool, toolInfo, err
			}

			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".json") {
					toolFilePath = path.Join(toolFilePath, file.Name())
					break
				}
			}
		}
	} else {
		if _, err := os.Stat(toolPath); err == nil {
			toolFilePath = toolPath
			cid, err = ipfs.WrapAndPinFile(toolFilePath)
			if err != nil {
				return tool, toolInfo, fmt.Errorf("failed to pin tool file")
			}
			toolInfo.IPFS = cid
		} else {
			return tool, toolInfo, fmt.Errorf("tool not found")
		}
	}

	file, err := os.Open(toolFilePath)
	if err != nil {
		return tool, toolInfo, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return tool, toolInfo, fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(bytes, &tool)
	if err != nil {
		return tool, toolInfo, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if tool.ToolType == "" {
		tool.ToolType = ToolTypeBacalhau
	}

	toolInfo.Name = tool.Name

	return tool, toolInfo, nil
}
