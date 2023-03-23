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

func CreateTool(toolName string, toolFilePath, cid string, paramOverrides map[string]string) (Tool, error) {
	tool, err := ReadTools(toolName, toolFilePath)
	if err != nil {
		return tool, err
	}
	tool.Arguments[0] = formatCmd(tool.Arguments[0], paramOverrides)
	tool.Requirements[0]["dockerPull"] = formatCmd(tool.Requirements[0]["dockerPull"], paramOverrides)
	return tool, nil
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

func formatCmd(cmd string, params map[string]string) (formatted string) {
	// this requires string inputs to have `%{paramKeyX}s %{paramKeyY}s"` formatting
	formatted = cmd
	for key, val := range params {
		formatted = strings.Replace(formatted, "%{"+key+"}s", val, -1)
	}
	return
}
