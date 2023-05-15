package ipwl

import (
	"io/ioutil"
	"log"
	"path/filepath"
)

func findMatchingFiles(inputDir string, tool Tool, layers int) (map[string][]string, error) {
	inputFilepaths := make(map[string][]string)

	err := searchMatchingFiles(inputDir, tool, layers, 0, inputFilepaths)
	if err != nil {
		return nil, err
	}

	return inputFilepaths, nil
}

func searchMatchingFiles(inputDir string, tool Tool, layers int, currentLayer int, inputFilepaths map[string][]string) error {
	if currentLayer > layers {
		return nil
	}

	for inputName, input := range tool.Inputs {
		if input.Type == "File" {
			for _, globPattern := range input.Glob {
				matches, err := filepath.Glob(filepath.Join(inputDir, globPattern))
				if err != nil {
					return err
				}

				inputFilepaths[inputName] = append(inputFilepaths[inputName], matches...)
			}
		}
	}

	// Search subdirectories
	subDirs, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return err
	}

	for _, subDir := range subDirs {
		if subDir.IsDir() {
			err := searchMatchingFiles(filepath.Join(inputDir, subDir.Name()), tool, layers, currentLayer+1, inputFilepaths)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func generateCombinationsRecursive(keys []string, values map[string][]string, index int, combination map[string]string, combinations *[]map[string]string) {
	if index == len(keys) {
		newCombination := make(map[string]string)
		for key, value := range combination {
			newCombination[key] = value
		}
		*combinations = append(*combinations, newCombination)
		return
	}

	key := keys[index]
	for _, value := range values[key] {
		combination[key] = value
		generateCombinationsRecursive(keys, values, index+1, combination, combinations)
	}
}

func generateInputCombinations(inputFilepaths map[string][]string) []map[string]string {
	var combinations []map[string]string
	keys := make([]string, 0, len(inputFilepaths))

	for key := range inputFilepaths {
		keys = append(keys, key)
	}

	generateCombinationsRecursive(keys, inputFilepaths, 0, make(map[string]string), &combinations)

	return combinations
}

func createIOEntries(toolPath string, tool Tool, inputCombinations []map[string]string) []IO {
	var ioData []IO

	for _, combination := range inputCombinations {
		ioEntry := IO{
			Tool:    toolPath,
			State:   "created",
			Inputs:  map[string]FileInput{},
			Outputs: map[string]CustomOutput{},
		}

		for inputName, path := range combination {
			absPath, err := filepath.Abs(path)
			if err != nil {
				log.Printf("Error converting to absolute path: %v", err)
				continue
			}

			ioEntry.Inputs[inputName] = FileInput{
				Class: "File",
				Address: FileAddress{
					FilePath: absPath,
					IPFS:     "",
				},
			}
		}

		for outputName, output := range tool.Outputs {
			if output.Type == "File" {
				ioEntry.Outputs[outputName] = CustomOutput{
					FileOutput: &FileOutput{
						Class: "File",
						Address: FileAddress{
							FilePath: "",
							IPFS:     "",
						},
					},
				}
			} else if output.Type == "Array" && output.Item == "File" {
				ioEntry.Outputs[outputName] = CustomOutput{
					ArrayFile: &ArrayFileOutput{
						Class: "Array",
						Files: []FileOutput{},
					},
				}
			}
		}

		ioData = append(ioData, ioEntry)
	}

	return ioData
}

func CreateIOJson(inputDir string, tool Tool, toolPath string, layers int) ([]IO, error) {
	inputFilepaths, err := findMatchingFiles(inputDir, tool, layers)
	if err != nil {
		return nil, err
	}

	inputCombinations := generateInputCombinations(inputFilepaths)
	ioData := createIOEntries(toolPath, tool, inputCombinations)

	return ioData, nil
}
