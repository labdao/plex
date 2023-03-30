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
