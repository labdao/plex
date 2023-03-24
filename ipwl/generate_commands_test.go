package main

import (
    "io/ioutil"
    "os"
    "testing"
)

func TestReadInputFile(t *testing.T) {
    inputJSON := []byte(`
        {
            "protein": {
                "class": "File",
                "basename": "7n9g.pdb",
                "size": 231008,
                "path": "7n9g.pdb"
            },
            "small_molecule": {
                "class": "File",
                "basename": "ZINC000003986735.sdf",
                "size": 732,
                "path": "ZINC000003986735.sdf"
            }
        }
    `)
    tempFile, err := ioutil.TempFile("", "test_input*.json")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tempFile.Name())
    if _, err := tempFile.Write(inputJSON); err != nil {
        t.Fatal(err)
    }
    if err := tempFile.Close(); err != nil {
        t.Fatal(err)
    }

    input, err := readInputFile(tempFile.Name())
    if err != nil {
        t.Fatal(err)
    }

    if _, ok := (*input)["protein"].(map[string]interface{}); !ok {
        t.Errorf("expected protein input to be a map[string]interface{}, got %T", (*input)["protein"])
    }

    if _, ok := (*input)["small_molecule"].(map[string]interface{}); !ok {
        t.Errorf("expected small_molecule input to be a map[string]interface{}, got %T", (*input)["small_molecule"])
    }
}