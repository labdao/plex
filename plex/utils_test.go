package main

import (
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
