package ipwl

import (
	"encoding/json"
	"fmt"
	"os"
)

func updateIOState(ioJsonPath string, index int, state string) error {
	ioList, err := readIOList(ioJsonPath)
	if err != nil {
		return fmt.Errorf("error reading IO list: %w", err)
	}

	if index >= len(ioList) {
		return fmt.Errorf("index out of range: %d", index)
	}

	ioList[index].State = state

	file, err := os.OpenFile(ioJsonPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening IO JSON file for writing: %w", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(ioList)
	if err != nil {
		return fmt.Errorf("error encoding updated IO list: %w", err)
	}

	return nil
}

/*
func processIOList(ioList []IO, jobDir string) error {
	for i, ioEntry := range ioList {
		err := processIOTask(ioEntry, i, jobDir)
		if err != nil {
			return fmt.Errorf("error processing IO task at index %d: %w", i, err)
		}
	}

	return nil
}

func processIOTask(ioEntry IO, index int, jobDir string, ioJsonPath string, state string) error {
	err := updateIOState(ioJsonPath, index, state)
	if err != nil {
		return fmt.Errorf("error updating IO state: %w", err)
	}

	outputDirPath := filepath.Join(jobDir, fmt.Sprintf("shard%d/outputs", index))

	err = os.MkdirAll(outputDirPath, 0755)
	if err != nil {
		updateIOWithError(index, err)
		return fmt.Errorf("error creating output directory: %w", err)
	}

	dockerCmd, err := toolToDockerCmd(toolConfig, ioEntry, outputDirPath)
	if err != nil {
		updateIOWithError(index, err)
		return fmt.Errorf("error converting tool to Docker cmd: %w", err)
	}

	err = runDockerCmd(dockerCmd)
	if err != nil {
		updateIOWithError(index, err)
		return fmt.Errorf("error running Docker cmd: %w", err)
	}

	updateIOWithResult(index, toolConfig, outputDirPath)

	return nil
}
*/
