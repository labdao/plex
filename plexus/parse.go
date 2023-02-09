package plexus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
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

func formatCmd(cmd string, params map[string]string) (formatted string) {
	// this requires string inputs to have `%{paramKeyX}s %{paramKeyY}s"` formatting
	formatted = cmd
	for key, val := range params {
		formatted = strings.Replace(formatted, "%{"+key+"}s", fmt.Sprintf("%s", val), -1)
	}
	return
}

func createInputCID(inputDirPath string, cmdHelper bool, cmd string) (cid string) {
	// if cmdHelper this will push formattedCmd to helper.sh
	// this will then use the 2 be merged ipfs function to return a cid
	cid = "QmZGavZusys5SrgyQB69iJwWL5tAbXrYeyoJBcjdJsp3mR"
	return
}

func createInstruction(app string, instuctionFilePath, inputDirPath string, paramOverrides map[string]string) (Instruction, error) {
	instruction, err := readInstructions(app, instuctionFilePath)
	if err != nil {
		return instruction, err
	}
	instruction.Params = overwriteParams(instruction.Params, paramOverrides)
	instruction.Cmd = formatCmd(instruction.Cmd, instruction.Params)
	instruction.InputCIDs = append(instruction.InputCIDs, createInputCID(inputDirPath, instruction.CmdHelper, instruction.Cmd))
	return instruction, nil
}

/*
func instructiontoBacalhauCmd(instructJSON: json) -> string {

}

func runBacalhauCmd(cmd: str) -> {

}
*/

func bacalhau() {
	// cmd := exec.Command("bacalhau", "docker", "run", "ubuntu", "echo", "Hello World")
	// cmd := exec.Command("bacalhau", "docker", "run", "--gpu", "1", "--memory", "12gb", "nvidia/cuda:11.0.3-base-ubuntu20.04", "nvidia-smi")
	// cmd := exec.Command("bacalhau", "docker", "run", "--gpu", "1", "--memory", "12gb", "-i", "QmZGavZusys5SrgyQB69iJwWL5tAbXrYeyoJBcjdJsp3mR", "nvidia/cuda:11.0.3-base-ubuntu20.04", "nvidia-smi")
	cmd := exec.Command("bacalhau", "docker", "run", "--gpu", "1", "--memory", "12gb", "-i", "bafybeidu5cds5bzjmlt2mmahsfsja6sn5hpdibsfrxwqu5clx5l7q3dvbq", "nvidia/cuda:11.0.3-base-ubuntu20.04", "ls", "-lah", "/inputs", ">", "/outputs/ls.txt")
	// cmd := exec.Command("bacalhau", "docker", "run", "--network", "full", "--gpu", "1", "--memory", "12gb", "-i", "QmZGavZusys5SrgyQB69iJwWL5tAbXrYeyoJBcjdJsp3mR", "openzyme/diffdock:latest", "head", "./test.sh")
	// cmd := exec.Command("bacalhau", "docker", "run", "--gpu", "1", "--memory", "12gb", "-i", "QmZGavZusys5SrgyQB69iJwWL5tAbXrYeyoJBcjdJsp3mR", "nvidia/cuda:11.0.3-base-ubuntu20.04", "nvidia-smi")
	// docker exec -i $container /bin/sh < ./your_script.sh

	// cmd := exec.Command("bacalhau", "docker", "run", "--network", "full", "--gpu", "1", "--memory", "12gb", "-i", "QmZGavZusys5SrgyQB69iJwWL5tAbXrYeyoJBcjdJsp3mR", "openzyme/labdaodiffdock", "./test.sh")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Command failed: %s\n", err)
		return
	}
	fmt.Printf("Output: %s\n", out)
}
