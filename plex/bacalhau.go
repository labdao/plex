package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/filecoin-project/bacalhau/pkg/model"
)

func InstructionToBacalhauCmd(cid, container, cmd string) string {
	// TODO allow overrides for gpu memory and network flags
	return `bacalhau docker run --network full --gpu 1 --memory 12gb -i ` + fmt.Sprintf(cid) + ` ` + fmt.Sprintf(container) + ` -- ` + fmt.Sprintf(cmd)
}

func createBacalhauJob(cid, container, cmd string) (job *model.Job, err error) {
	job, err = model.NewJobWithSaneProductionDefaults()
	if err != nil {
		return nil, err
	}
	job.Spec.Docker.Image = container
	job.Spec.Docker.Entrypoint = []string{cmd}
	job.Spec.Network = model.NetworkConfig{Type: model.NetworkFull}
	job.Spec.Resources.Memory = "12gb"
	job.Spec.Resources.GPU = "1"
	job.Spec.Inputs = []model.StorageSpec{{CID: cid}}
	return job, err
}

func RunBacalhauCmd(cmdString string) (out []byte, err error) {
	args := strings.Fields(cmdString)
	fmt.Println(args)
	cmd := exec.Command(args[0], args[1:]...)
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Command failed: %s\n", err)
		return
	}
	fmt.Printf("Output: %s\n", out)
	return
}
