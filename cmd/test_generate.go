package cmd

import (
	"reflect"
	"testing"

	"github.com/labdao/plex/internal/ipwl"
)

func TestGenerateIOGraphFromToolDotProduct(t *testing.T) {
	// Mock the data
	inputVectors := map[string][]string{
		"testkeyA": {"value1A", "value2A"},
		"testkeyB": {"value1B", "value2B"},
	}

	result, err := GenerateIOGraphFromTool("./testTool.json", "dotProduct", inputVectors)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedResult := []ipwl.IO{
		{
			Tool: ipwl.ToolInfo{
				Name: "testTool",
				IPFS: "Qm123",
			},
			Inputs: map[string]ipwl.FileInput{
				"testkeyA": {Class: "File", FilePath: "value1A"},
				"testkeyB": {Class: "File", FilePath: "value1B"},
			},
			Outputs: map[string]ipwl.Output{
				"outputkey1": ipwl.FileOutput{Class: "File", FilePath: "", IPFS: ""},
				"outputkey2": ipwl.FileOutput{Class: "File", FilePath: "", IPFS: ""},
			},
			State:  "created",
			ErrMsg: "",
		},
		{
			Tool: ipwl.ToolInfo{
				Name: "testTool",
				IPFS: "Qm123",
			},
			Inputs: map[string]ipwl.FileInput{
				"testkeyA": {Class: "File", FilePath: "value2A"},
				"testkeyB": {Class: "File", FilePath: "value2B"},
			},
			Outputs: map[string]ipwl.Output{
				"outputkey1": ipwl.FileOutput{Class: "File", FilePath: "", IPFS: ""},
				"outputkey2": ipwl.FileOutput{Class: "File", FilePath: "", IPFS: ""},
			},
			State:  "created",
			ErrMsg: "",
		},
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %v, but got %v", expectedResult, result)
	}
}

func TestGenerateIOGraphFromToolCrossProduct(t *testing.T) {
	// Mock the data
	inputVectors := map[string][]string{
		"testkeyA": {"value1A", "value2A"},
		"testkeyB": {"value1B", "value2B"},
	}

	result, err := GenerateIOGraphFromTool("./testTool.json", "crossProduct", inputVectors)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedResult := []ipwl.IO{
		{
			Tool: ipwl.ToolInfo{
				Name: "testTool",
				IPFS: "Qm123",
			},
			Inputs: map[string]ipwl.FileInput{
				"testkeyA": {Class: "File", FilePath: "value1A"},
				"testkeyB": {Class: "File", FilePath: "value1B"},
			},
			Outputs: map[string]ipwl.Output{
				"outputkey1": ipwl.FileOutput{Class: "File", FilePath: "", IPFS: ""},
				"outputkey2": ipwl.FileOutput{Class: "File", FilePath: "", IPFS: ""},
			},
			State:  "created",
			ErrMsg: "",
		},
		{
			Tool: ipwl.ToolInfo{
				Name: "testTool",
				IPFS: "Qm123",
			},
			Inputs: map[string]ipwl.FileInput{
				"testkeyA": {Class: "File", FilePath: "value1A"},
				"testkeyB": {Class: "File", FilePath: "value2B"},
			},
			Outputs: map[string]ipwl.Output{
				"outputkey1": ipwl.FileOutput{Class: "File", FilePath: "", IPFS: ""},
				"outputkey2": ipwl.FileOutput{Class: "File", FilePath: "", IPFS: ""},
			},
			State:  "created",
			ErrMsg: "",
		},
		{
			Tool: ipwl.ToolInfo{
				Name: "testTool",
				IPFS: "Qm123",
			},
			Inputs: map[string]ipwl.FileInput{
				"testkeyA": {Class: "File", FilePath: "value2A"},
				"testkeyB": {Class: "File", FilePath: "value1B"},
			},
			Outputs: map[string]ipwl.Output{
				"outputkey1": ipwl.FileOutput{Class: "File", FilePath: "", IPFS: ""},
				"outputkey2": ipwl.FileOutput{Class: "File", FilePath: "", IPFS: ""},
			},
			State:  "created",
			ErrMsg: "",
		},
		{
			Tool: ipwl.ToolInfo{
				Name: "testTool",
				IPFS: "Qm123",
			},
			Inputs: map[string]ipwl.FileInput{
				"testkeyA": {Class: "File", FilePath: "value2A"},
				"testkeyB": {Class: "File", FilePath: "value2B"},
			},
			Outputs: map[string]ipwl.Output{
				"outputkey1": ipwl.FileOutput{Class: "File", FilePath: "", IPFS: ""},
				"outputkey2": ipwl.FileOutput{Class: "File", FilePath: "", IPFS: ""},
			},
			State:  "created",
			ErrMsg: "",
		},
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %v, but got %v", expectedResult, result)
	}
}
