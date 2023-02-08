package plexus

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCanary(t *testing.T) {
	got := "desci"
	want := "desci"
	if got != want {
		t.Errorf("got = %s; wanted %s", got, want)
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
