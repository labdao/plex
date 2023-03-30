package ipwl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	BaseCommand []string               `json:"baseCommand"`
	Arguments   []string               `json:"arguments"`
	DockerPull  string                 `json:"dockerPull"`
	GpuBool     bool                   `json:"gpuBool"`
	Inputs      map[string]interface{} `json:"inputs"`
}

func readToolConfig(filePath string) (Tool, error) {
	var tool Tool

	file, err := os.Open(filePath)
	if err != nil {
		return tool, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return tool, fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(bytes, &tool)
	if err != nil {
		return tool, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return tool, nil
}
