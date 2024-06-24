package ipwl

// func TestReadModelConfig(t *testing.T) {
// 	filePath := "testdata/example_model.json"
// 	expected := Model{
// 		Name:        "equibind",
// 		Author:      "@misc{stärk2022equibind,\n      title={EquiBind: Geometric Deep Learning for Drug Binding Structure Prediction}, \n      author={Hannes Stärk and Octavian-Eugen Ganea and Lagnajit Pattanaik and Regina Barzilay and Tommi Jaakkola},\n      year={2022},\n      eprint={2202.05146},\n      archivePrefix={arXiv},\n      primaryClass={q-bio.BM}\n}",
// 		Description: "Docking of small molecules to a protein",
// 		BaseCommand: []string{"/bin/bash", "-c"},
// 		Arguments: []string{
// 			"mkdir -p /tmp-inputs/tmp;",
// 			"mkdir -p /tmp-outputs/tmp;",
// 			"cp /inputs/* /tmp-inputs/tmp/;",
// 			"ls /tmp-inputs/tmp;",
// 			"cd /src && python /src/inference.py --config=/src/configs_clean/bacalhau.yml;",
// 			"mv /tmp-outputs/tmp/* /outputs/;",
// 			"mv /outputs/lig_equibind_corrected.sdf /outputs/$(inputs.protein.basename)_$(inputs.small_molecule.basename)_docked.$(inputs.small_molecule.ext);",
// 			"mv /tmp-inputs/tmp/*.pdb /outputs/;",
// 		},
// 		DockerPull: "ghcr.io/labdao/equibind:main@sha256:21a381d9ab1ff047565685044569c8536a55e489c9531326498b28d6b3cc244f",
// 		GpuBool:    false,
// 		Inputs: map[string]ModelInput{
// 			"protein": {
// 				Type: "File",
// 				Glob: []string{"*.pdb"},
// 			},
// 			"small_molecule": {
// 				Type: "File",
// 				Glob: []string{"*.sdf", "*.mol2"},
// 			},
// 		},
// 		Outputs: map[string]ModelOutput{
// 			"best_docked_small_molecule": {
// 				Type: "File",
// 				Glob: []string{"*_docked.sdf", "*_docked.mol2"},
// 			},
// 			"protein": {
// 				Type: "File",
// 				Glob: []string{"*.pdb"},
// 			},
// 		},
// 	}

// 	model, _, err := ReadModelConfig(filePath)
// 	if err != nil {
// 		t.Fatalf("Error reading model config: %v", err)
// 	}

// 	if !reflect.DeepEqual(model, expected) {
// 		t.Errorf("Expected:\n%v\nGot:\n%v", expected, model)
// 	}
// }
