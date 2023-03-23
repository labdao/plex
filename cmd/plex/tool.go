package plex

import (
	"encoding/json"
	"fmt"
	"os"
)

type Tool struct {
	Class       string                 `json:"class"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	BaseCommand []string               `json:"baseCommand"`
	Arguments   []string               `json:"arguments"`
	Requirements []map[string]string   `json:"requirements"`
	Inputs      map[string]interface{} `json:"inputs"`
	Outputs     map[string]interface{} `json:"outputs"`
}

func ReadTool(toolName string, filepath string) (Tool, error) {
	fileContents, err := os.ReadFile(filepath)
	var tool Tool
	if err != nil {
		return tool, err
	}
	err = json.Unmarshal(fileContents, &tool)
	if err != nil {
		return tool, err
	}

	if tool.Name == toolName {
		return tool, nil
	}

	return tool, fmt.Errorf("no tool found for tool name %s", toolName)
}