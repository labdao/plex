package main

import (
	"flag"
	"fmt"
	"os"
)

// go run main.go --app diffdock --input-dir /path/to/pdbs (mode, network and other flags in future)

func main() {
	// define the flags
	app := flag.String("app", "", "Application name")
	inputDir := flag.String("input-dir", "", "Input directory path")

	// optional flags
	appConfigsFilePath := flag.String("app-configs", "app.jsonl", "App Configurations file")
	layers := flag.Int("layers", 2, "number of layers to search in the directory path")
	flag.Parse()

	// print the values of the flags
	fmt.Println("## User input ##")
	fmt.Println("Provided application name:", *app)
	fmt.Println("Provided directory path:", *inputDir)

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
	// ValidateDirectoryPath(*inputDir)
	// ValidateAppConfig(*appConfig)

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

	// TODO enable passing an array of multiple input directories
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
	createIndex(movedFiles, "app.jsonl", jobDir)

	// create instructions
	fmt.Println("## Creating instruction ##")
	instruction, err := CreateInstruction(*app, "instruction_template.jsonl", jobDir, map[string]string{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bacalhauCmd := InstructionToBacalhauCmd(instruction.InputCIDs[0], instruction.Container, instruction.Cmd)
	fmt.Println(bacalhauCmd)
}

/*
func main2() {
	client, err := w3s.NewClient(
		w3s.WithEndpoint("https://api.web3.storage"),
		w3s.WithToken(os.Getenv("WEB3STORAGE_TOKEN")),
	)
	errorCheck(err)

	if len(os.Args) < 2 {
		fmt.Println("Error: Please specify a command (putFile, putDirectory, getFiles)")
		os.Exit(1)
	}

	command := strings.ToLower(os.Args[1])

	switch command {
	case "putfile":
		if len(os.Args) != 3 {
			fmt.Println("Error: Please specify a file path")
			os.Exit(1)
		}
		filePath := os.Args[2]
		file, err := os.Open(filePath)
		errorCheck(err)
		defer file.Close()
		putFile(client, file)
	case "putdirectory":
		if len(os.Args) != 3 {
			fmt.Println("Error: Please specify a directory path")
			os.Exit(1)
		}
		directoryPath := os.Args[2]
		putDirectory(client, directoryPath)
	case "getfiles":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please specify a CID")
			os.Exit(1)
		}
		cidString := os.Args[2]
		getFiles(client, cidString)
	}
}
*/
