package ipwl

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/labdao/plex/internal/web3"
	"gorm.io/gorm"
)

var (
	inputs           string
	scatteringMethod string
)

func InitializeIo(toolPath string, scatteringMethod string, inputVectors map[string][]interface{}, db *gorm.DB) ([]IO, error) {
	// Open the file and load its content
	tool, toolInfo, err := ReadToolConfig(toolPath, db)
	if err != nil {
		return nil, err
	}

	// Validate input keys
	err = validateInputKeys(inputVectors, tool.Inputs)
	if err != nil {
		return nil, err
	}

	// Handle scattering methods and create the ist
	var inputsList [][]interface{}
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

func validateInputKeys(inputVectors map[string][]interface{}, toolInputs map[string]ToolInput) error {
	for inputKey := range inputVectors {
		if _, exists := toolInputs[inputKey]; !exists {
			log.Printf("The argument %s is not in the tool inputs.\n", inputKey)
			log.Printf("Available keys: %v\n", toolInputs)
			return fmt.Errorf("the argument %s is not in the tool inputs", inputKey)
		}
	}
	return nil
}

func dotProductScattering(inputVectors map[string][]interface{}) ([][]interface{}, error) {
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

	var inputsList [][]interface{}
	keys := make([]string, 0, len(inputVectors))
	for k := range inputVectors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := 0; i < vectorLength; i++ {
		var tmp []interface{}
		for _, k := range keys {
			tmp = append(tmp, inputVectors[k][i])
		}
		inputsList = append(inputsList, tmp)
	}

	return inputsList, nil
}

func crossProductScattering(inputVectors map[string][]interface{}) ([][]interface{}, error) {
	// Cartesian product function adapted for slices of interfaces
	cartesian := func(arrs ...[]interface{}) [][]interface{} {
		result := [][]interface{}{{}}
		for _, arr := range arrs {
			var temp [][]interface{}
			for _, res := range result {
				for _, item := range arr {
					product := append([]interface{}{}, res...) // Copy current slice of interfaces
					product = append(product, item)            // Append the new item
					temp = append(temp, product)               // Append the new product to the temporary result
				}
			}
			result = temp // Set the result to the temporary result
		}
		return result
	}

	// Extracting the keys and sorting them
	keys := make([]string, 0, len(inputVectors))
	for k := range inputVectors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Preparing the input for the cartesian product function
	arrays := make([][]interface{}, len(keys))
	for i, key := range keys {
		arrays[i] = inputVectors[key] // Directly assign the slice from the map
	}

	// Computing the Cartesian product
	inputsList := cartesian(arrays...)

	return inputsList, nil
}

func createSingleIo(inputs []interface{}, tool Tool, toolInfo ToolInfo, userId string, inputVectors map[string][]interface{}) (IO, error) {
	io := IO{
		Tool:    toolInfo,
		Inputs:  make(map[string]interface{}),
		Outputs: make(map[string]interface{}),
		State:   "created",
		ErrMsg:  "",
		UserID:  userId,
	}

	inputKeys := make([]string, 0, len(inputVectors))
	for k := range inputVectors {
		inputKeys = append(inputKeys, k)
	}

	// Sort the inputKeys slice to ensure a consistent order
	sort.Strings(inputKeys)

	for i, inputValue := range inputs {
		inputKey := inputKeys[i]
		io.Inputs[inputKey] = inputValue
	}

	for outputKey, outputValue := range tool.Outputs {
		io.Outputs[outputKey] = outputValue
	}

	return io, nil
}
