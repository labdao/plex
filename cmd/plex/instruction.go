package plex

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Instruction struct {
	App       string            `json:"app"`
	InputCIDs []string          `json:"input_cids"`
	Container string            `json:"container"`
	Params    map[string]string `json:"params"`
	Cmd       string            `json:"cmd"`
}

func CreateInstruction(app string, instuctionFilePath, inputDirPath string, paramOverrides map[string]string) (Instruction, error) {
	instruction, err := ReadInstructions(app, instuctionFilePath)
	if err != nil {
		return instruction, err
	}
	instruction.Params = overwriteParams(instruction.Params, paramOverrides)
	instruction.Cmd = formatCmd(instruction.Cmd, instruction.Params)
	cid, err := CreateInputCID(inputDirPath, instruction.Cmd)
	if err != nil {
		return instruction, err
	}
	instruction.InputCIDs = append(instruction.InputCIDs, cid)
	return instruction, nil
}

func ReadInstructions(app string, filepath string) (Instruction, error) {
	fileContents, err := os.ReadFile(filepath)
	var instruction Instruction
	if err != nil {
		return instruction, err
	}
	lines := strings.Split(string(fileContents), "\n")
	for _, line := range lines {
		err := json.Unmarshal([]byte(line), &instruction)
		if err != nil {
			return instruction, err
		}

		if instruction.App == app {
			return instruction, nil
		}
	}

	return instruction, fmt.Errorf("no instruction found for app %s", app)
}

func overwriteParams(defaultParams, overrideParams map[string]string) (finalParams map[string]string) {
	finalParams = make(map[string]string)
	for key, defaultVal := range defaultParams {
		if overrideVal, ok := overrideParams[key]; ok {
			finalParams[key] = overrideVal
		} else {
			finalParams[key] = defaultVal
		}
	}
	return
}

func formatCmd(cmd string, params map[string]string) (formatted string) {
	// this requires string inputs to have `%{paramKeyX}s %{paramKeyY}s"` formatting
	formatted = cmd
	for key, val := range params {
		formatted = strings.Replace(formatted, "%{"+key+"}s", val, -1)
	}
	return
}
