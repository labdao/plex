package ipwl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func updateIOWithError(ioJsonPath string, index int, err error, fileMutex *sync.Mutex) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	ioList, errRead := ReadIOList(ioJsonPath)
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

func updateIOWithBacalhauJob(ioJsonPath string, index int, bacalhauJobId string, fileMutex *sync.Mutex) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	ioList, errRead := ReadIOList(ioJsonPath)
	if errRead != nil {
		return fmt.Errorf("failed to read IO list: %w", errRead)
	}

	if index < 0 || index >= len(ioList) {
		return fmt.Errorf("index out of range: %d", index)
	}

	ioList[index].BacalhauJobId = bacalhauJobId

	errWrite := WriteIOList(ioJsonPath, ioList)
	if errWrite != nil {
		return fmt.Errorf("failed to write updated IO list: %w", errWrite)
	}

	return nil
}

func updateIOState(ioJsonPath string, index int, state string, fileMutex *sync.Mutex) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	ioList, err := ReadIOList(ioJsonPath)
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

func findMatchingFilesForPatterns(outputDirPath string, patterns []string) ([]string, error) {
	var matchingFiles []string

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(outputDirPath, pattern))
		if err != nil {
			return nil, fmt.Errorf("error while matching pattern: %w", err)
		}

		matchingFiles = append(matchingFiles, matches...)
	}

	return matchingFiles, nil
}
