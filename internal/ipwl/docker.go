package ipwl

import (
	"fmt"
	"os/exec"
)

func runDockerCmd(dockerCmd string) ([]byte, error) {
	cmd := exec.Command("/bin/sh", "-c", dockerCmd)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return output, fmt.Errorf("error running Docker command: %w, output: %s", err, string(output))
	}

	return output, nil
}
