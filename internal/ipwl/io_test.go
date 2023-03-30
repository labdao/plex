package ipwl

import (
	"reflect"
	"testing"
)

func TestReadIOLibrary(t *testing.T) {
	filePath := "testdata/example_io.jsonl"
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
