package bacalhau

import (
	"fmt"
	"testing"

	"github.com/labdao/plex/internal/ipfs"
)

func DonotTestCreateBacalhauJob(t *testing.T) {
	cid, err := ipfs.WrapAndPinFile("./testdata/design/insulin_target.pdb")
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}

	inputs := map[string]interface{}{
		"binder_length":        50,
		"contigs_override":     "",
		"hotspot":              "",
		"number_of_binders":    2,
		"target_chain":         "D",
		"target_end_residue":   200,
		"target_protein":       cid + "/insulin_target.pdb",
		"target_start_residue": 50,
	}

	container := "ubuntu"
	selector := ""
	maxTime := 60
	memory := 12
	cpu := 1.2
	gpu := true
	network := true
	annotations := []string{"labdaolocal"}
	flowUUID := "1234"
	jobUUID := "5678"
	job, err := CreateBacalhauJob(inputs, container, selector, maxTime, memory, cpu, gpu, network, annotations, flowUUID, jobUUID)
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	if job.Spec.Timeout != int64(maxTime*60) {
		t.Errorf("got = %d; wanted %d", job.Spec.Timeout, maxTime)
	}
	if job.Spec.Resources.Memory != fmt.Sprintf("%dgb", memory) {
		t.Errorf("got = %s; wanted %s", job.Spec.Resources.Memory, fmt.Sprintf("%dgb", memory))
	}
	if job.Spec.Resources.CPU != fmt.Sprintf("%f", cpu) {
		t.Errorf("got = %s; wanted %s", job.Spec.Resources.CPU, fmt.Sprintf("%f", cpu))
	}
	if job.Spec.Resources.GPU != "1" {
		t.Errorf("got = %s; wanted 1", job.Spec.Resources.GPU)
	}
	if job.Spec.EngineSpec.Type != "docker" {
		t.Errorf("got = %s; wanted docker", job.Spec.EngineSpec.Type)
	}
}
