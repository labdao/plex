package plex

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCreateInstruction(t *testing.T) {
	want := Instruction{
		Tool:       "simpdock",
		InputCIDs: []string{"QmWVKoVYBWHWdRLrL8Td5kUpqN2qH6zQ5piwtdCE1fjSYt", "QmAnotherCIDHere123456789"},
		Container: "simpdock:v1",
		Cmd:       "python -m inference -l 33 -s 9000 && python -m run --protein /inputs/7n9g.pdb --small_molecule_library /inputs/ZINC000003986735.sdf --scifimode Y",
	}
	got, err := CreateInstruction("simpdock", "../../testdata/simpdock.json", "../../testdata/simpdock-input.json")
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(got), fmt.Sprint(want))
	}
}

// TODO not sure we need this function and test anymore
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
	want := "bafybeifzg6egpgb6wi47cayzlltjcdlglls7qtteuqzbrecpiyzyfipuzi"
	got, err := CreateInputCID("ipfs_test", "python -m desci --decent-lvl 11")
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	if want != got {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(got), fmt.Sprint(want))
	}
}
