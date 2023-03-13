package docker

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
)

func InstructionToDockerCmd(container, cmd, jobDir string, gpu bool) string {
	outputsDir := filepath.Join(jobDir, "outputs")
	gpuFlag := ""
	if gpu {
		gpuFlag = "--gpus"
	}
	return fmt.Sprintf("docker run %s -v %s:/inputs -v %s:/outputs %s /bin/bash -c '%s'", gpuFlag, jobDir, outputsDir, container, cmd)
}

func RunDockerCmd(container, cmd, jobDir string, gpu bool) error {
	dockerCmd := InstructionToDockerCmd(container, cmd, jobDir, gpu)

	// make output dir

	cmdExec := exec.Command(dockerCmd)

	stdout, err := cmdExec.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmdExec.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := cmdExec.Wait(); err != nil {
		return err
	}
	return err
}
