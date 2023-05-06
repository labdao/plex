package ipwl

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/labdao/plex/internal/bacalhau"
	"github.com/labdao/plex/internal/ipfs"
)

var errOutputPathEmpty = errors.New("output file path is empty, still waiting")

func ProcessIOList(jobDir, ioJsonPath string, retry, verbose, local bool, maxConcurrency int) {
	// Use a buffered channel as a semaphore to limit the number of concurrent tasks
	semaphore := make(chan struct{}, maxConcurrency)

	// Mutex to synchronize file access
	var fileMutex sync.Mutex

	if retry {
		setRetryState(ioJsonPath)
	}

	for {
		var pendingIOs []int
		ioList, err := ReadIOList(ioJsonPath)
		if err != nil {
			fmt.Printf("Error reading IO list: %v\n", err)
			return
		}

		for i, ioEntry := range ioList {
			if ioEntry.State == "" || ioEntry.State == "retry" || ioEntry.State == "created" || ioEntry.State == "processing" || ioEntry.State == "waiting" {
				pendingIOs = append(pendingIOs, i)
			}
		}

		if len(pendingIOs) == 0 {
			break
		}

		var wg sync.WaitGroup
		for _, i := range pendingIOs {
			ioEntry := ioList[i]
			wg.Add(1)
			go func(index int, entry IO) {
				defer wg.Done()

				// Acquire the semaphore
				semaphore <- struct{}{}

				fmt.Printf("Starting to process IO entry %d \n", index)

				// add retry and resume check
				err := processIOTask(entry, index, jobDir, ioJsonPath, retry, verbose, local, &fileMutex)
				if err != nil {
					fmt.Printf("Error processing IO entry %d \n", index)
					fmt.Println(err)
				} else if errors.Is(err, errOutputPathEmpty) {
					fmt.Printf("Success processing IO entry %d \n", index)
				} else {
					fmt.Printf("Waiting to process IO entry %d \n", index)
				}

				// Release the semaphore
				<-semaphore
			}(i, ioEntry)
		}

		// Wait for all goroutines to finish
		wg.Wait()

		// Wait before re-checking chain dependecies
		time.Sleep(500 * time.Millisecond)
	}
}

func processIOTask(ioEntry IO, index int, jobDir, ioJsonPath string, retry, verbose, local bool, fileMutex *sync.Mutex) error {
	fileMutex.Lock()
	ioGraph, err := ReadIOList(ioJsonPath)
	fileMutex.Unlock()

	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error reading IO graph: %w", err)
	}

	dependsReady, err := checkSubgraphDepends(ioEntry, ioGraph)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error updating IO state: %w", err)
	} else if !dependsReady {
		err := updateIOState(ioJsonPath, index, "waiting", fileMutex)
		if err != nil {
			updateIOWithError(ioJsonPath, index, err, fileMutex)
			return fmt.Errorf("error updating IO state: %w", err)
		}
		fmt.Printf("IO Subgraph at %d is still waiting on inputs to complete \n", index)
		return errOutputPathEmpty
	}

	err = updateIOState(ioJsonPath, index, "processing", fileMutex)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
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

	err = copyInputFilesToDir(ioEntry, ioGraph, inputsDirPath)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error copying files to results input directory: %w", err)
	}

	if local {
		dockerCmd, err := toolToDockerCmd(toolConfig, ioEntry, ioGraph, inputsDirPath, outputsDirPath)
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
	} else {
		ipfsNodeUrl, err := ipfs.DeriveIpfsNodeUrl()
		if err != nil {
			updateIOWithError(ioJsonPath, index, err, fileMutex)
			return fmt.Errorf("error deriving IPFS Url: %w", err)
		}

		cid, err := ipfs.AddDirHttp(ipfsNodeUrl, inputsDirPath)
		if err != nil {
			updateIOWithError(ioJsonPath, index, err, fileMutex)
			return fmt.Errorf("error adding inputs to IPFS: %w", err)
		}

		if verbose {
			fmt.Printf("Added inputs directory to IPFS with CID: %s\n", cid)
		}

		cmd, err := toolToCmd(toolConfig, ioEntry, ioGraph)
		if err != nil {
			updateIOWithError(ioJsonPath, index, err, fileMutex)
			return fmt.Errorf("error converting tool to cmd: %w", err)
		}
		if verbose {
			fmt.Printf("Generated cmd: %s\n", cmd)
		}

		// this memory type conversion is for backwards compatibility with the -app flag
		var memory int
		if toolConfig.MemoryGB == nil {
			memory = 0
		} else {
			memory = *toolConfig.MemoryGB
		}

		bacalhauJob, err := bacalhau.CreateBacalhauJob(cid, toolConfig.DockerPull, cmd, memory, toolConfig.GpuBool, toolConfig.NetworkBool)
		if err != nil {
			updateIOWithError(ioJsonPath, index, err, fileMutex)
			return fmt.Errorf("error creating Bacalhau job: %w", err)
		}

		if verbose {
			fmt.Println("Submitting Bacalhau job")
		}
		submittedJob, err := bacalhau.SubmitBacalhauJob(bacalhauJob)
		if err != nil {
			updateIOWithError(ioJsonPath, index, err, fileMutex)
			return fmt.Errorf("error submitting Bacalhau job: %w", err)
		}

		if verbose {
			fmt.Println("Getting Bacalhau job")
		}
		results, err := bacalhau.GetBacalhauJobResults(submittedJob)
		if err != nil {
			updateIOWithError(ioJsonPath, index, err, fileMutex)
			return fmt.Errorf("error getting Bacalhau job results: %w", err)
		}

		if verbose {
			fmt.Println("Downloading Bacalhau job")
			fmt.Printf("Output dir of %s \n", outputsDirPath)
		}
		err = bacalhau.DownloadBacalhauResults(outputsDirPath, submittedJob, results)
		if err != nil {
			updateIOWithError(ioJsonPath, index, err, fileMutex)
			return fmt.Errorf("error downloading Bacalhau results: %w", err)
		}

		if verbose {
			fmt.Println("Cleaning Bacalhau job")
		}
		err = cleanBacalhauOutputDir(outputsDirPath)
		if err != nil {
			updateIOWithError(ioJsonPath, index, err, fileMutex)
			return fmt.Errorf("error cleaning Bacalhau output directory: %w", err)
		}
	}

	err = updateIOWithResult(ioJsonPath, toolConfig, index, outputsDirPath, fileMutex)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error updating IO with result: %w", err)
	}

	return nil
}

func copyInputFilesToDir(ioEntry IO, ioGraph []IO, dirPath string) error {
	// Ensure the destination directory exists
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	for _, input := range ioEntry.Inputs {
		srcPath, err := DetermineSrcPath(input, ioGraph)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dirPath, filepath.Base(srcPath))

		err = copyFile(srcPath, destPath)
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

func cleanBacalhauOutputDir(outputsDirPath string) error {
	bacalOutputsDirPath := filepath.Join(outputsDirPath, "outputs")

	// Move files from /outputs to outputsDirPath
	files, err := ioutil.ReadDir(bacalOutputsDirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		src := filepath.Join(bacalOutputsDirPath, file.Name())
		dst := filepath.Join(outputsDirPath, file.Name())
		fmt.Printf("Moving %s to %s", src, dst)
		if err := os.Rename(src, dst); err != nil {
			return err
		}
	}

	// if err := os.RemoveAll(bacalOutputsDirPath); err != nil {
	// 	return err
	// }

	return nil
}

func checkSubgraphDepends(ioEntry IO, ioGraph []IO) (bool, error) {
	dependsReady := true

	for _, input := range ioEntry.Inputs {
		_, err := DetermineSrcPath(input, ioGraph)
		if err != nil {
			if errors.Is(err, errOutputPathEmpty) {
				dependsReady = false
				break
			}
			return false, fmt.Errorf("failed to determine source path: %w", err)
		}
	}

	return dependsReady, nil
}

func DetermineSrcPath(input FileInput, ioGraph []IO) (string, error) {
	// Check if the input.FilePath has the format ${i[key]}
	re := regexp.MustCompile(`^\$\{(\d+)\[(\w+)\]\}$`)
	match := re.FindStringSubmatch(input.FilePath)

	if match == nil {
		// The input.FilePath is a normal file path
		return input.FilePath, nil
	}

	// Extract the index and key from the matched pattern
	indexStr, key := match[1], match[2]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", fmt.Errorf("invalid index in input.FilePath: %s", input.FilePath)
	}

	if index < 0 || index >= len(ioGraph) {
		return "", fmt.Errorf("index out of range: %d", index)
	}

	// Check that dependent subgraph has not failed
	if ioGraph[index].State == "failed" {
		return "", fmt.Errorf("dependent subgraph %d failed", index)
	}

	// Get the output Filepath of the corresponding key
	output, ok := ioGraph[index].Outputs[key]
	if !ok {
		return "", fmt.Errorf("key not found in outputs: %s", key)
	}

	outputFilepath := ""
	switch output := output.(type) {
	case FileOutput:
		outputFilepath = output.FilePath
	case ArrayFileOutput:
		return "", fmt.Errorf("PLEx does not currently support ArrayFileOutput as an input")
	default:
		return "", fmt.Errorf("unknown output type")
	}

	if outputFilepath == "" {
		return "", errOutputPathEmpty
	}

	return outputFilepath, nil
}

func setRetryState(ioJsonPath string) error {
	// Read the IO list from the file
	ioList, err := ReadIOList(ioJsonPath)
	if err != nil {
		return fmt.Errorf("failed to read IO list: %w", err)
	}

	// Iterate through the IO list and update the state of failed entries to "retry"
	for i, ioEntry := range ioList {
		if ioEntry.State == "failed" {
			ioList[i].State = "retry"
		}
	}

	// Write the updated IO list back to the file
	err = WriteIOList(ioJsonPath, ioList)
	if err != nil {
		return fmt.Errorf("failed to write updated IO list: %w", err)
	}

	return nil
}
