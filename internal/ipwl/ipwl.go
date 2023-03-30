package ipwl

import (
	"fmt"
	"os"
	"path/filepath"
)

func processIOList(ioList []IO, jobDir, ioJsonPath string) error {
	for i, ioEntry := range ioList {
		err := processIOTask(ioEntry, i, jobDir, ioJsonPath)
		if err != nil {
			return fmt.Errorf("error processing IO task at index %d: %w", i, err)
		}
	}

	return nil
}

func processIOTask(ioEntry IO, index int, jobDir string, ioJsonPath string) error {
	err := updateIOState(ioJsonPath, index, "processing")
	if err != nil {
		return fmt.Errorf("error updating IO state: %w", err)
	}

	outputDirPath := filepath.Join(jobDir, fmt.Sprintf("shard%d/outputs", index))

	err = os.MkdirAll(outputDirPath, 0755)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err)
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Call readToolConfig to get the toolConfig
	toolConfig, err := readToolConfig(ioEntry.Tool)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err)
		return fmt.Errorf("error reading tool config: %w", err)
	}

	dockerCmd, err := toolToDockerCmd(toolConfig, ioEntry, outputDirPath)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err)
		return fmt.Errorf("error converting tool to Docker cmd: %w", err)
	}

	err = runDockerCmd(dockerCmd)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err)
		return fmt.Errorf("error running Docker cmd: %w", err)
	}

	updateIOWithResult(ioJsonPath, toolConfig, index, outputDirPath)

	return nil
}
