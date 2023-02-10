package main

import (
	"os"
	"testing"
)

func TestValidateAppConfig(t *testing.T) {
	// Test cases
	testCases := []struct {
		appConfig string
		want      bool
	}{
		{"testdata/invalid_app.jsonl", false},
		{"testdata/valid_app.jsonl", true},
		{"nonexistent.jsonl", false},
		{"app.jsonl", true},
	}

	// Test each test case
	for _, tc := range testCases {
		appConfig := tc.appConfig
		want := tc.want

		// Call the function to test
		got, _ := validateAppConfig(appConfig)

		// Assert the expected result
		if got != want {
			t.Errorf("validateAppConfig(%q) = %v, want %v", appConfig, got, want)
		}
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
		got, _ := validateDirectoryPath(directory)

		// Assert the expected result
		if got != want {
			t.Errorf("validateDirectoryPath(%q) = %v, want %v", directory, got, want)
		}
	}
}
