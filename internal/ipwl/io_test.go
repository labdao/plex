package ipwl

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func loadJSONFile(filePath string, target interface{}) error {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileBytes, target)
	if err != nil {
		return err
	}

	return nil
}

func TestReadIOList(t *testing.T) {
	filePath := "testdata/example_io.json"

	var expected []IO
	err := loadJSONFile(filePath, &expected)
	if err != nil {
		t.Fatalf("Error loading example_io.json: %v", err)
	}

	ioList, err := ReadIOList(filePath)
	if err != nil {
		t.Fatalf("Error in ReadIOList: %v", err)
	}

	if !reflect.DeepEqual(ioList, expected) {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, ioList)
	}
}

func TestWriteIOList(t *testing.T) {
	ioJsonPath := "testdata/temp_io.json"
	defer os.Remove(ioJsonPath)

	var ioList []IO
	err := loadJSONFile("testdata/example_io.json", &ioList)
	if err != nil {
		t.Fatalf("Error loading example_io.json: %v", err)
	}

	err = WriteIOList(ioJsonPath, ioList)
	if err != nil {
		t.Fatalf("Error in WriteIOList: %v", err)
	}

	var writtenIOList []IO
	err = loadJSONFile(ioJsonPath, &writtenIOList)
	if err != nil {
		t.Fatalf("Error loading temp_io.json: %v", err)
	}

	if !reflect.DeepEqual(writtenIOList, ioList) {
		t.Errorf("Expected:\n%v\nGot:\n%v", ioList, writtenIOList)
	}
}
