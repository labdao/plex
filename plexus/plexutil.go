package main

import (
	"fmt"
	"flag"
	"os"
	"bufio"
	"encoding/json"
	//ipfsapi "github.com/ipfs/go-ipfs-api"
	//"io/ioutil"
)

// appStruct for the application file in the app.jsonl file
type appStruct struct {
	App     string `json:"app"`
	Inputs  []struct {
		protein_path string `json:"protein_path"`
		ligand      string `json:"ligand"`
	} `json:"inputs"`
	Outputs []string `json:"outputs"`
}

// validate that the directory path exists and is a directory
func validateDirectoryPath(directory *string){
	if _, err := os.Stat(*directory); os.IsNotExist(err) {
		fmt.Println("Error: the directory path does not exist.")
		os.Exit(1)
	}
	if fileInfo, err := os.Stat(*directory); err == nil && !fileInfo.Mode().IsDir() {
		fmt.Println("Error: the path provided is not a directory.")
		os.Exit(1)
	}
	fmt.Println("Directory found:", *directory)
}

// validate that the application is supported based on the app.jsonl file
func validateApplication(application *string){	
	file, err := os.Open("plexus/app.jsonl")
	if err != nil {
		fmt.Println("Error opening application file:", err)
		return
	}
	defer file.Close()

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var appData appStruct
		err = json.Unmarshal([]byte(scanner.Text()), &appData)
		if err != nil {
			fmt.Println("Error unmarshalling application file JSON:", err)
			return
		}
		if appData.App == *application {
			fmt.Println("Application found:", appData.App)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}
}

// index the directory path and return the files that match the input of the specified application
func indexDirectoryPath(directory *string, layers int) {
	// read the app.jsonl file
	file, err := os.Open("plexus/app.jsonl")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// read the file line by line
	var appData appStruct
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = json.Unmarshal([]byte(scanner.Text()), &appData)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		break
	}
	// additional errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}
	if _, err := os.Stat(*directory); os.IsNotExist(err) {
		fmt.Println("Error: the directory path does not exist.")
		os.Exit(1)
	}

	// walk the directory path and return the files that match the input file extensions of the specified application
	var files []string
	err = filepath.Walk(*directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if layers > 0 {
				layers--
				return nil
			} else {
				return filepath.SkipDir
			}
		}
		for _, input := range appData.Inputs {
			for ext := range input {
				if strings.HasSuffix(path, input[ext]) {
					files = append(files, path)
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Files:", files)
}


func main() {
	// define the flags
	app := flag.String("application", "", "Application name")
	dir := flag.String("directory", "", "Directory path")
	flag.Parse()


	// print the values of the flags
	fmt.Println("Provided application name:", *app)
	fmt.Println("Provided directory path:", *dir)

	// validate the flags
	validateDirectoryPath(dir)
	validateApplication(app)

	// creating index file
	indexDirectoryPath(dir, 1)

	

	
}
