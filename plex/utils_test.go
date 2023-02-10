package main

import (
	"os"
	"testing"
)

func TestFindAppConfig(t *testing.T) {
	// it errors for a nonexistent filepath
	_, err := findAppConfig("diffdock", "nonexistent.jsonl")
	if err == nil {
		t.Errorf("findAppConfig should error for nonexsitent jsonl filepath")
	}

	// it errors for invalid jsonl
	_, err = findAppConfig("diffdock", "testdata/invalid_app.jsonl")
	if err == nil {
		t.Errorf("findAppConfig should error for invald jsonl")
	}

	// it errors if app doesn't exist
	_, err = findAppConfig("theranos", "testdata/invalid_app.jsonl")
	if err == nil {
		t.Errorf("findAppConfig should error for nonexistent app")
	}

	// it finds valid app
	appConfig, err := findAppConfig("diffdock", "testdata/valid_app.jsonl")
	if err != nil {
		t.Errorf(err.Error())
	}
	if appConfig.App != "diffdock" {
		t.Errorf("got = %s; wanted diffdock", appConfig.App)
	}
}

func TestValidateDirectoryPath(t *testing.T) {
	// Test cases
	got, _ := os.Getwd()
	testCases := []struct {
		directory string
		want      bool
	}{
		{got, true},
		{"/nonexistent", false},
		{"/etc/passwd", false},
	}

	// Test each test case
	for _, tc := range testCases {
		directory := tc.directory
		want := tc.want

		// Call the function to test
		got, _ := ValidateDirectoryPath(directory)

		// Assert the expected result
		if got != want {
			t.Errorf("validateDirectoryPath(%q) = %v, want %v", directory, got, want)
		}
	}
}
