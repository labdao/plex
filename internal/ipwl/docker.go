package ipwl

import (
	"fmt"
	"os/exec"
)

func runDockerCmd(dockerCmd string) error {
	cmd := exec.Command("/bin/sh", "-c", dockerCmd)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("error running Docker command: %w, output: %s", err, string(output))
	}

	fmt.Println("Docker command output:", string(output))

	return nil
}
