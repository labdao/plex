package main

import (
	"encoding/json"
//	"flag"
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
	Class        string   `json:"class"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	BaseCommand  []string `json:"baseCommand"`
	Arguments    []string `json:"arguments"`
	Requirements []struct {
		Class      string `json:"class"`
		DockerPull string `json:"dockerPull"`
	} `json:"requirements"`
	Inputs  map[string]InputFile    `json:"inputs"`
	Outputs map[string]interface{} `json:"outputs"`
}

func readToolConfig(path string) (*ToolConfig, error) {
	// Read tool.json file
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Unmarshal tool.json to ToolConfig type
	var toolConfig ToolConfig
	err = json.Unmarshal(byteValue, &toolConfig)
	if err != nil {
		return nil, err
	}
	return &toolConfig, nil
}

func readInputFile(path string) (*InputJSON, error) {
	// Read input.json file
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Unmarshal input.json to InputJSON type
	var input InputJSON
	err = json.Unmarshal(byteValue, &input)
	if err != nil {
		return nil, err
	}
	return &input, nil
}

//do not use this function
func enrichInputFile(inputPath string) error {
    // Read input.json file
    input, err := readInputFile(inputPath)
    if err != nil {
        return err
    }

    // Enrich input files
    for inputName, inputFile := range *input {
        inputFileMap := inputFile.(map[string]interface{})
        inputBasename := inputFileMap["basename"].(string)
        inputNameRoot := strings.TrimSuffix(inputBasename, filepath.Ext(inputBasename))
        inputNameExt := filepath.Ext(inputBasename)

        inputFileMap["nameroot"] = inputNameRoot
        inputFileMap["nameext"] = inputNameExt
        inputFileMap["path"] = fmt.Sprintf("/inputs/%s", inputBasename)

        // Update input object with enriched input file
        (*input)[inputName] = inputFileMap
    }

    // Write enriched input to new file
    outputFilename := strings.TrimSuffix(inputPath, ".json") + "_enriched.json"
    outputFile, err := os.Create(outputFilename)
    if err != nil {
        return err
    }
    defer outputFile.Close()

    outputBytes, err := json.MarshalIndent(input, "", "    ")
    if err != nil {
        return err
    }

    _, err = outputFile.Write(outputBytes)
    if err != nil {
        return err
    }

    fmt.Printf("Enriched input written to %s\n", outputFilename)

    return nil
}


func printInputs() error {
	// Read input.json file
	input, err := readInputFile("inputs.json")
	if err != nil {
		return err
	}

	// Loop through each input file and print out the basename
	for inputName, inputFile := range *input {
		switch inputFile.(type) {
		case map[string]interface{}:
			inputPath := inputFile.(map[string]interface{})["basename"].(string)
			fmt.Printf("%s: %s\n", inputName, inputPath)
		case []interface{}:
			for _, file := range inputFile.([]interface{}) {
				inputPath := file.(map[string]interface{})["basename"].(string)
				fmt.Printf("%s: %s\n", inputName, inputPath)
			}
		}
	}
	return nil
}


func main() {
    printInputs()
}
