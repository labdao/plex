package plex

import (
	"fmt"
	"os"

	"github.com/labdao/plex/internal/bacalhau"
)

func Execute(app, inputDir, gpu, appConfigsFilePath string, layers int) {
	// validate the flags
	fmt.Println("## Validating ##")
	appConfig, err := FindAppConfig(app, appConfigsFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// creating index file
	fmt.Println("## Searching input files ##")
	identifiedFiles, err := searchDirectoryPath(inputDir, appConfig, layers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Found", len(identifiedFiles), "matching files")
	for _, fileName := range identifiedFiles {
		fmt.Println(fileName)
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, movedFiles, jobDir, err := createInputsDirectory(dir, identifiedFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Created job directory", jobDir)

	createIndex(movedFiles, appConfig, jobDir)

	// create instructions
	instruction, err := CreateInstruction(app, "config/instruction_template.jsonl", jobDir, map[string]string{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := bacalhau.InstructionToBacalhauCmd(instruction.InputCIDs[0], instruction.Container, instruction.Cmd, "true")
	fmt.Println(cmd)

	/*
		// create bacalhau job
		fmt.Println("## Creating Bacalhau Job ##")
		job, err := bacalhau.CreateBacalhauJob(instruction.InputCIDs[0], instruction.Container, instruction.Cmd, gpu)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		submittedJob, err := bacalhau.SubmitBacalhauJob(job)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Bacalhau Job Id: " + submittedJob.Metadata.ID)
		results, err := bacalhau.GetBacalhauJobResults(submittedJob)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bacalhau.DownloadBacalhauResults(jobDir, submittedJob, results)
		fmt.Println("Your job results have been downloaded to " + jobDir)
	*/
}
