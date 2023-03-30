package ipwl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type IO struct {
	Tool    string                 `json:"tool"`
	Inputs  map[string]interface{} `json:"inputs"`
	Outputs map[string]interface{} `json:"outputs"`
	State   string                 `json:"state"`
	ErrMsg  string                 `json:"errMsg`
}

func readIOList(filePath string) ([]IO, error) {
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

func writeIOList(ioJsonPath string, ioList []IO) error {
	file, err := os.OpenFile(ioJsonPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	for _, ioEntry := range ioList {
		err = encoder.Encode(ioEntry)
		if err != nil {
			return fmt.Errorf("failed to encode IO entry: %w", err)
		}
	}

	return nil
}
