// package docker

// import (
// 	"fmt"
// 	"os"
// 	"path"
// 	"testing"
// )

// func TestInstructionToDockerCmd(t *testing.T) {
// 	want := "docker run --gpus -v home/job-dir:/inputs -v home/job-dir/outputs:/outputs mycontainer /bin/bash -c 'python -m molbind'"
// 	got := InstructionToDockerCmd("mycontainer", "python -m molbind", "home/job-dir", true)
// 	if want != got {
// 		t.Errorf("got = %s; wanted %s", fmt.Sprint(got), fmt.Sprint(want))
// 	}
// }

// func TestRunDockerJob(t *testing.T) {
// 	dir, err := os.Getwd()
// 	if err != nil {
// 		t.Errorf(fmt.Sprint(err))
// 	}

// 	jobDir := path.Join(dir, "testJobDir")
// 	err = RunDockerJob("ubuntu", "ls", jobDir, false)
// 	if err != nil {
// 		t.Errorf(fmt.Sprint(err))
// 	}
// }
