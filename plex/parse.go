package main

import (
	"os"
	"path/filepath"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

/*
docker container setup

docker build -t diffdock -f Docker/DiffDock/Dockerfile Docker/DiffDock/
docker tag diffdock openzyme/diffdock
docker push openzyme/diffdock

bacalhau server setup:
sudo sysctl -w net.core.rmem_max=2500000

ipfs id

# change 4001 to 5001
export IPFS_CONNECT=/ip4/127.0.0.1/tcp/5001/p2p/12D3KooWPH1BpPfNXwkf778GMP2H5z7pwjKVQFnA5NS3DngU7pxG

LOG_LEVEL=debug bacalhau serve --job-selection-accept-networked --limit-total-gpu 1 --limit-total-memory 12gb --ipfs-connect $IPFS_CONNECT

# also make sure data is on ipfs node
*/

type Instruction struct {
	App       string            `json:"app"`
	InputCIDs []string          `json:"input_cids"`
	Container string            `json:"container"`
	Params    map[string]string `json:"params"`
	Cmd       string            `json:"cmd"`
	CmdHelper bool              `json:"cmd_helper"`
}

func CreateInstruction(app string, instuctionFilePath, inputDirPath string, paramOverrides map[string]string) (Instruction, error) {
	instruction, err := readInstructions(app, instuctionFilePath)
	if err != nil {
		return instruction, err
	}
	instruction.Params = overwriteParams(instruction.Params, paramOverrides)
	instruction.Cmd = formatCmd(instruction.Cmd, instruction.Params)
	cid, err := createInputCID(inputDirPath, instruction.CmdHelper, instruction.Cmd)
	if err != nil {
		return instruction, err
	}
	instruction.InputCIDs = append(instruction.InputCIDs, cid)
	return instruction, nil
}

func readInstructions(app string, filepath string) (Instruction, error) {
	fileContents, err := ioutil.ReadFile(filepath)
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

	return instruction, fmt.Errorf("No instruction found for app %s", app)
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
		formatted = strings.Replace(formatted, "%{"+key+"}s", fmt.Sprintf("%s", val), -1)
	}
	return
}

func createInputCID(inputDirPath string, cmdHelper bool, cmd string) (string, error) {
	// if cmdHelper this will push formattedCmd to helper.sh
	// this will then use the 2 be merged ipfs function to return a cid
	if (cmdHelper) {
		err := createHelperFile(inputDirPath, cmd)
		if err != nil {
			return "", err
		}
	}
	return "QmZGavZusys5SrgyQB69iJwWL5tAbXrYeyoJBcjdJsp3mR", nil
}

func createHelperFile(dirPath string, contents string) error {
	fileName := filepath.Join(dirPath, "helper.sh")
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("#!/bin/bash\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString(contents)
	if err != nil {
		return err
	}

	err = os.Chmod(fileName, 0755)
	if err != nil {
		return err
	}

	return nil
}
