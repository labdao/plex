package ipwl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

func ReadToolConfig(filePath string) (Tool, error) {
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

func toolToDockerCmd(toolConfig Tool, ioEntry IO, outputDirPath string) (string, error) {
	inputVolumes := ""
	for inputName, input := range ioEntry.Inputs {
		inputVolumes += fmt.Sprintf("--mount type=bind,source=%s,target=/inputs/%s ", input.FilePath, inputName)
	}

	arguments := strings.Join(toolConfig.Arguments, " ")

	placeholderRegex := regexp.MustCompile(`\$\((inputs\..+?(\.filepath|\.basename|\.ext))\)`)
	matches := placeholderRegex.FindAllStringSubmatch(arguments, -1)

	for _, match := range matches {
		placeholder := match[0]
		key := strings.TrimSuffix(strings.TrimPrefix(match[1], "inputs."), ".filepath")
		key = strings.TrimSuffix(key, ".basename")
		key = strings.TrimSuffix(key, ".ext")

		var replacement string
		input := ioEntry.Inputs[key]

		switch match[2] {
		case ".filepath":
			replacement = input.FilePath
		case ".basename":
			replacement = filepath.Base(input.FilePath)
		case ".ext":
			ext := filepath.Ext(input.FilePath)
			replacement = strings.TrimPrefix(ext, ".")
		}

		arguments = strings.Replace(arguments, placeholder, replacement, -1)
	}

	dockerCmd := fmt.Sprintf(`docker run %s-v %s:/outputs %s %s "%s"`, inputVolumes, outputDirPath, toolConfig.DockerPull, strings.Join(toolConfig.BaseCommand, " "), arguments)

	return dockerCmd, nil
}
