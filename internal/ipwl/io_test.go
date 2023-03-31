package ipwl

import (
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

	toolPath := "../ipwl/equibind"

	expected := []IO{
		{
			Tool:  toolPath,
			State: "created",
			Inputs: map[string]FileInput{
				"protein": {
					Class:    "File",
					FilePath: filepath.FromSlash("testdata/binding/abl/7n9g.pdb"),
				},
				"small_molecule": {
					Class:    "File",
					FilePath: filepath.FromSlash("testdata/binding/abl/ZINC000003986735.sdf"),
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
			Tool:  toolPath,
			State: "created",
			Inputs: map[string]FileInput{
				"protein": {
					Class:    "File",
					FilePath: filepath.FromSlash("testdata/binding/abl/7n9g.pdb"),
				},
				"small_molecule": {
					Class:    "File",
					FilePath: filepath.FromSlash("testdata/binding/abl/ZINC000019632618.sdf"),
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

	ioEntries := createIOEntries(toolPath, tool, inputCombinations)

	if !reflect.DeepEqual(ioEntries, expected) {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, ioEntries)
	}
}

func TestCreateIOJson(t *testing.T) {
	inputDir := "testdata/binding/abl/"
	toolPath := "../ipwl/equibind"

	tool := Tool{
		Name: "equibind",
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
		Outputs: map[string]interface{}{
			"docked_small_molecule": map[string]interface{}{
				"type": "File",
				"glob": []interface{}{"*_docked.sdf"},
			},
			"protein": map[string]interface{}{
				"type": "File",
				"glob": []interface{}{"*.pdb"},
			},
		},
	}

	expected := []IO{
		{
			Tool:  toolPath,
			State: "created",
			Inputs: map[string]FileInput{
				"protein": {
					Class:    "File",
					FilePath: filepath.FromSlash("testdata/binding/abl/7n9g.pdb"),
				},
				"small_molecule": {
					Class:    "File",
					FilePath: filepath.FromSlash("testdata/binding/abl/ZINC000003986735.sdf"),
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
			Tool:  toolPath,
			State: "created",
			Inputs: map[string]FileInput{
				"protein": {
					Class:    "File",
					FilePath: filepath.FromSlash("testdata/binding/abl/7n9g.pdb"),
				},
				"small_molecule": {
					Class:    "File",
					FilePath: filepath.FromSlash("testdata/binding/abl/ZINC000019632618.sdf"),
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

	ioData, err := CreateIOJson(inputDir, tool, toolPath)
	if err != nil {
		t.Fatalf("CreateIOJson returned an error: %v", err)
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
			Inputs: map[string]FileInput{
				"protein": {
					Class:    "File",
					FilePath: "testdata/binding/abl/7n9g.pdb",
				},
				"small_molecule": {
					Class:    "File",
					FilePath: "testdata/binding/abl/ZINC000003986735.sdf",
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
			Inputs: map[string]FileInput{
				"protein": {
					Class:    "File",
					FilePath: "testdata/binding/abl/7n9g.pdb",
				},
				"small_molecule": {
					Class:    "File",
					FilePath: "testdata/binding/abl/ZINC000019632618.sdf",
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

func TestWriteIOList(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "io_json_test_*.json")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	ioList := []IO{
		{
			Tool: "testdata/ipwl_test/equibind",
			Inputs: map[string]FileInput{
				"protein": {
					Class:    "File",
					FilePath: "testdata/binding/abl/7n9g.pdb",
				},
				"small_molecule": {
					Class:    "File",
					FilePath: "testdata/binding/abl/ZINC000003986735.sdf",
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

	err = WriteIOList(tmpFile.Name(), ioList)
	if err != nil {
		t.Fatalf("Error writing IO list: %v", err)
	}

	readIOList, err := readIOList(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error reading IO list: %v", err)
	}

	if !reflect.DeepEqual(ioList, readIOList) {
		t.Errorf("Expected:\n%v\nGot:\n%v", ioList, readIOList)
	}
}
