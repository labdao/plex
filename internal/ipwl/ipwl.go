package ipwl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func ProcessIOList(ioList []IO, jobDir, ioJsonPath string, verbose bool, maxConcurrency int) {
	// Use a buffered channel as a semaphore to limit the number of concurrent tasks
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	// Mutex to synchronize file access
	var fileMutex sync.Mutex

	for i, ioEntry := range ioList {
		wg.Add(1)
		go func(index int, entry IO) {
			defer wg.Done()

			// Acquire the semaphore
			semaphore <- struct{}{}

			fmt.Printf("Starting to process IO entry %d \n", index)
			err := processIOTask(entry, index, jobDir, ioJsonPath, verbose, &fileMutex)
			if err != nil {
				fmt.Printf("Error processing IO entry %d \n", index)
			} else {
				fmt.Printf("Success processing IO entry %d \n", index)
			}

			// Release the semaphore
			<-semaphore
		}(i, ioEntry)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func processIOTask(ioEntry IO, index int, jobDir, ioJsonPath string, verbose bool, fileMutex *sync.Mutex) error {
	err := updateIOState(ioJsonPath, index, "processing", fileMutex)
	if err != nil {
		return fmt.Errorf("error updating IO state: %w", err)
	}

	workingDirPath := filepath.Join(jobDir, fmt.Sprintf("entry-%d", index))
	err = os.MkdirAll(workingDirPath, 0755)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error creating working directory: %w", err)
	}

	outputsDirPath := filepath.Join(workingDirPath, "outputs")
	err = os.MkdirAll(outputsDirPath, 0755)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error creating outputs directory: %w", err)
	}

	inputsDirPath := filepath.Join(workingDirPath, "inputs")
	err = os.MkdirAll(inputsDirPath, 0755)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error creating output directory: %w", err)
	}

	toolConfig, err := ReadToolConfig(ioEntry.Tool)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error reading tool config: %w", err)
	}

	err = copyInputFilesToDir(ioEntry, inputsDirPath)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error copying files to results input directory: %w", err)
	}

	dockerCmd, err := toolToDockerCmd(toolConfig, ioEntry, inputsDirPath, outputsDirPath)
	if verbose {
		fmt.Println("Generated docker cmd:", dockerCmd)
	}
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error converting tool to Docker cmd: %w", err)
	}

	output, err := runDockerCmd(dockerCmd)
	if verbose {
		fmt.Printf("Docker ran with output: %s \n", output)
	}
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error running Docker cmd: %w", err)
	}

	err = updateIOWithResult(ioJsonPath, toolConfig, index, outputsDirPath, fileMutex)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error updating IO with result: %w", err)
	}

	return nil
}

func copyInputFilesToDir(ioEntry IO, dirPath string) error {
	// Ensure the destination directory exists
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	for _, input := range ioEntry.Inputs {
		srcPath := input.FilePath
		destPath := filepath.Join(dirPath, filepath.Base(srcPath))

		err := copyFile(srcPath, destPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
