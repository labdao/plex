package ipwl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func updateIOWithError(ioJsonPath string, index int, err error) error {
	ioList, errRead := readIOList(ioJsonPath)
	if errRead != nil {
		return fmt.Errorf("failed to read IO list: %w", errRead)
	}

	if index < 0 || index >= len(ioList) {
		return fmt.Errorf("index out of range: %d", index)
	}

	ioList[index].State = "failed"
	ioList[index].ErrMsg = err.Error()

	errWrite := WriteIOList(ioJsonPath, ioList)
	if errWrite != nil {
		return fmt.Errorf("failed to write updated IO list: %w", errWrite)
	}

	return nil
}

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

func updateIOWithResult(ioJsonPath string, toolConfig Tool, index int, outputDirPath string) error {
	ioList, err := readIOList(ioJsonPath)
	if err != nil {
		return fmt.Errorf("error reading IO list: %w", err)
	}

	for outputKey, output := range toolConfig.Outputs {
		if output.Type != "File" {
			continue
		}

		matchingFiles := findMatchingFilesForPatterns(outputDirPath, output.Glob)
		if len(matchingFiles) == 0 {
			continue
		}

		// Assume there is only one matching file per output key
		filePath := matchingFiles[0]
		filename := filepath.Base(filePath)

		// Update IO entry
		ioList[index].Outputs[outputKey] = FileOutput{
			Class:    "File",
			FilePath: filePath,
			Basename: filename,
		}
	}

	ioList[index].State = "completed"

	err = WriteIOList(ioJsonPath, ioList)
	if err != nil {
		return fmt.Errorf("error writing updated IO list: %w", err)
	}

	return nil
}

func findMatchingFilesForPatterns(outputDirPath string, globPatterns []string) []string {
	var matchingFiles []string

	for _, globPattern := range globPatterns {
		patternMatches, err := filepath.Glob(filepath.Join(outputDirPath, globPattern))
		if err != nil {
			fmt.Printf("Error matching glob pattern: %v", err)
			continue
		}

		matchingFiles = append(matchingFiles, patternMatches...)
	}

	return matchingFiles
}
