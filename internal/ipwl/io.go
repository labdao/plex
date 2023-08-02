package ipwl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type FileInput struct {
	Class    string `json:"class"`
	FilePath string `json:"filepath"`
	IPFS     string `json:"ipfs"`
}

type FileOutput struct {
	Class    string `json:"class"`
	FilePath string `json:"filepath"`
	IPFS     string `json:"ipfs"`
}

func (fo FileOutput) OutputType() string {
	return fo.Class
}

type Output interface {
	OutputType() string
}

type ArrayFileOutput struct {
	Class string       `json:"class"`
	Files []FileOutput `json:"files"`
}

func (afo ArrayFileOutput) OutputType() string {
	return afo.Class
}

type ToolInfo struct {
	Name string `json:"name"`
	IPFS string `json:"ipfs"`
}

type IO struct {
	Tool          ToolInfo             `json:"tool"`
	Inputs        map[string]FileInput `json:"inputs"`
	Outputs       map[string]Output    `json:"outputs"`
	State         string               `json:"state"`
	ErrMsg        string               `json:"errMsg"`
	BacalhauJobId string               `json:"bacalhauJobId"`
}

func (io *IO) UnmarshalJSON(data []byte) error {
	type Alias IO
	aux := struct {
		Outputs map[string]json.RawMessage `json:"outputs"`
		*Alias
	}{
		Alias: (*Alias)(io),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	io.Outputs = make(map[string]Output)
	for k, v := range aux.Outputs {
		var fileOutput FileOutput
		if err := json.Unmarshal(v, &fileOutput); err != nil {
			return err
		}

		if fileOutput.Class == "Array" {
			var arrayFileOutput ArrayFileOutput
			if err := json.Unmarshal(v, &arrayFileOutput); err != nil {
				return err
			}
			io.Outputs[k] = arrayFileOutput
		} else {
			io.Outputs[k] = fileOutput
		}
	}

	return nil
}

func (io IO) MarshalJSON() ([]byte, error) {
	type Alias IO
	aux := struct {
		Outputs map[string]interface{} `json:"outputs"`
		*Alias
	}{
		Outputs: make(map[string]interface{}),
		Alias:   (*Alias)(&io),
	}
	for k, v := range io.Outputs {
		aux.Outputs[k] = v
	}
	return json.Marshal(aux)
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
