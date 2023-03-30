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

/*
<plex-uuid>
	io.json // slowly be updated
	/shard_0 7n9g.pdb & ZINC000003986735.sdf
		/outputs <- where the docker mount dumps all files
			/7n9g.pdb
			/ZINC000003986735-docked.sdf
	/shard_1 7n9g.pdb & ZINC000019632618.sdf
		/outputs
			/7n9g.pdb
			/ZINC000019632618-docked.sdf
*/

/*
func processIOTask(ioEntry IO, ioJsonPath string, index int, jobDir string) error {
	toolConfig, err := readToolConfig(ioEntry.Tool)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err)
		return fmt.Errorf("error reading tool config: %w", err)
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
