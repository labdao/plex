package ipwl

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/web3"
)

var (
	inputs           string
	scatteringMethod string
)

func InitializeIo(toolPath string, scatteringMethod string, inputVectors map[string][]string) ([]IO, error) {
	// Open the file and load its content
	tool, toolInfo, err := ReadToolConfig(toolPath)
	if err != nil {
		return nil, err
	}

	// Validate input keys
	err = validateInputKeys(inputVectors, tool.Inputs)
	if err != nil {
		return nil, err
	}

	// Handle scattering methods and create the ist
	var inputsList [][]string
	switch scatteringMethod {
	case "dotProduct":
		inputsList, err = dotProductScattering(inputVectors)
		if err != nil {
			return nil, err
		}
	case "crossProduct":
		inputsList, err = crossProductScattering(inputVectors)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid scattering method: %s", scatteringMethod)
	}

	var userId string

	if web3.IsValidEthereumAddress(os.Getenv("RECIPIENT_WALLET")) {
		userId = os.Getenv("RECIPIENT_WALLET")
	} else {
		fmt.Print("Invalid wallet address detected. Using empty string for user ID.\n")
		userId = ""
	}

	var ioList []IO

	for _, inputs := range inputsList {
		io, err := createSingleIo(inputs, tool, toolInfo, userId, inputVectors)
		if err != nil {
			return nil, err
		}
		ioList = append(ioList, io)
	}

	return ioList, nil
}

func validateInputKeys(inputVectors map[string][]string, toolInputs map[string]ToolInput) error {
	for inputKey := range inputVectors {
		if _, exists := toolInputs[inputKey]; !exists {
			log.Printf("The argument %s is not in the tool inputs.\n", inputKey)
			log.Printf("Available keys: %v\n", toolInputs)
			return fmt.Errorf("the argument %s is not in the tool inputs", inputKey)
		}
	}
	return nil
}

func dotProductScattering(inputVectors map[string][]string) ([][]string, error) {
	var vectorLength int
	for _, v := range inputVectors {
		if vectorLength == 0 {
			vectorLength = len(v)
			continue
		}
		if len(v) != vectorLength {
			return nil, fmt.Errorf("all input arguments must have the same length for dotProduct scattering method")
		}
	}

	var inputsList [][]string
	for i := 0; i < vectorLength; i++ {
		tmp := []string{}
		for _, v := range inputVectors {
			tmp = append(tmp, v[i])
		}
		inputsList = append(inputsList, tmp)
	}

	return inputsList, nil
}

func crossProductScattering(inputVectors map[string][]string) ([][]string, error) {
	cartesian := func(arrs ...[]string) [][]string {
		result := [][]string{{}}
		for _, arr := range arrs {
			var temp [][]string
			for _, res := range result {
				for _, str := range arr {
					product := append([]string{}, res...)
					product = append(product, str)
					temp = append(temp, product)
				}
			}
			result = temp
		}
		return result
	}

	keys := make([]string, 0, len(inputVectors))
	for k := range inputVectors {
		keys = append(keys, k)
	}
	arrays := make([][]string, len(inputVectors))
	for i, k := range keys {
		arrays[i] = inputVectors[k]
	}
	inputsList := cartesian(arrays...)

	return inputsList, nil
}

func createSingleIo(inputs []string, tool Tool, toolInfo ToolInfo, userId string, inputVectors map[string][]string) (IO, error) {
	io := IO{
		Tool:    toolInfo,
		Inputs:  make(map[string]FileInput),
		Outputs: make(map[string]Output),
		State:   "created",
		ErrMsg:  "",
		UserID:  userId,
	}

	inputKeys := make([]string, 0, len(inputVectors))
	for k := range inputVectors {
		inputKeys = append(inputKeys, k)
	}

	for i, inputValue := range inputs {
		inputKey := inputKeys[i]

		if strings.Count(inputValue, "/") == 1 {
			parts := strings.Split(inputValue, "/")
			cid := parts[0]
			fileName := parts[1]
			if !ipfs.IsValidCID(cid) {
				return io, fmt.Errorf("invalid CID: %s", cid)
			}
			io.Inputs[inputKey] = FileInput{
				Class:    tool.Inputs[inputKey].Type,
				FilePath: fileName,
				IPFS:     cid,
			}
		} else {
			cid, err := ipfs.WrapAndPinFile(inputValue)
			if err != nil {
				return io, err
			}
			io.Inputs[inputKey] = FileInput{
				Class:    tool.Inputs[inputKey].Type,
				FilePath: filepath.Base(inputValue),
				IPFS:     cid,
			}
		}
	}

	for outputKey, outputValue := range tool.Outputs {
		io.Outputs[outputKey] = FileOutput{
			Class:    outputValue.Type,
			FilePath: "",
			IPFS:     "",
		}
	}

	return io, nil
}

// NOTE
// all functions below will be deprecated soon with the removal of plex create

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

func createIOEntries(toolInfo ToolInfo, tool Tool, inputCombinations []map[string]string) []IO {
	var ioData []IO

	var userID string

	if web3.IsValidEthereumAddress(os.Getenv("RECIPIENT_WALLET")) {
		userID = os.Getenv("RECIPIENT_WALLET")
	} else {
		fmt.Println("RECIPIENT_WALLET is not a valid Ethereum address. Using empty string as user ID.")
		userID = ""
	}

	for _, combination := range inputCombinations {
		ioEntry := IO{
			Tool:    toolInfo,
			State:   "created",
			Inputs:  map[string]FileInput{},
			Outputs: map[string]Output{},
			UserID:  userID,
		}

		for inputName, path := range combination {
			_, fileName := filepath.Split(path)

			cid, err := ipfs.WrapAndPinFile(path)
			if err != nil {
				log.Printf("Error getting CID for file %s: %v", path, err)
				continue
			}

			ioEntry.Inputs[inputName] = FileInput{
				Class:    "File",
				FilePath: fileName,
				IPFS:     cid,
			}
		}

		for outputName, output := range tool.Outputs {
			if output.Type == "File" {
				ioEntry.Outputs[outputName] = FileOutput{
					Class: "File",
				}
			} else if output.Type == "Array" && output.Item == "File" {
				ioEntry.Outputs[outputName] = ArrayFileOutput{
					Class: "Array",
					Files: []FileOutput{},
				}
			}
		}

		ioData = append(ioData, ioEntry)
	}

	return ioData
}

func CreateIOJson(inputDir string, tool Tool, toolInfo ToolInfo, layers int) ([]IO, error) {
	inputFilepaths, err := findMatchingFiles(inputDir, tool, layers)
	if err != nil {
		return nil, err
	}

	inputCombinations := generateInputCombinations(inputFilepaths)
	ioData := createIOEntries(toolInfo, tool, inputCombinations)

	return ioData, nil
}
