package docker

import (
	"fmt"
	"path/filepath"
)

func InstructionToDockerCmd(container, cmd, jobDir string, gpu bool) string {
	outputsDir := filepath.Join(jobDir, "outputs")
	gpuFlag := ""
	if gpu {
		gpuFlag = "--gpus"
	}
	return fmt.Sprintf("docker run %s -v %s:/inputs -v %s:/outputs %s -- /bin/bash -c '%s'", gpuFlag, jobDir, outputsDir, container, cmd)
}
