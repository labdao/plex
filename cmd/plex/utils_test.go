// package plex

// import (
// 	"testing"
// )

// func TestFindAppConfig(t *testing.T) {
// 	// it errors for a nonexistent filepath
// 	_, err := FindAppConfig("diffdock", "../../testdata/nonexistent.jsonl")
// 	if err == nil {
// 		t.Errorf("findAppConfig should error for nonexsitent jsonl filepath")
// 	}

// 	// it errors for invalid jsonl
// 	_, err = FindAppConfig("diffdock", "../../testdata/invalid_app.jsonl")
// 	if err == nil {
// 		t.Errorf("findAppConfig should error for invald jsonl")
// 	}

// 	// it errors if app doesn't exist
// 	_, err = FindAppConfig("theranos", "../../testdata/invalid_app.jsonl")
// 	if err == nil {
// 		t.Errorf("findAppConfig should error for nonexistent app")
// 	}

// 	// it finds valid app
// 	appConfig, err := FindAppConfig("diffdock", "../../testdata/valid_app.jsonl")
// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}
// 	if appConfig.App != "diffdock" {
// 		t.Errorf("got = %s; wanted diffdock", appConfig.App)
// 	}
// }
