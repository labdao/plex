package ipwl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ModelInfo struct {
	Name string `json:"name"`
	S3   string `json:"s3"`
}

type NullableMap map[string]interface{}

func (nm *NullableMap) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*nm = nil
		return nil
	}
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	*nm = m
	return nil
}

func (nm NullableMap) MarshalJSON() ([]byte, error) {
	if nm == nil {
		return []byte("null"), nil
	}
	return json.Marshal(map[string]interface{}(nm))
}

type IO struct {
	Model         ModelInfo   `json:"modelInfo"`
	Inputs        NullableMap `json:"inputs"`
	Outputs       NullableMap `json:"outputs"`
	State         string      `json:"state"`
	ErrMsg        string      `json:"errMsg"`
	WalletAddress string      `json:"walletAddress"`
	RayJobID      string      `json:"rayJobId"`
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
