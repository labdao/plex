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
	"encoding/csv"
	//"encoding/csv"
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

func indexCreateInputsVolume(volume_directory *string, files []string, prefix string) (string, []string, string) {
	// creating uuid for the volume
	id := uuid.New()
	//create a volume directory
	//TODO: add safeguard to prevent the creation of a volume directory if it already exists
	//TODO: find elegant solution for "nil"
	volume_path := *volume_directory + "/" + id.String()
	err := os.Mkdir(volume_path, 0755)
	if err != nil {
		fmt.Println("Error creating a volume directory:", err)
		return "nil", files, "nil"
	}
	os.Mkdir(volume_path + prefix, 0755)
	// copy the files to the volume directory
	new_files := make([]string, 0)
	for _, file := range files {
		_, err = fileutils.CopyFile(file, volume_path + prefix + "/" + filepath.Base(file))
		if err != nil {
			fmt.Println("Error copying file to volume directory:", err)
			return "nil", files, "nil"
		}
		new_files = append(new_files, prefix + "/" + filepath.Base(file))
	}
	print("Volume created: ", volume_path, "\n")
	return id.String(), new_files, volume_path
}

// create a csv file that lists the indexed files in an application-specific format
// the paths to the input files and the app config are given as input 
// the path to the index.csv file is returned
func indexCreateIndexJSONL(new_files []string, app_config *string, volume_path string) (string, []map[string]string) {
	// read the app.jsonl file
	file, err := os.Open(*app_config)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	
	// parse the json object
    var appData appStruct
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = json.Unmarshal([]byte(scanner.Text()), &appData)
		if err != nil {
			panic(err)
		}
		break
	}

    // map the input new_files to their respective columns based on the config
    // Create an empty 2D slice to store the columns and their corresponding files
	var result []map[string]string // Define a slice to store the maps

	for _, file := range new_files {
		m := make(map[string]string)
		for _, mapping := range appData.Inputs {
			if strings.HasSuffix(file, mapping[1]) {
				m[mapping[0]] = file
				break
			}
		}
		result = append(result, m)
	}
	//fmt.Println(result)

	// generate combinations of the mapping
	//TODO implement generalisable version
	combinations := []map[string]string{}
	for _, r := range result {
		if r["ligand"] != "" {
			for _, r2 := range result {
				if r2["protein_path"] != "" && r2["ligand"] == "" {
					m := make(map[string]string)
					m["ligand"] = r["ligand"]
					m["protein_path"] = r2["protein_path"]
					combinations = append(combinations, m)
				}
			}
		}
	}
	//fmt.Println(combinations)

	// Open the file for writing
	index_file := volume_path + "/index.jsonl"
	index_dict, err := os.Create(index_file)
	if err != nil {
		panic(err)
	}
	defer index_dict.Close()

	// Write each JSON object as a separate line in the file
	for _, m := range combinations {
		b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		_, err = index_dict.Write(b)
		if err != nil {
			panic(err)
		}

		_, err = index_dict.WriteString("\n")
		if err != nil {
			panic(err)
		}
	}
    return index_file, combinations
}

func indexCreateIndexCSV(index []map[string]string, volume_path string) string {
	// todo generalise the function
	index_csv := volume_path + "/index.csv"
	file, err := os.Create(index_csv)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"protein_path", "ligand"}
	if err := writer.Write(header); err != nil {
		panic(err)
	}

	for _, row := range index {
		proteinPath := row["protein_path"]
		ligand := row["ligand"]
		record := []string{proteinPath, ligand}
		if err := writer.Write(record); err != nil {
			panic(err)
		}
	}
	return index_csv
}

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
	fmt.Println("## Seaching input files ##")
	out := indexSearchDirectoryPath(in_dir, app_config, *layers)
	// TODO pass separate volume directory
	// TODO create dedicated id generator
	// TODO enable passing an array of multiple input directories
	fmt.Println("## Creating volume ##")
	_, new_out, volume := indexCreateInputsVolume(in_dir, out, "/inputs")
	fmt.Println("## Creating index ##")
	index_file, index := indexCreateIndexJSONL(new_out, app_config, volume)
	index_csv := indexCreateIndexCSV(index, volume)
	fmt.Println("Index jsonl created:", index_file)
	fmt.Println("Index csv created:", index_csv)
}
