package ipwl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Output interface {
	OutputType() string
}

type FileAddress struct {
	FilePath string `json:"filepath"`
	IPFS     string `json:"ipfs"`
}

type FileInput struct {
	Class   string      `json:"class"`
	Address FileAddress `json:"address"`
}

type FileOutput struct {
	Class   string      `json:"class"`
	Address FileAddress `json:"address"`
}

func (fo FileOutput) OutputType() string {
	return fo.Class
}

type ArrayFileOutput struct {
	Class string       `json:"class"`
	Files []FileOutput `json:"files"`
}

func (afo ArrayFileOutput) OutputType() string {
	return afo.Class
}

type CustomOutput struct {
	FileOutput *FileOutput
	ArrayFile  *ArrayFileOutput
}

func (co CustomOutput) OutputType() string {
	if co.FileOutput != nil {
		return co.FileOutput.Class
	}
	return co.ArrayFile.Class
}

type IO struct {
	Tool    string                  `json:"tool"`
	Inputs  map[string]FileInput    `json:"inputs"`
	Outputs map[string]CustomOutput `json:"outputs"`
	State   string                  `json:"state"`
	ErrMsg  string                  `json:"errMsg"`
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

	io.Outputs = make(map[string]CustomOutput)
	for k, v := range aux.Outputs {
		var fileOutput FileOutput
		if err := json.Unmarshal(v, &fileOutput); err == nil && fileOutput.Class != "Array" {
			io.Outputs[k] = CustomOutput{FileOutput: &fileOutput}
		} else {
			var arrayFileOutput ArrayFileOutput
			if err := json.Unmarshal(v, &arrayFileOutput); err == nil && arrayFileOutput.Class == "Array" {
				io.Outputs[k] = CustomOutput{ArrayFile: &arrayFileOutput}
			} else {
				return err
			}
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
