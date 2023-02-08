package plexus

import (
	"fmt"
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

// Structure for AppConfig
type appStruct2 struct {
	app     string      `json:"app"`
	inputs  [][2]string `json:"inputs"`
	outputs []string    `json:"outputs"`
	container string
	params map[string]string
	cmd string
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

/*
func createInstruction(appConfig appStruct2, inputCIDs []string, paramOverrides map[string][string]) (appStruct2) {
	var instruction appStruct2
	json.Unmarshal([]byte(appConfig), &instruct)
	instruction.params = overrideParams(instruction.params, paramOverrides)
	instruction.cmd = formatCmd
}
*/

/*
func instructJSONtoBacalhauCmd(instructJSON: json) -> string {

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
