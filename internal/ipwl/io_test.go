package ipwl

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFindMatchingFiles(t *testing.T) {
	tool := Tool{
		Inputs: map[string]interface{}{
			"protein": map[string]interface{}{
				"type": "File",
				"glob": []interface{}{"*.pdb"},
			},
			"small_molecule": map[string]interface{}{
				"type": "File",
				"glob": []interface{}{"*.sdf", "*.mol2"},
			},
		},
	}

	inputDir := "../../testdata/binding/abl/"

	expected := map[string][]string{
		"protein":        {filepath.Join(inputDir, "7n9g.pdb")},
		"small_molecule": {filepath.Join(inputDir, "ZINC000003986735.sdf"), filepath.Join(inputDir, "ZINC000019632618.sdf")},
	}

	inputFilepaths, err := findMatchingFiles(inputDir, tool)
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

func TestCreateIOEntries(t *testing.T) {
	inputCombinations := []map[string]string{
		{
			"protein":        "testdata/binding/abl/7n9g.pdb",
			"small_molecule": "testdata/binding/abl/ZINC000003986735.sdf",
		},
		{
			"protein":        "testdata/binding/abl/7n9g.pdb",
			"small_molecule": "testdata/binding/abl/ZINC000019632618.sdf",
		},
	}

	tool := Tool{
		Name: "equibind",
		Inputs: map[string]interface{}{
			"protein": map[string]interface{}{
				"type": "File",
				"glob": []string{"*.pdb"},
			},
			"small_molecule": map[string]interface{}{
				"type": "File",
				"glob": []string{"*.sdf", "*.mol2"},
			},
		},
		Outputs: map[string]interface{}{
			"docked_small_molecule": map[string]interface{}{
				"type": "File",
				"glob": []string{"*_docked.sdf"},
			},
			"protein": map[string]interface{}{
				"type": "File",
				"glob": []string{"*.pdb"},
			},
		},
	}

	expected := []IO{
		{
			Tool:  "equibind",
			State: "created",
			Inputs: map[string]interface{}{
				"protein": map[string]interface{}{
					"class":    "File",
					"filepath": "testdata/binding/abl/7n9g.pdb",
				},
				"small_molecule": map[string]interface{}{
					"class":    "File",
					"filepath": "testdata/binding/abl/ZINC000003986735.sdf",
				},
			},
			Outputs: map[string]interface{}{
				"docked_small_molecule": map[string]interface{}{
					"class": "File",
				},
				"protein": map[string]interface{}{
					"class": "File",
				},
			},
		},
		{
			Tool:  "equibind",
			State: "created",
			Inputs: map[string]interface{}{
				"protein": map[string]interface{}{
					"class":    "File",
					"filepath": "testdata/binding/abl/7n9g.pdb",
				},
				"small_molecule": map[string]interface{}{
					"class":    "File",
					"filepath": "testdata/binding/abl/ZINC000019632618.sdf",
				},
			},
			Outputs: map[string]interface{}{
				"docked_small_molecule": map[string]interface{}{
					"class": "File",
				},
				"protein": map[string]interface{}{
					"class": "File",
				},
			},
		},
	}

	ioEntries := createIOEntries(tool, inputCombinations)

	if !reflect.DeepEqual(ioEntries, expected) {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, ioEntries)
	}
}

func TestCreateIOJson(t *testing.T) {
	inputDir := "../../testdata/binding/abl/"
	toolFilePath := "testdata/example_tool.json"

	// Read and unmarshal the Tool object from the example_tool.json file
	file, err := os.Open(toolFilePath)
	if err != nil {
		t.Fatalf("Error opening tool file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("Error reading tool file: %v", err)
	}

	var tool Tool
	err = json.Unmarshal(data, &tool)
	if err != nil {
		t.Fatalf("Error unmarshalling tool JSON: %v", err)
	}

	expected := []IO{
		{
			Tool:  "equibind",
			State: "created",
			Inputs: map[string]interface{}{
				"protein": map[string]interface{}{
					"class":    "File",
					"filepath": "testdata/binding/abl/7n9g.pdb",
				},
				"small_molecule": map[string]interface{}{
					"class":    "File",
					"filepath": "testdata/binding/abl/ZINC000003986735.sdf",
				},
			},
			Outputs: map[string]interface{}{
				"docked_small_molecule": map[string]interface{}{
					"class": "File",
				},
				"protein": map[string]interface{}{
					"class": "File",
				},
			},
		},
		{
			Tool:  "equibind",
			State: "created",
			Inputs: map[string]interface{}{
				"protein": map[string]interface{}{
					"class":    "File",
					"filepath": "testdata/binding/abl/7n9g.pdb",
				},
				"small_molecule": map[string]interface{}{
					"class":    "File",
					"filepath": "testdata/binding/abl/ZINC000019632618.sdf",
				},
			},
			Outputs: map[string]interface{}{
				"docked_small_molecule": map[string]interface{}{
					"class": "File",
				},
				"protein": map[string]interface{}{
					"class": "File",
				},
			},
		},
	}

	ioData, err := createIOJson(inputDir, tool)
	if err != nil {
		t.Fatalf("Error creating IO JSON: %v", err)
	}

	if !reflect.DeepEqual(ioData, expected) {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, ioData)
	}
}

func TestReadIOLibrary(t *testing.T) {
	filePath := "testdata/example_io.json"
	expected := []IO{
		{
			Tool: "testdata/ipwl_test/equibind",
			Inputs: map[string]interface{}{
				"protein": map[string]interface{}{
					"class":    "File",
					"filepath": "testdata/binding/abl/7n9g.pdb",
				},
				"small_molecule": map[string]interface{}{
					"class":    "File",
					"filepath": "testdata/binding/abl/ZINC000003986735.sdf",
				},
			},
			Outputs: map[string]interface{}{
				"docked_small_molecule": map[string]interface{}{
					"class": "File",
				},
				"protein": map[string]interface{}{
					"class": "File",
				},
			},
			State: "created",
		},
		{
			Tool: "testdata/ipwl_test/equibind",
			Inputs: map[string]interface{}{
				"protein": map[string]interface{}{
					"class":    "File",
					"basename": "7n9g.pdb",
					"filepath": "testdata/binding/abl/7n9g.pdb",
				},
				"small_molecule": map[string]interface{}{
					"class":    "File",
					"filepath": "testdata/binding/abl/ZINC000019632618.sdf",
				},
			},
			Outputs: map[string]interface{}{
				"docked_small_molecule": map[string]interface{}{
					"class": "File",
				},
				"protein": map[string]interface{}{
					"class": "File",
				},
			},
			State: "created",
		},
	}

	ioLibrary, err := readIOList(filePath)
	if err != nil {
		t.Fatalf("Error reading IO library: %v", err)
	}

	if !reflect.DeepEqual(ioLibrary, expected) {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, ioLibrary)
	}
}
