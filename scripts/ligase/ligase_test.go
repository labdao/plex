package main

import (
	"testing"
)

func TestLoadToolsConfig(t *testing.T) {
	filePath := "testdata/tools.json"
	toolsConfig, err := loadToolsConfig(filePath)
	if err != nil {
		t.Errorf("loadToolsConfig() failed with error: %v", err)
	}

	// Verify that some values in the toolsConfig struct are as expected
	expectedName := "equibind"
	if toolsConfig.Name != expectedName {
		t.Errorf("Expected name to be %s, but got %s", expectedName, toolsConfig.Name)
	}

	expectedDescription := "Docking of small molecules to a protein"
	if toolsConfig.Description != expectedDescription {
		t.Errorf("Expected description to be %s, but got %s", expectedDescription, toolsConfig.Description)
	}

	if len(toolsConfig.Requirements) != 2 {
		t.Errorf("Expected requirements length to be 2, but got %d", len(toolsConfig.Requirements))
	}
}

func TestFindOutputs(t *testing.T) {
	filePath := "testdata/tools.json"
	toolsConfig, err := loadToolsConfig(filePath)
	if err != nil {
		t.Fatalf("loadToolsConfig() failed with error: %v", err)
	}

	dir := "testdata/"
	outputMap, err := findOutputs(dir, toolsConfig)
	if err != nil {
		t.Fatalf("findOutputs() failed with error: %v", err)
	}

	// Verify that the expected keys exist in the output map
	expectedKeys := []string{"best_docked_small_molecule", "all_docked_small_molecules"}
	for _, key := range expectedKeys {
		if _, ok := outputMap[key]; !ok {
			t.Errorf("Expected key %s not found in output map", key)
		}
	}

	// Verify that the correct files are matched for each key
	expectedFileCounts := map[string]int{
		"best_docked_small_molecule": 1,
		"all_docked_small_molecules": 2,
	}
	for key, expectedCount := range expectedFileCounts {
		if output, ok := outputMap[key]; ok {
			if len(output.Address) != expectedCount {
				t.Errorf("Expected %d files for key %s, but got %d", expectedCount, key, len(output.Address))
			}
		}
	}
}
