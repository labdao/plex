package bacalhau

import (
	"fmt"
	"testing"
)

func TestCreateBalhauJob(t *testing.T) {
	cid := "bafybeibuivbwkakim3hkgphffaipknuhw4epjfu3bstvsuv577spjhbvju"
	container := "ubuntu"
	cmd := "echo DeSci"
	memory := "12gb"
	gpu := "1"
	job, err := CreateBacalhauJob(cid, container, cmd, 12, true, true)
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	if job.Spec.Resources.Memory != memory {
		t.Errorf("got = %s; wanted %s", job.Spec.Resources.Memory, memory)
	}
	if job.Spec.Resources.GPU != gpu {
		t.Errorf("got = %s; wanted %s", job.Spec.Resources.GPU, gpu)
	}
	if job.Spec.Resources.GPU != gpu {
		t.Errorf("got = %s; wanted %s", job.Spec.Resources.GPU, gpu)
	}
	if job.Spec.Docker.Image != container {
		t.Errorf("got = %s; wanted %s", job.Spec.Docker.Image, container)
	}
	if fmt.Sprint(job.Spec.Docker.Entrypoint) != fmt.Sprint([]string{"/bin/bash", "-c", cmd}) {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(job.Spec.Docker.Entrypoint), fmt.Sprint(cmd))
	}
	if job.Spec.Inputs[0].CID != cid {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(job.Spec.Inputs[0].CID), fmt.Sprint(cid))
	}
}

/*
func TestGetBacalhauJobResults(t *testing.T) {
	cid := "bafybeibuivbwkakim3hkgphffaipknuhw4epjfu3bstvsuv577spjhbvju"
	container := "ubuntu"
	cmd := "printenv && echo DeSci"
	job, err := CreateBacalhauJob(cid, container, cmd, 0, false, false)
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	submittedBacalhauJob, err := SubmitBacalhauJob(job)
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	fmt.Println(submittedBacalhauJob.Metadata.ID)
	results, err := GetBacalhauJobResults(submittedBacalhauJob)
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	if len(results) == 0 {
		t.Errorf("Bacalhau failed to find completed job")
	}
}
*/

func TestInstructionToBacalhauCmd(t *testing.T) {
	want := "bacalhau docker run --selector owner=labdao --memory 4gb --network full -i QmZGavZu mycontainer -- /bin/bash -c 'python -m molbind'"
	got := InstructionToBacalhauCmd("QmZGavZu", "mycontainer", "python -m molbind", 4, false, true)
	if want != got {
		t.Errorf("got = %s; wanted %s", fmt.Sprint(got), fmt.Sprint(want))
	}
}
