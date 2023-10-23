package ipwl

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labdao/plex/internal/bacalhau"
	"github.com/labdao/plex/internal/ipfs"
)

var errOutputPathEmpty = errors.New("output file path is empty, still waiting")

func RunIO(ioJsonCid, outputDir, selector string, verbose, showAnimation bool, maxTime int, annotations []string) (completedIoJsonCid, ioJsonPath string, err error) {
	id := uuid.New()
	var cwd string
	if outputDir != "" {
		absPath, err := filepath.Abs(outputDir)
		if err != nil {
			return completedIoJsonCid, ioJsonPath, err
		}
		cwd = absPath
	} else {
		cwd, err = os.Getwd()
		if err != nil {
			return completedIoJsonCid, ioJsonPath, err
		}
		cwd = path.Join(cwd, "jobs")
	}
	workDirPath := path.Join(cwd, id.String())
	err = os.MkdirAll(workDirPath, 0755)
	if err != nil {
		return completedIoJsonCid, ioJsonPath, err
	}
	fmt.Println("Created working directory:", workDirPath)

	ioJsonPath = path.Join(workDirPath, "io.json")
	err = ipfs.DownloadFileContents(ioJsonCid, ioJsonPath)
	if err != nil {
		return completedIoJsonCid, ioJsonPath, err
	}
	fmt.Println("Initialized IO file at:", ioJsonPath)

	fmt.Println("Extracting user ID from IO JSON...")
	userID, err := ExtractUserIDFromIOJson(ioJsonPath)
	if err != nil {
		return completedIoJsonCid, ioJsonPath, err
	}

	if userID != "" && !ContainsUserIdAnnotation(annotations) {
		annotations = append(annotations, fmt.Sprintf("userId=%s", userID))
	}

	if maxTime > 60 {
		fmt.Println("Error: maxTime cannot exceed 60 minutes")
		os.Exit(1)
	}

	retry := false
	fmt.Println("Processing IO Entries")
	ProcessIOList(workDirPath, ioJsonPath, selector, retry, verbose, showAnimation, maxTime, annotations)
	fmt.Printf("Finished processing, results written to %s\n", ioJsonPath)
	completedIoJsonCid, err = ipfs.PinFile(ioJsonPath)
	if err != nil {
		return completedIoJsonCid, ioJsonPath, err
	}

	fmt.Println("Completed IO JSON CID:", completedIoJsonCid)
	return completedIoJsonCid, ioJsonPath, nil
}

func SubmitIoList(ioList []IO, selector string, maxTime int, annotations []string) []IO {
	submittedIOList := make([]IO, len(ioList))
	for i, ioEntry := range ioList {
		log.Printf("Submitting IO entry %d \n", i)
		submittedIOList[i] = ioEntry
		log.Println("Reading tool config")
		toolConfig, _, err := ReadToolConfig(ioEntry.Tool.IPFS)
		if err != nil {
			submittedIOList[i].State = "failed"
			submittedIOList[i].ErrMsg = fmt.Sprintf("error reading tool config: %v", err)
			continue
		}
		log.Println("Creating cmd")
		cmd, err := toolToCmd2(toolConfig, ioEntry, ioList)
		if err != nil {
			submittedIOList[i].State = "failed"
			submittedIOList[i].ErrMsg = fmt.Sprintf("error reading tool config: %v", err)
			continue
		}
		log.Printf("cmd: %s \n", cmd)
		log.Println("mapping inputs")
		bacalhauInputs := make(map[string]string)
		for key, input := range ioEntry.Inputs {
			bacalhauInputs[key] = input.IPFS
		}
		log.Println("creating bacalhau job")
		// this memory type conversion is for backwards compatibility with the -app flag
		var memory int
		if toolConfig.MemoryGB == nil {
			memory = 0
		} else {
			memory = *toolConfig.MemoryGB
		}
		log.Println("creating bacalhau job")
		bacalhauJob, err := bacalhau.CreateBacalhauJobV2(bacalhauInputs, toolConfig.DockerPull, selector, cmd, maxTime, memory, toolConfig.GpuBool, toolConfig.NetworkBool, annotations)
		if err != nil {
			submittedIOList[i].State = "failed"
			submittedIOList[i].ErrMsg = fmt.Sprintf("error creating Bacalhau job: %v", err)
			continue
		}

		log.Println("submitting bacalhau job")
		submittedJob, err := bacalhau.SubmitBacalhauJob(bacalhauJob)
		if err != nil {
			submittedIOList[i].State = "failed"
			submittedIOList[i].ErrMsg = fmt.Sprintf("error submitting Bacalhau job: %v", err)
			continue
		}
		submittedIOList[i].State = "new"
		submittedIOList[i].BacalhauJobId = submittedJob.Metadata.ID
	}
	log.Println("returning io submited list")
	return submittedIOList
}

func ProcessIOList(jobDir, ioJsonPath, selector string, retry, verbose, showAnimation bool, maxTime int, annotations []string) {
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

				fmt.Printf("Starting to process IO entry %d \n", index)

				// add retry and resume check
				err := processIOTask(entry, index, maxTime, jobDir, ioJsonPath, selector, retry, verbose, showAnimation, annotations, &fileMutex)
				if errors.Is(err, errOutputPathEmpty) {
					fmt.Printf("Waiting to process IO entry %d \n", index)
				} else if err != nil {
					fmt.Printf("Error processing IO entry %d \n", index)
					fmt.Println(err)
				} else {
					fmt.Printf("Success processing IO entry %d \n", index)
				}
			}(i, ioEntry)
		}

		// Wait for all goroutines to finish
		wg.Wait()

		// Wait before re-checking chain dependencies
		// Todo: switch to event driven checking so go-routine
		// can stop once bacalhau job is created instead of completed
		time.Sleep(500 * time.Millisecond)
	}
}

func processIOTask(ioEntry IO, index, maxTime int, jobDir, ioJsonPath, selector string, retry, verbose, showAnimation bool, annotations []string, fileMutex *sync.Mutex) error {
	fileMutex.Lock()
	ioGraph, err := ReadIOList(ioJsonPath)
	fileMutex.Unlock()

	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error reading IO graph: %w", err)
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

	toolConfig, _, err := ReadToolConfig(ioEntry.Tool.IPFS)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error reading tool config: %w", err)
	}

	err = downloadInputFilesToDir(ioEntry, ioGraph, inputsDirPath)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error copying files to results input directory: %w", err)
	}

	cid, err := ipfs.PinDir(inputsDirPath)
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

	bacalhauJob, err := bacalhau.CreateBacalhauJob(cid, toolConfig.DockerPull, cmd, selector, maxTime, memory, toolConfig.GpuBool, toolConfig.NetworkBool, annotations)
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
	} else {
		err = updateIOWithBacalhauJob(ioJsonPath, index, submittedJob.Metadata.ID, fileMutex)
		if err != nil {
			updateIOWithError(ioJsonPath, index, err, fileMutex)
			return fmt.Errorf("error updating IO with Bacalhau job: %w", err)
		}
	}

	if verbose {
		fmt.Println("Getting Bacalhau job")
	}
	results, err := bacalhau.GetBacalhauJobResults(submittedJob, showAnimation, maxTime)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error getting Bacalhau job results: %w", err)
	}

	if verbose {
		fmt.Println("Downloading Bacalhau job")
		fmt.Printf("Output dir of %s \n", outputsDirPath)
	}

	for _, result := range results {
		fmt.Printf("Downloading result %s to %s \n", result.Data.CID, outputsDirPath)
		err = ipfs.DownloadToDirectory(result.Data.CID, outputsDirPath)
	}

	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error downloading Bacalhau results: %w", err)
	}

	if verbose {
		fmt.Println("Cleaning Bacalhau job")
	}
	err = cleanBacalhauOutputDir(outputsDirPath, verbose)

	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error cleaning Bacalhau output directory: %w", err)
	}

	err = updateIOWithResult(ioJsonPath, toolConfig, index, outputsDirPath, fileMutex)
	if err != nil {
		updateIOWithError(ioJsonPath, index, err, fileMutex)
		return fmt.Errorf("error updating IO with result: %w", err)
	}

	return nil
}

func downloadInputFilesToDir(ioEntry IO, ioGraph []IO, dirPath string) error {
	// Ensure the destination directory exists
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	for _, input := range ioEntry.Inputs {
		destPath := filepath.Join(dirPath, input.FilePath)
		cidPath := input.IPFS + "/" + input.FilePath
		err = ipfs.DownloadFileContents(cidPath, destPath)
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

func cleanBacalhauOutputDir(outputsDirPath string, verbose bool) error {
	bacalOutputsDirPath := filepath.Join(outputsDirPath, "outputs")

	// Move files from /outputs to outputsDirPath
	files, err := ioutil.ReadDir(bacalOutputsDirPath)
	if err != nil {
		return fmt.Errorf("error reading Bacalhau output directory: %w", err)
	}

	for _, file := range files {
		src := filepath.Join(bacalOutputsDirPath, file.Name())
		dst := filepath.Join(outputsDirPath, file.Name())

		if verbose {
			fmt.Printf("Moving %s to %s\n", src, dst)
		}
		if err := os.Rename(src, dst); err != nil {
			return fmt.Errorf("error moving file from %s to %s: %w", src, dst, err)
		}
	}

	// remove empty outputs folder now that files have been moved
	err = os.Remove(bacalOutputsDirPath)
	if err != nil {
		return fmt.Errorf("error removing Bacalhau output directory: %w", err)
	}

	return nil
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
