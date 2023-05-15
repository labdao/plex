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
