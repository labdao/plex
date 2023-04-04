package ipwl

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestUpdateIOState(t *testing.T) {
	// Create a temporary copy of the original JSON file
	origFilePath := "testdata/example_io.json"
	tempFilePath := "testdata/temp_io_list.json"

	input, err := ioutil.ReadFile(origFilePath)
	if err != nil {
		t.Fatalf("Failed to read original file: %v", err)
	}

	err = ioutil.WriteFile(tempFilePath, input, 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	defer func() {
		// Restore the original file after the test is completed
		err = ioutil.WriteFile(origFilePath, input, 0644)
		if err != nil {
			t.Fatalf("Failed to restore original file: %v", err)
		}
		os.Remove(tempFilePath)
	}()

	// Test the updateIOState function
	index := 0
	newState := "testing"

	err = updateIOState(tempFilePath, index, newState)
	if err != nil {
		t.Fatalf("Error updating IO state: %v", err)
	}

	ioList, err := ReadIOList(tempFilePath)
	if err != nil {
		t.Fatalf("Error reading IO list: %v", err)
	}

	if ioList[index].State != newState {
		t.Errorf("Expected state to be '%v', got '%v'", newState, ioList[index].State)
	}
}
