package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	fileutils "github.com/docker/docker/pkg/fileutils"
	"github.com/google/uuid"
)

// Structure for AppConfig
type appStruct struct {
	App     string      `json:"app"`
	Inputs  [][2]string `json:"inputs"`
	Outputs []string    `json:"outputs"`
}

func ValidateDirectoryPath(directory string) (bool, error) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return false, err
	}
	if fileInfo, err := os.Stat(directory); err == nil && !fileInfo.Mode().IsDir() {
		return false, err
	}
	return true, nil
}

func ValidateAppConfig(appConfig string) (bool, error) {
	file, err := os.Open(appConfig)
	if err != nil {
		return false, err
	}
	defer file.Close()

	var appData appStruct
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = json.Unmarshal([]byte(scanner.Text()), &appData)
		if err != nil {
			return false, err
		}
		break
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	if _, err := os.Stat(appConfig); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

func ValidateApplication(application string, appConfig string) {
	ValidateAppConfig(appConfig)
	file, err := os.Open(appConfig)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var appData appStruct
		err = json.Unmarshal([]byte(scanner.Text()), &appData)
		if err != nil {
			fmt.Println("Error unmarshalling application file JSON:", err)
			return
		}
		if appData.App == application {
			fmt.Println("Application found:", appData.App)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}
}

func writeJSONL(index_map []map[string]string, file string) {
	// Open the file for writing
	file_dict, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer file_dict.Close()

	// Write each JSON object as a separate line in the file
	for _, m := range index_map {
		b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		_, err = file_dict.Write(b)
		if err != nil {
			panic(err)
		}

		_, err = file_dict.WriteString("\n")
		if err != nil {
			panic(err)
		}
	}
}

func writeCSV(index_map []map[string]string, file string) string {
	// todo generalise the function beyond diffdock
	file_dict, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer file_dict.Close()

	writer := csv.NewWriter(file_dict)
	defer writer.Flush()

	header := []string{"protein_path", "ligand"}
	if err := writer.Write(header); err != nil {
		panic(err)
	}

	for _, row := range index_map {
		proteinPath := row["protein_path"]
		ligand := row["ligand"]
		record := []string{proteinPath, ligand}
		if err := writer.Write(record); err != nil {
			panic(err)
		}
	}
	return file
}

func searchDirectoryPath(directory *string, appConfig string, layers int) []string {
	// validate config file
	ValidateAppConfig(appConfig)

	// read the app.jsonl file
	file, err := os.Open(appConfig)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// read the file line by line
	var appData appStruct
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = json.Unmarshal([]byte(scanner.Text()), &appData)
		if err != nil {
			panic(err)
		}
		break
	}

	// additional errors
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// walk the directory path
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

		//keep files that match the input file extensions of the specified application
		//TODO: create safeguard to constrain ext to the 2nd element input array
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
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
	return files
}

func createInputsDirectory(inputs_basedir string, files []string, prefix string) (string, []string, string) {
	// create job directory
	id := uuid.New()
	inputs_path := inputs_basedir + "/" + id.String()
	err := os.Mkdir(inputs_path, 0755)
	if err != nil {
		panic(err)
	}

	// create the inputs directory within the job directory
	os.Mkdir(inputs_path+prefix, 0755)

	// copy files to the inputs directory
	new_files := make([]string, 0)
	for _, file := range files {
		_, err = fileutils.CopyFile(file, inputs_path+prefix+"/"+filepath.Base(file))
		if err != nil {
			fmt.Println("Error copying file to inputs directory")
			panic(err)
		}
		new_files = append(new_files, prefix+"/"+filepath.Base(file))
	}
	print("job directory created: ", inputs_path, "\n")
	return id.String(), new_files, inputs_path
}

func createCombinations(index_map []map[string]string) []map[string]string {
	// generate combinations of the mapping
	//TODO implement generalisable version
	combinations := []map[string]string{}
	for _, r := range index_map {
		if r["ligand"] != "" {
			for _, r2 := range index_map {
				if r2["protein_path"] != "" && r2["ligand"] == "" {
					m := make(map[string]string)
					m["ligand"] = r["ligand"]
					m["protein_path"] = r2["protein_path"]
					combinations = append(combinations, m)
				}
			}
		}
	}
	return combinations
}

func createIndex(newFiles []string, appConfig string, volumePath string) (string, []map[string]string) {
	// read the app.jsonl file
	file, err := os.Open(appConfig)
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

	// categorise the input files based on the app config specifications
	var sorted []map[string]string // Define a slice to store the maps

	for _, file := range newFiles {
		m := make(map[string]string)
		for _, mapping := range appData.Inputs {
			if strings.HasSuffix(file, mapping[1]) {
				m[mapping[0]] = file
				break
			}
		}
		sorted = append(sorted, m)
	}

	combinations := createCombinations(sorted)
	writeJSONL(combinations, volumePath+"/index.jsonl")
	writeCSV(combinations, volumePath+"/index.csv")
	return volumePath + "/index.csv", combinations
}

func errorCheck(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
