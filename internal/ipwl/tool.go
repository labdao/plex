package ipwl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
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

type Tool struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Author      string                `json:"author"`
	GitHub      string                `json:"github"`
	Paper       string                `json:"paper"`
	Task        string                `json:"task"`
	BaseCommand []string              `json:"baseCommand"`
	Arguments   []string              `json:"arguments"`
	DockerPull  string                `json:"dockerPull"`
	GpuBool     bool                  `json:"gpuBool"`
	MemoryGB    *int                  `json:"memoryGB"`
	Cpu         *float64              `json:"cpu"`
	NetworkBool bool                  `json:"networkBool"`
	Inputs      map[string]ToolInput  `json:"inputs"`
	Outputs     map[string]ToolOutput `json:"outputs"`
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

	toolInfo.Name = tool.Name

	return tool, toolInfo, nil
}

func toolToCmd(toolConfig Tool, ioEntry IO) ([]string, error) {
	var arguments []string

	for _, arg := range toolConfig.Arguments {
		placeholderRegex := regexp.MustCompile(`\$\((inputs\.[a-zA-Z0-9_]+)(\.value)?\)`)
		matches := placeholderRegex.FindAllStringSubmatch(arg, -1)

		for _, match := range matches {
			placeholder := match[0]
			inputKey := strings.TrimPrefix(match[1], "inputs.")

			inputValue, ok := ioEntry.Inputs[inputKey]
			if !ok {
				return nil, fmt.Errorf("input key %s not found in IO entry", inputKey)
			}

			var replacement string
			// Determine the type of inputValue and process accordingly
			switch v := inputValue.(type) {
			case []interface{}:
				// Directly marshal to JSON, as we're already asserting types on insertion
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal array for key %s: %s", inputKey, err)
				}
				replacement = string(jsonBytes)
			case string, bool, float64, int: // Add other types as needed
				// Single values are not JSON arrays, so marshal directly to JSON
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal value for key %s: %s", inputKey, err)
				}
				replacement = string(jsonBytes)
			default:
				return nil, fmt.Errorf("unsupported type for key %s", inputKey)
			}

			arg = strings.Replace(arg, placeholder, replacement, -1)
		}

		arguments = append(arguments, arg)
	}

	return arguments, nil
}
