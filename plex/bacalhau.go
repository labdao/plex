package main

import (
	"fmt"
	"os/exec"
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

func main() {
	// cmd := exec.Command("bacalhau", "docker", "run", "ubuntu", "echo", "Hello World")
	// cmd := exec.Command("bacalhau", "docker", "run", "--gpu", "1", "--memory", "12gb", "nvidia/cuda:11.0.3-base-ubuntu20.04", "nvidia-smi")
	// cmd := exec.Command("bacalhau", "docker", "run", "--gpu", "1", "--memory", "12gb", "-i", "QmZGavZusys5SrgyQB69iJwWL5tAbXrYeyoJBcjdJsp3mR", "nvidia/cuda:11.0.3-base-ubuntu20.04", "nvidia-smi")
	// cmd := exec.Command("bacalhau", "docker", "run", "--gpu", "1", "--memory", "12gb", "-i", "QmZGavZusys5SrgyQB69iJwWL5tAbXrYeyoJBcjdJsp3mR", "openzyme/diffdock", "ls", "-lah", "/inputs", ">", "/outputs/ls.txt")
	cmd := exec.Command("bacalhau", "docker", "run", "--network", "full", "--gpu", "1", "--memory", "12gb", "-i", "QmZGavZusys5SrgyQB69iJwWL5tAbXrYeyoJBcjdJsp3mR", "openzyme/diffdock", "./test.sh")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Command failed: %s\n", err)
		return
	}
	fmt.Printf("Output: %s\n", out)
}
