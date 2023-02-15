package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func InstructionToBacalhauCmd(cid, container, cmd string, gpu bool) string {
	// TODO allow overrides for gpu memory and network flags
	gpuFlag := ""
	if gpu {
		gpuFlag = "--gpu 1 "
	}
	return `bacalhau docker run --network full ` + gpuFlag + `--memory 12gb -i ` + fmt.Sprintf(cid) + ` ` + fmt.Sprintf(container) + ` -- ` + fmt.Sprintf(cmd)
}


func RunBacalhauCmd(cmdString string) {
	args := strings.Fields(cmdString)
	fmt.Println(args)
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Command failed: %s\n", err)
		return
	}
	fmt.Printf("Output: %s\n", out)
}
