package ipwl

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestFindMatchingFiles(t *testing.T) {
	tool := Tool{
		Inputs: map[string]ToolInput{
			"protein": {
				Type: "File",
				Glob: []string{"*.pdb"},
			},
			"small_molecule": {
				Type: "File",
				Glob: []string{"*.sdf", "*.mol2"},
			},
		},
	}

	inputDir := "testdata/binding/abl/"

	expected := map[string][]string{
		"protein":        {filepath.Join(inputDir, "7n9g.pdb")},
		"small_molecule": {filepath.Join(inputDir, "ZINC000003986735.sdf"), filepath.Join(inputDir, "ZINC000019632618.sdf")},
	}

	layers := 1
	inputFilepaths, err := findMatchingFiles(inputDir, tool, layers)
	if err != nil {
		t.Fatalf("findMatchingFiles returned an error: %v", err)
	}

	if !reflect.DeepEqual(inputFilepaths, expected) {
		t.Errorf("findMatchingFiles returned unexpected results\nGot: %v\nExpected: %v", inputFilepaths, expected)
	}
}

func TestGenerateInputCombinations(t *testing.T) {
	inputFilepaths := map[string][]string{
		"protein": {
			"testdata/binding/abl/7n9g.pdb",
		},
		"small_molecule": {
			"testdata/binding/abl/ZINC000003986735.sdf",
			"testdata/binding/abl/ZINC000019632618.sdf",
		},
	}

	expected := []map[string]string{
		{
			"protein":        "testdata/binding/abl/7n9g.pdb",
			"small_molecule": "testdata/binding/abl/ZINC000003986735.sdf",
		},
		{
			"protein":        "testdata/binding/abl/7n9g.pdb",
			"small_molecule": "testdata/binding/abl/ZINC000019632618.sdf",
		},
	}

	combinations := generateInputCombinations(inputFilepaths)

	if !reflect.DeepEqual(combinations, expected) {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, combinations)
	}
}

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

func TestCreateIOEntries(t *testing.T) {
	var ios []IO
	err := loadJSONFile("testdata/example_initial_io.json", &ios)
	if err != nil {
		t.Fatalf("Error loading example_initiail_io.json: %v", err)
	}

	var tool Tool
	err = loadJSONFile("testdata/example_tool.json", &tool)
	if err != nil {
		t.Fatalf("Error loading example_tool.json: %v", err)
	}

	toolPath := ios[0].Tool

	inputCombinations := make([]map[string]string, len(ios))
	for i, io := range ios {
		inputCombination := make(map[string]string)
		for key, fileInput := range io.Inputs {
			inputCombination[key] = fileInput.FilePath
		}
		inputCombinations[i] = inputCombination
	}

	expected := make([]IO, len(ios))
	for i := range ios {
		expected[i] = IO{
			Tool:    ios[i].Tool,
			Inputs:  ios[i].Inputs,
			Outputs: ios[i].Outputs,
			State:   ios[i].State,
			ErrMsg:  ios[i].ErrMsg,
		}
	}

	ioEntries := createIOEntries(toolPath, tool, inputCombinations)

	if !reflect.DeepEqual(ioEntries, expected) {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, ioEntries)
	}
}

func TestCreateIOJson(t *testing.T) {
	inputDir := "testdata/binding/abl"

	var ios []IO
	err := loadJSONFile("testdata/example_equibind_io.json", &ios)
	if err != nil {
		t.Fatalf("Error loading example_equibind_io.json: %v", err)
	}

	var tool Tool
	err = loadJSONFile("testdata/example_tool.json", &tool)
	if err != nil {
		t.Fatalf("Error loading example_tool.json: %v", err)
	}

	// Extract the toolPath from the first item in the ios
	toolPath := ios[0].Tool

	// Get the expected inputs and outputs
	expected := ios

	// Remove the "tool" key from each map in ios
	for i := range ios {
		ios[i].Tool = ""
	}

	generatedIOData, err := CreateIOJson(inputDir, tool, toolPath, 2)
	if err != nil {
		t.Fatalf("Error in CreateIOJson: %v", err)
	}

	// Compare the paths after the asterisk
	for i := range generatedIOData {
		for k, v := range generatedIOData[i].Inputs {
			expectedPath := strings.Split(expected[i].Inputs[k].FilePath, "*")[1]
			if !strings.HasSuffix(v.FilePath, expectedPath) {
				t.Errorf("Expected path suffix:\n%v\nGot:\n%v", expectedPath, v.FilePath)
			}
		}
	}
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
