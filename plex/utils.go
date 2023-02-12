package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	fileutils "github.com/docker/docker/pkg/fileutils"
	"github.com/google/uuid"
)

type AppConfig struct {
	App    string `json:"app"`
	Inputs []struct {
		Field     string   `json:"field"`
		Filetypes []string `json:"filetypes"`
	} `json:"inputs"`
	Outputs []string `json:"outputs"`
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

/*
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
*/

func findAppConfig(app string, appConfigsFilePath string) (AppConfig, error) {
	appConfig := AppConfig{}
	file, err := os.Open(appConfigsFilePath)
	if err != nil {
		return appConfig, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = json.Unmarshal([]byte(scanner.Text()), &appConfig)
		if err != nil {
			return appConfig, err
		}
		if appConfig.App == app {
			fmt.Println("App found:", appConfig.App)
			return appConfig, nil
		}
	}
	return appConfig, err
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

func searchDirectoryPath(directory *string, appConfig AppConfig, layers int) (files []string, err error) {
	// validate config file
	// ValidateAppConfig(appConfig)

	// walk the directory path
	err = filepath.Walk(*directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if layers > 0 {
				layers--
				return nil
			}
		}

		// keep files that match the input filetypes of the specified application
		fmt.Println("appConfig.Inputs", appConfig.Inputs)
		for _, input := range appConfig.Inputs {
			for _, filetype := range input.Filetypes {
				fmt.Println("path:", path, "extension search", filetype)
				fmt.Println(filetype)
				if strings.HasSuffix(path, filetype) {
					files = append(files, path)
				}
			}
		}
		return nil
	})

	if err != nil {
		return
	}
	return
}

func createInputsDirectory(inputsBasedir string, files []string) (string, []string, string, error) {
	// create job directory
	id := uuid.New()
	inputsPath := path.Join(inputsBasedir, id.String())
	err := os.Mkdir(inputsPath, 0755)
	if err != nil {
		return id.String(), []string{}, inputsPath, err
	}

	// create the inputs directory within the job directory
	os.Mkdir(inputsPath, 0755)

	// copy files to the inputs directory
	newFiles := make([]string, 0)
	for _, file := range files {
		_, err = fileutils.CopyFile(file, path.Join(inputsPath, filepath.Base(file)))
		if err != nil {
			return id.String(), newFiles, inputsPath, err
		}
		newFiles = append(newFiles, filepath.Base(file))
	}
	return id.String(), newFiles, inputsPath, nil
}

func createCombinations(indexMap map[string][]string, fieldA, fieldB string) []map[string]string {
	// generate combinations of the mapping
	combinations := []map[string]string{}
	for _, valA := range indexMap[fieldA] {
		for _, valB := range indexMap[fieldB] {
			row := map[string]string{fieldA: valA, fieldB: valB}
			combinations = append(combinations, row)
		}
	}
	return combinations
}

func createIndex(filePaths []string, appConfig AppConfig, jobDirPath string) (string, []map[string]string) {
	// categorise the input files based on the app config specifications
	indexMap := map[string][]string{}
	for _, filePath := range filePaths {
		for _, input := range appConfig.Inputs {
			for _, filetype := range input.Filetypes {
				if strings.HasSuffix(filePath, filetype) {
					indexMap[input.Field] = append(indexMap[input.Field], filePath)
				}
			}
		}
	}

	fieldA := appConfig.Inputs[0].Field
	fieldB := appConfig.Inputs[1].Field
	combinations := createCombinations(indexMap, fieldA, fieldB)
	writeJSONL(combinations, path.Join(jobDirPath, "index.jsonl"))
	writeCSV(combinations, path.Join(jobDirPath, "index.csv"))
	return path.Join(jobDirPath, "index.csv"), combinations
}
