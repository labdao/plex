package ipwl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	BaseCommand []string               `json:"baseCommand"`
	Arguments   []string               `json:"arguments"`
	DockerPull  string                 `json:"dockerPull"`
	GpuBool     bool                   `json:"gpuBool"`
	Inputs      map[string]interface{} `json:"inputs"`
	Outputs     map[string]interface{} `json:"outputs"`
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

func ToolToDockerCmd(toolConfig Tool, ioEntry IO, outputDirPath string) (string, error) {
	inputVolumes := ""
	for _, input := range ioEntry.Inputs {
		inputFilepath := input.(map[string]interface{})["filepath"].(string)
		inputVolumes += fmt.Sprintf("-v %s:/inputs ", inputFilepath)
	}

	arguments := strings.Join(toolConfig.Arguments, " ")

	placeholderRegex := regexp.MustCompile(`\$\((inputs\..+?\.filepath)\)`)
	matches := placeholderRegex.FindAllStringSubmatch(arguments, -1)

	for _, match := range matches {
		placeholder := match[0]
		key := strings.TrimSuffix(strings.TrimPrefix(match[1], "inputs."), ".filepath")
		arguments = strings.Replace(arguments, placeholder, ioEntry.Inputs[key].(map[string]interface{})["filepath"].(string), -1)
	}

	dockerCmd := fmt.Sprintf("docker %s-v %s:/outputs run %s %s %s", inputVolumes, outputDirPath, toolConfig.DockerPull, strings.Join(toolConfig.BaseCommand, " "), arguments)

	return dockerCmd, nil
}
