package ipwl

import (
	"reflect"
	"testing"
)

func TestReadToolConfig(t *testing.T) {
	filePath := "testdata/example_tool.json"
	expected := Tool{
		Name:        "test",
		Description: "test tool",
		BaseCommand: []string{"/bin/bash", "-c"},
		Arguments: []string{
			"python main.py --protein $(inputs.protein.path) --small_molecule_library $(inputs.small_molecule.path);",
			"cp /inputs /outputs;",
		},
		DockerPull: "ghcr.io/labdao/testtool",
		GpuBool:    false,
		Inputs: map[string]interface{}{
			"protein": map[string]interface{}{
				"type": "File",
				"glob": "*.pdb",
			},
			"small_molecule": map[string]interface{}{
				"type": "File",
				"glob": []interface{}{"*.sdf", "*.mol"},
			},
			"outputs": map[string]interface{}{
				"docked_small_molecule": map[string]interface{}{
					"type": "File",
					"glob": "*_docked.sdf",
					"protein": map[string]interface{}{
						"type": "File",
						"glob": "*.pdb",
					},
				},
			},
		},
	}

	tool, err := readToolConfig(filePath)
	if err != nil {
		t.Fatalf("Error reading tool config: %v", err)
	}

	if !reflect.DeepEqual(tool, expected) {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, tool)
	}
}
