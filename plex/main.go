package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// required flags
	app := flag.String("app", "", "Application name")
	inputDir := flag.String("input-dir", "", "Input directory path")
	// todo: needs to be a boolean flag
	gpu := flag.String("gpu", "true", "Use GPU")

	// optional flags
	appConfigsFilePath := flag.String("app-configs", "app.jsonl", "App Configurations file")
	layers := flag.Int("layers", 2, "number of layers to search in the directory path")
	flag.Parse()

	// print the values of the flags
	fmt.Println("## User input ##")
	fmt.Println("Provided application name:", *app)
	fmt.Println("Provided directory path:", *inputDir)
	fmt.Println("Using GPU:", *gpu)

	fmt.Println("## Default parameters ##")
	fmt.Println("Using app configs:", *appConfigsFilePath)
	fmt.Println("Setting layers to:", *layers)

	// validate the flags
	fmt.Println("## Validating ##")
	appConfig, err := findAppConfig(*app, *appConfigsFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// creating index file
	fmt.Println("## Seaching input files ##")
	identifiedFiles, err := searchDirectoryPath(inputDir, appConfig, *layers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Found", len(identifiedFiles), "matching files")
	for _, fileName := range identifiedFiles {
		fmt.Println(fileName)
	}

	fmt.Println("## Creating job directory ##")
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
	fmt.Println("## Creating index ##")
	createIndex(movedFiles, appConfig, jobDir)

	// create instructions
	fmt.Println("## Creating instruction ##")
	instruction, err := CreateInstruction(*app, "instruction_template.jsonl", jobDir, map[string]string{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bacalhauCmd := InstructionToBacalhauCmd(instruction.InputCIDs[0], instruction.Container, instruction.Cmd, *gpu)
	fmt.Println(bacalhauCmd)
}
