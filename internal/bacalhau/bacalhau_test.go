package bacalhau

import (
	"fmt"
	"testing"
)

func TestCreateBalhauJob(t *testing.T) {
	cid := "bafybeibuivbwkakim3hkgphffaipknuhw4epjfu3bstvsuv577spjhbvju"
	container := "ubuntu"
	cmd := "echo DeSci"
	maxTime := 60
	timeOut := maxTime * 60 // Bacalhau timeout is in seconds, so we need to multiply by 60
	memory := "12gb"
	gpu := "1"
	networkFlag := true
	job, err := CreateBacalhauJob(cid, container, cmd, 60, 12, true, networkFlag, []string{})
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	if job.Spec.Timeout != float64(timeOut) {
		t.Errorf("got = %f; wanted %d", job.Spec.Timeout, timeOut)
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
