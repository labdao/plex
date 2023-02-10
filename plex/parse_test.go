package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCreateInstruction(t *testing.T) {
	want := Instruction{
		App:       "simpdock",
		InputCIDs: []string{"bafybeifzg6egpgb6wi47cayzlltjcdlglls7qtteuqzbrecpiyzyfipuzi"},
		Container: "simpdock:v1",
		Params:    map[string]string{"layers": "33", "steps": "9000", "scifimode": "Y"},
		Cmd:       "python -m inference -l 33 -s 9000 && python -m run --scifimode Y",
		CmdHelper: false,
	}
	type Instruction struct {
		App       string            `json:"app"`
		InputCIDs []string          `json:"input_cids"`
		Container string            `json:"container"`
		Params    map[string]string `json:"params"`
		Cmd       string            `json:"cmd"`
		CmdHelper bool              `json:"cmd_helper"`
	}
	got, err := CreateInstruction("simpdock", "testdata/test_instruction_template.jsonl", "test-directory/", map[string]string{"steps": "9000", "scifimode": "Y"})
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(got), fmt.Sprint(want))
	}
}

func TestOverwriteParams(t *testing.T) {
	defaultParams := map[string]string{"iterations": "42", "repeats": "32", "batch_size": "12"}
	overrideParams := map[string]string{"iterations": "22", "batch_size": "16"}
	want := map[string]string{"iterations": "22", "repeats": "32", "batch_size": "16"}
	got := overwriteParams(defaultParams, overrideParams)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(got), fmt.Sprint(want))
	}
}

func TestFormatCmd(t *testing.T) {
	want := "python -m solvescience --iterations 42 -extra-fast YESSS --batch_size 12"
	unformmatted := "python -m solvescience --iterations %{iterations}s -extra-fast %{fast}s --batch_size %{batch_size}s"
	params := map[string]string{"fast": "YESSS", "batch_size": "12", "iterations": "42"}
	got := formatCmd(unformmatted, params)
	if want != got {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(got), fmt.Sprint(want))
	}
}

func TestCreateInputCID(t *testing.T) {
	want := "bafybeiafphev4hxg5potnih4add2l55qpigjtynyhs7dof4jtbx6fylgei"
	got, err := createInputCID("testdata", true, "python -m desci --decent-lvl 11")
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	if want != got {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(got), fmt.Sprint(want))
	}
}
