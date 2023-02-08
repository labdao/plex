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
