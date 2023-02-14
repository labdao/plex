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

func TestCreatevalhauJob(t *testing.T) {
	cid := "bafybeig7rsafgrtwzivrorumixcqxpwmje7cp56eoxzg3jbwxxyy26xgue"
	container := "ubuntu"
	cmd := "echo DeSci"
	job, err := createBacalhauJob(cid, container, cmd)
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	if job.Spec.Docker.Image != container {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(job.Spec.Docker.Image), fmt.Sprint(container))
	}
	if job.Spec.Docker.Entrypoint[0] != cmd {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(job.Spec.Docker.Entrypoint), fmt.Sprint(cmd))
	}
	if job.Spec.Inputs[0].CID != cid {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(job.Spec.Inputs[0].CID), fmt.Sprint(cid))
	}
}

func TestRunBacalhauCmd(t *testing.T) {
	cmd := "bacalhau docker run ubuntu echo Hello World"
	out, err := RunBacalhauCmd(cmd)
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	fmt.Printf("Output: %s\n", out)
}
