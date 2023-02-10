package main

import (
	"fmt"
	"testing"
)

func TestInstructionToBacalhauCmd(t *testing.T) {
	// it uses cmd when cmdHelper is false
	want := "bacalhau docker run --network full --gpu 1 --memory 12gb -i QmZGavZu mycontainer python -m molbind"
	got := InstructionToBacalhauCmd("QmZGavZu", "mycontainer", "python -m molbind", false)
	if want != got {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(got), fmt.Sprint(want))
	}

	// it uses helper.sh cmd when cmdHelper is true
	want = "bacalhau docker run --network full --gpu 1 --memory 12gb -i QmZGavZu mycontainer ./helper.sh"
	got = InstructionToBacalhauCmd("QmZGavZu", "mycontainer", "python -m molbind", true)
	if want != got {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(got), fmt.Sprint(want))
	}
}
