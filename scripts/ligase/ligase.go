package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Requirement struct {
	DockerPull  string `json:"dockerPull,omitempty"`
	GpuBool     bool   `json:"gpuBool,omitempty"`
	NetworkBool bool   `json:"networkBool,omitempty"`
}

type InputOutput struct {
	Type  string `json:"type"`
	Items string `json:"items,omitempty"`
	Glob  string `json:"glob,omitempty"`
}

type ToolsConfig struct {
	Class        string                 `json:"class"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	BaseCommand  []string               `json:"baseCommand"`
	Arguments    []string               `json:"arguments"`
	Requirements []Requirement          `json:"requirements"`
	Inputs       map[string]InputOutput `json:"inputs"`
	Outputs      map[string]InputOutput `json:"outputs"`
}

type Address struct {
	File string `json:"file"`
}

type Output struct {
	Class   string    `json:"class"`
	Item    string    `json:"item,omitempty"`
	Address []Address `json:"address"`
}

func loadToolsConfig(filePath string) (ToolsConfig, error) {
	var toolsConfig ToolsConfig

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return toolsConfig, fmt.Errorf("failed to open file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return toolsConfig, fmt.Errorf("failed to read file content: %v", err)
	}

	err = json.Unmarshal(byteValue, &toolsConfig)
	if err != nil {
		return toolsConfig, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return toolsConfig, nil
}

func findOutputs(dir string, toolsConfig ToolsConfig) (map[string]Output, error) {
	outputMap := make(map[string]Output)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return outputMap, fmt.Errorf("failed to read directory: %v", err)
	}

	for key, io := range toolsConfig.Outputs {
		var matchingFiles []Address
		for _, file := range files {
			if !file.IsDir() {
				match, err := filepath.Match(io.Glob, file.Name())
				if err != nil {
					return outputMap, fmt.Errorf("failed to match glob pattern: %v", err)
				}
				if match {
					matchingFiles = append(matchingFiles, Address{File: file.Name()})
				}
			}
		}

		output := Output{
			Class:   io.Type,
			Item:    io.Items,
			Address: matchingFiles,
		}

		outputMap[key] = output
	}

	return outputMap, nil
}

func saveOutputMap(dirPath string, outputMap map[string]Output) error {
	outputJSON, err := json.MarshalIndent(outputMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output map: %v", err)
	}

	outputFilePath := filepath.Join(dirPath, "outputs.json")
	err = ioutil.WriteFile(outputFilePath, outputJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write outputs.json file: %v", err)
	}

	return nil
}

func main() {
	toolsConfig, err := loadToolsConfig("inputs/tools.json")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	dir := "inputs"
	outputMap, err := findOutputs(dir, toolsConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Outputs:")
	for key, output := range outputMap {
		fmt.Printf("%s: %+v\n", key, output)
	}

	fmt.Printf("ToolsConfig: %+v\n", toolsConfig)

	outputDir := "outputs"
	err = saveOutputMap(outputDir, outputMap)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Output map saved to outputs/outputs.json")
}
