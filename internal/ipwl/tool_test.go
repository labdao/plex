package ipwl

import (
	"reflect"
	"testing"
)

func TestReadToolConfig(t *testing.T) {
	filePath := "testdata/example_tool.json"
	expected := Tool{
		Name:        "equibind",
		Description: "Docking of small molecules to a protein",
		BaseCommand: []string{"/bin/bash", "-c"},
		Arguments: []string{
			"python main.py --protein $(inputs.protein.filepath) --small_molecule_library $(inputs.small_molecule.filepath);",
			"mv /outputs/ligands_predicted.sdf /outputs/$(inputs.protein.basename)_$(inputs.small_molecule.basename)_docked.$(inputs.small_molecule.ext);",
			"cp $(inputs.protein.filepath) /outputs/;",
			"rmdir /outputs/dummy;",
		},
		DockerPull: "ghcr.io/labdao/equibind@sha256:ae2cec63b3924774727ed1c6c8af95cf4aaea2d3f0c5acbec56478505ccb2b07",
		GpuBool:    false,
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
		Outputs: map[string]ToolOutput{
			"best_docked_small_molecule": {
				Type: "File",
				Glob: []string{"*_docked.sdf"},
			},
			"protein": {
				Type: "File",
				Glob: []string{"*.pdb"},
			},
		},
	}

	tool, err := ReadToolConfig(filePath)
	if err != nil {
		t.Fatalf("Error reading tool config: %v", err)
	}

	if !reflect.DeepEqual(tool, expected) {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, tool)
	}
}
