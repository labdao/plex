package main

import (
	"fmt"
	"testing"
)

func TestInstructionToBacalhauCmd(t *testing.T) {
	want := "bacalhau docker run --network full --gpu 1 --memory 12gb -i QmZGavZu mycontainer -- python -m molbind"
	got := InstructionToBacalhauCmd("QmZGavZu", "mycontainer", "python -m molbind")
	if want != got {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(got), fmt.Sprint(want))
	}
}
