package ipwl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/labdao/plex/internal/ipfs"
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

func updateIOWithResult(ioJsonPath string, toolConfig Tool, index int, outputDirPath string, fileMutex *sync.Mutex) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	ioList, err := ReadIOList(ioJsonPath)
	if err != nil {
		return fmt.Errorf("error reading IO list: %w", err)
	}

	var outputsWithNoData []string

	for outputKey, output := range toolConfig.Outputs {
		matchingFiles, err := findMatchingFilesForPatterns(outputDirPath, output.Glob)
		if err != nil {
			return fmt.Errorf("error matching output files: %w", err)
		}

		if len(matchingFiles) == 0 {
			outputsWithNoData = append(outputsWithNoData, outputKey)
			continue
		}

		if output.Type == "File" {
			filePath := matchingFiles[0]

			// Update IO entry
			cid, err := ipfs.GetFileCid(filePath)
			if err != nil {
				return fmt.Errorf("error generating file IPFS cid: %w", err)
			}

			ioList[index].Outputs[outputKey] = FileOutput{
				Class:    "File",
				FilePath: filepath.Base(filePath),
				IPFS:     cid,
			}
		} else if output.Type == "Array" && output.Item == "File" {
			var files []FileOutput
			for _, filePath := range matchingFiles {
				cid, err := ipfs.GetFileCid(filePath)
				if err != nil {
					return fmt.Errorf("error generating file IPFS cid: %w", err)
				}

				files = append(files, FileOutput{
					Class:    "File",
					FilePath: filepath.Base(filePath),
					IPFS:     cid,
				})
			}

			ioList[index].Outputs[outputKey] = ArrayFileOutput{
				Class: "Array",
				Files: files,
			}
		} else {
			return fmt.Errorf("unsupported output Type and Item combination: Type=%s, Item=%s", output.Type, output.Item)
		}
	}

	if len(outputsWithNoData) > 0 {
		ioList[index].State = "failed"
	} else {
		ioList[index].State = "completed"
	}

	err = WriteIOList(ioJsonPath, ioList)
	if err != nil {
		return fmt.Errorf("error writing updated IO list: %w", err)
	}

	if len(outputsWithNoData) > 0 {
		return fmt.Errorf("no output data found for: %v", outputsWithNoData)
	}

	return nil
}
