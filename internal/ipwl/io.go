package ipwl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type ToolInfo struct {
	Name string `json:"name"`
	IPFS string `json:"ipfs"`
}

type IO struct {
	Tool          ToolInfo               `json:"tool"`
	Inputs        map[string]interface{} `json:"inputs"`
	Outputs       map[string]interface{} `json:"outputs"`
	State         string                 `json:"state"`
	ErrMsg        string                 `json:"errMsg"`
	UserID        string                 `json:"userId"`
	BacalhauJobId string                 `json:"bacalhauJobId"`
}

func ReadIOList(filePath string) ([]IO, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var ioLibrary []IO
	err = json.Unmarshal(data, &ioLibrary)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return ioLibrary, nil
}

func WriteIOList(ioJsonPath string, ioList []IO) error {
	file, err := os.OpenFile(ioJsonPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(ioList)
	if err != nil {
		return fmt.Errorf("failed to encode IO list: %w", err)
	}

	return nil
}

func ContainsUserIdAnnotation(slice []string) bool {
	for _, a := range slice {
		if strings.HasPrefix(a, "userId=") {
			return true
		}
	}
	return false
}

func ExtractUserIDFromIOJson(ioJsonPath string) (string, error) {
	file, err := os.Open(ioJsonPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// TODO: this is where we were previously failing on the gateway side
	var ioEntries []IO
	err = json.Unmarshal(data, &ioEntries)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if len(ioEntries) == 0 {
		return "", fmt.Errorf("no IO entries found")
	}

	return ioEntries[0].UserID, nil
}

func PrintIOGraphStatus(ioList []IO) {
	stateCount := make(map[string]int)

	// Iterate through the ioList and count the occurrences of each state
	for _, io := range ioList {
		stateCount[io.State]++
	}

	// Print the total number of IOs
	fmt.Printf("Total IOs: %d\n", len(ioList))

	// Print the number of IOs in each state
	for state, count := range stateCount {
		fmt.Printf("IOs in %s state: %d\n", state, count)
	}
}

type OutputValues struct {
	FilePaths []string `json:"filePaths"`
	CIDs      []string `json:"cids"`
	CidPaths  []string `json:"cidPaths"`
}
