package main

import (
	"fmt"
	"testing"
)

func TestCreateBalhauJob(t *testing.T) {
	cid := "bafybeig7rsafgrtwzivrorumixcqxpwmje7cp56eoxzg3jbwxxyy26xgue"
	container := "ubuntu"
	cmd := "echo DeSci"
	job, err := createBacalhauJob(cid, container, cmd, "true")
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	if job.Spec.Docker.Image != container {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(job.Spec.Docker.Image), fmt.Sprint(container))
	}
	if fmt.Sprint(job.Spec.Docker.Entrypoint) != fmt.Sprint([]string{"/bin/bash", "-c", cmd}) {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(job.Spec.Docker.Entrypoint), fmt.Sprint(cmd))
	}
	if job.Spec.Inputs[0].CID != cid {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(job.Spec.Inputs[0].CID), fmt.Sprint(cid))
	}
}

func TestGetBacalhauJobResults(t *testing.T) {
	cid := "bafybeig7rsafgrtwzivrorumixcqxpwmje7cp56eoxzg3jbwxxyy26xgue"
	container := "ubuntu"
	cmd := "printenv && echo DeSci"
	job, err := createBacalhauJob(cid, container, cmd, "true")
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	submittedBacalhauJob, err := submitBacalhauJob(job)
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	fmt.Println(submittedBacalhauJob.Metadata.ID)
	results, err := getBacalhauJobResults(submittedBacalhauJob)
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	if len(results) == 0 {
		t.Errorf("Bacalhau failed to find completed job")
	}
}

func TestInstructionToBacalhauCmd(t *testing.T) {
	want := "bacalhau docker run --network full --gpu 1 --memory 12gb -i QmZGavZu mycontainer -- python -m molbind"
	got := InstructionToBacalhauCmd("QmZGavZu", "mycontainer", "python -m molbind", "true")
	if want != got {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(got), fmt.Sprint(want))
	}
}
