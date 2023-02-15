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

func TestCreateBalhauJob(t *testing.T) {
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
	job, err := createBacalhauJob(cid, container, cmd)
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

/*
func TestRunBacalhauCmd(t *testing.T) {
	cmd := "bacalhau docker run ubuntu echo Hello World"
	out, err := RunBacalhauCmd(cmd)
	if err != nil {
		t.Fatalf(fmt.Sprint(err))
	}
	fmt.Printf("Output: %s\n", out)
}
*/

/*
    // Loop until the job is no longer in progress
    for {
        // Get the status of the job using the model's JobStatus method
        status, err := client.GetJobStatus(jobID)
        if err != nil {
            panic(err)
        }

        fmt.Printf("Job Status: %v\n", status.State)

        // If the job has completed or failed, break out of the loop
        if status.State == model.JobStateCompleted || status.State == model.JobStateFailed {
            break
        }

        // Wait for a few seconds before checking the job status again
        time.Sleep(5 * time.Second)
    }
}
*/
