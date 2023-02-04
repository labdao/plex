package main

import (
	"fmt"
	"flag"
	"os"
	"bufio"
	"encoding/json"
	"strings"
	"path/filepath"
	"github.com/google/uuid"
	fileutils "github.com/docker/docker/pkg/fileutils"
	//ipfsapi "github.com/ipfs/go-ipfs-api"
	//"io/ioutil"
)

// appStruct for the application file in the app.jsonl file
//TODO #65 #64
type appStruct struct {
	App     string `json:"app"`
	Inputs  [][2]string `json:"inputs"`
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

// validate the app.jsonl file
func validateAppConfig(app_config *string){
	file, err := os.Open(*app_config)
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
	if _, err := os.Stat(*app_config); os.IsNotExist(err) {
		fmt.Println("Error: the directory path does not exist.")
		os.Exit(1)
	}
	//fmt.Println("App Config found:", *app_config)
}

// validate that the application is supported based on the app.jsonl file
func validateApplication(application *string, app_config *string){
	// validate config file
	validateAppConfig(app_config)
	file, err := os.Open(*app_config)

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
func indexSearchDirectoryPath(directory *string, app_config *string, layers int) []string {
	// validate config file
	validateAppConfig(app_config)
	// read the app.jsonl file
	file, err := os.Open(*app_config)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// read the file line by line
	var appData appStruct
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = json.Unmarshal([]byte(scanner.Text()), &appData)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return nil
		}
		break
	}
	// additional errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return nil
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
		//TODO: create safeguard to prevent "ligand" or "protein_path" from being added to the index
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

	for _, file := range files {
		fmt.Println(file)
	}
	return files
}

func indexCreateInputsVolume(volume_directory *string, files []string) {
	// creating uuid for the volume
	id := uuid.New()
	//create a volume directory
	volume_path := *volume_directory + "/" + id.String()
	err := os.Mkdir(volume_path, 0755)
	if err != nil {
		fmt.Println("Error creating volume directory:", err)
		return
	}
	// copy the files to the volume directory
	for _, file := range files {
		_, err = fileutils.CopyFile(file, volume_path)
		if err != nil {
			fmt.Println("Error copying file to volume directory:", err)
			return
		}
	}
	print("Volume created:", volume_path)
}

// create a csv file that lists the indexed files in an application-specific format
//func indexCreateIndexCSV(files, app_config *string, layers int) []string {
//}

func main() {
	// define the flags
	app := flag.String("application", "", "Application name")
	in_dir := flag.String("input_directory", "", "Input directory path")
	// additional flags
	app_config := flag.String("app_config", "app.jsonl", "App Config file")
	layers := flag.Int("layers", 2, "number of layers to search in the directory path")
	flag.Parse()


	// print the values of the flags
	fmt.Println("## User input ##")
	fmt.Println("Provided application name:", *app)
	fmt.Println("Provided directory path:", *in_dir)
	fmt.Println("## Default parameters ##")
	fmt.Println("Using app config:", *app_config)
	fmt.Println("Setting layers to:", *layers)

	// validate the flags
	fmt.Println("## Validating ##")
	validateApplication(app, app_config)
	validateDirectoryPath(in_dir)
	validateAppConfig(app_config)
	fmt.Println("App Config found:", *app_config)
	
	

	// creating index file
	fmt.Println("## Creating index ##")
	out := indexSearchDirectoryPath(in_dir, app_config, *layers)
	// TODO pass separate volume directory
	indexCreateInputsVolume(in_dir, out)
	//indexCreateIndexCSV(out, app_config)
	// TODO create indexCreateIndexJSONL(out, app_config)
	

	
}
