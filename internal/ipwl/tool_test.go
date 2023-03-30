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
			"python main.py --protein $(inputs.protein.filepath) --small_molecule_library $(inputs.small_molecule.filepath);",
			"mv /outputs/ligands_predicted.sdf /outputs/$(inputs.protein.filepath)_$(inputs.small_molecule.filepath)_docked.$(inputs.small_molecule.filepath);",
			"cp /inputs/$(inputs.protein.filepath) /outputs/",
		},
		DockerPull: "ghcr.io/labdao/equibind@sha256:ae2cec63b3924774727ed1c6c8af95cf4aaea2d3f0c5acbec56478505ccb2b07",
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

func TestToolToDockerCmd(t *testing.T) {
	toolConfig := Tool{
		DockerPull:  "some_docker_pull",
		BaseCommand: []string{"some_base_command"},
		Arguments:   []string{"--protein", "$(inputs.protein.filepath)", "--small_molecule", "$(inputs.small_molecule.filepath)"},
	}

	ioEntry := IO{
		Inputs: map[string]interface{}{
			"protein": map[string]interface{}{
				"filepath": "testdata/binding/abl/7n9g.pdb",
			},
			"small_molecule": map[string]interface{}{
				"filepath": "testdata/binding/abl/ZINC000003986735.sdf",
			},
		},
	}

	outputDirPath := "some_output_dir"

	dockerCmd, err := ToolToDockerCmd(toolConfig, ioEntry, outputDirPath)
	if err != nil {
		t.Errorf("Error generating Docker command: %v", err)
	}

	expectedDockerCmd := "docker -v testdata/binding/abl/7n9g.pdb:/inputs -v testdata/binding/abl/ZINC000003986735.sdf:/inputs -v some_output_dir:/outputs run some_docker_pull some_base_command --protein testdata/binding/abl/7n9g.pdb --small_molecule testdata/binding/abl/ZINC000003986735.sdf"

	if dockerCmd != expectedDockerCmd {
		t.Errorf("Expected Docker command: %s, got: %s", expectedDockerCmd, dockerCmd)
	}
}
