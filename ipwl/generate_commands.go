package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type InputFile struct {
	Class    string `json:"class"`
	Basename string `json:"basename"`
	Address  struct {
		File string `json:"file"`
		IPFS string `json:"ipfs"`
	} `json:"address"`
}

type InputJSON map[string]interface{}

type ToolConfig struct {
	Class       string   `json:"class"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	BaseCommand []string `json:"baseCommand"`
	Arguments   []string `json:"arguments"`
	Requirements []struct {
		Class      string `json:"class"`
		DockerPull string `json:"dockerPull"`
	} `json:"requirements"`
	Inputs  map[string]InputFile `json:"inputs"`
	Outputs map[string]interface{} `json:"outputs"`
}

func main() {
	// Parse command line arguments
	inputsFile := flag.String("inputs", "inputs.json", "Path to the inputs.json file")
	toolConfigFile := flag.String("tool", "equibind.json", "Path to the tool configuration JSON file")
	outputFile := flag.String("output", "command.sh", "Path to the output shell script")
	flag.Parse()

	// Read input.json file
	jsonFile, err := os.Open(*inputsFile)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Unmarshal input.json to InputJSON type
	var input InputJSON
	err = json.Unmarshal(byteValue, &input)
	if err != nil {
		panic(err)
	}

	// Read tool.json file
	jsonFile, err = os.Open(*toolConfigFile)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ = ioutil.ReadAll(jsonFile)

	// Unmarshal tool.json to ToolConfig type
	var toolConfig ToolConfig
	err = json.Unmarshal(byteValue, &toolConfig)
	if err != nil {
		panic(err)
	}

	// Create output file
	output, err := os.Create(*outputFile)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	// Loop through each input file and generate shell commands
	for inputName, inputFile := range toolConfig.Inputs {
		inputPath := fmt.Sprintf("/inputs/%s", inputFile.Basename)
		inputNameRoot := strings.TrimSuffix(inputFile.Basename, filepath.Ext(inputFile.Basename))
		inputExt := filepath.Ext(inputFile.Basename)

		// Generate shell command based on baseCommand and arguments in tool config
		command := strings.Join(toolConfig.BaseCommand, " ")
		for _, arg := range toolConfig.Arguments {
			arg = strings.ReplaceAll(arg, "$(inputs."+inputName+".path)", inputPath)
			arg = strings.ReplaceAll(arg, "$(inputs."+inputName+".nameroot)", inputNameRoot)
			arg = strings.ReplaceAll(arg, "$(inputs."+inputName+".nameext)", inputExt)
			command += " " + arg
		}

		// Write the generated shell command to the output file
		_, err := output.WriteString(command + "\n")
		if err != nil {
			panic(err)
		}
	}
}
