package bacalhau

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/filecoin-project/bacalhau/pkg/downloader"
	"github.com/filecoin-project/bacalhau/pkg/downloader/util"
	"github.com/filecoin-project/bacalhau/pkg/model"
	"github.com/filecoin-project/bacalhau/pkg/requester/publicapi"
	"github.com/filecoin-project/bacalhau/pkg/system"
	"k8s.io/apimachinery/pkg/selection"
)

func CreateBacalhauJob(cid, container, cmd string, memory int, gpu, network bool) (job *model.Job, err error) {
	job, err = model.NewJobWithSaneProductionDefaults()
	if err != nil {
		return nil, err
	}
	job.Spec.Engine = model.EngineDocker
	job.Spec.Docker.Image = container
	job.Spec.Docker.Entrypoint = []string{"/bin/bash", "-c", cmd}
	selector := model.LabelSelectorRequirement{Key: "owner", Operator: selection.Equals, Values: []string{"labdao"}}
	job.Spec.NodeSelectors = []model.LabelSelectorRequirement{selector}
	if memory > 0 {
		job.Spec.Resources.Memory = fmt.Sprintf("%dgb", memory)
	}
	if gpu {
		job.Spec.Resources.GPU = "1"
	}
	if network {
		job.Spec.Network = model.NetworkConfig{Type: model.NetworkFull}
	}
	job.Spec.Inputs = []model.StorageSpec{{StorageSource: model.StorageSourceIPFS, CID: cid, Path: "/inputs"}}
	job.Spec.Outputs = []model.StorageSpec{{Name: "outputs", StorageSource: model.StorageSourceIPFS, Path: "/outputs"}}
	return job, err
}

func SubmitBacalhauJob(job *model.Job) (submittedJob *model.Job, err error) {
	system.InitConfig()
	apiPort := 1234
	apiHost := "35.245.115.191"
	client := publicapi.NewRequesterAPIClient(fmt.Sprintf("http://%s:%d", apiHost, apiPort))
	submittedJob, err = client.Submit(context.Background(), job)
	return submittedJob, err
}

func GetBacalhauJobResults(submittedJob *model.Job) (results []model.PublishedResult, err error) {
	system.InitConfig()
	apiPort := 1234
	apiHost := "35.245.115.191"
	client := publicapi.NewRequesterAPIClient(fmt.Sprintf("http://%s:%d", apiHost, apiPort))
	maxTrys := 360 // 30 minutes divided by 5 seconds is 360 iterations
	animation := []string{"\U0001F331", "_", "_", "_", "_"}
	fmt.Println("Job running...")
	for i := 0; i < maxTrys; i++ {
		saplingIndex := i % 5
		jobState, err := client.GetJobState(context.Background(), submittedJob.Metadata.ID)
		if err != nil {
			return results, err
		}

		// check once to see if job is already complete
		results, _ = client.GetResults(context.Background(), submittedJob.Metadata.ID)
		if len(results) > 0 {
			return results, err
		}

		// check to see if any node shards have finished or errored while running job
		// this assumes the job spec will only have one shard attempt to run the job
		completedShardRuns := []model.JobShardState{}
		erroredShardRuns := []model.JobShardState{}
		cancelledShardRuns := []model.JobShardState{}
		for _, jobNodeState := range jobState.Nodes {
			for _, jobShardState := range jobNodeState.Shards {
				if jobShardState.State == model.JobStateCompleted {
					completedShardRuns = append(completedShardRuns, jobShardState)
				} else if jobShardState.State == model.JobStateError {
					erroredShardRuns = append(erroredShardRuns, jobShardState)
				} else if jobShardState.State == model.JobStateCancelled {
					cancelledShardRuns = append(cancelledShardRuns, jobShardState)
				}
			}
		}
		if len(completedShardRuns) > 0 {
			fmt.Println("")
			fmt.Println("\U0001F332 Job run complete")
			results, err = client.GetResults(context.Background(), submittedJob.Metadata.ID)
			return results, err
		}
		if len(cancelledShardRuns) > 0 || len(erroredShardRuns) > 0 {
			fmt.Println("")
			fmt.Println("\U0001F342 Job failed to complete")
			results, err = client.GetResults(context.Background(), submittedJob.Metadata.ID)
			return results, err
		}
		animation[saplingIndex] = "\U0001F331"
		fmt.Printf("////%s////\r", strings.Join(animation, ""))
		animation[saplingIndex] = "_"
		time.Sleep(2 * time.Second)
	}
	return results, err
}

func DownloadBacalhauResults(dir string, submittedJob *model.Job, results []model.PublishedResult) error {
	downloadSettings := util.NewDownloadSettings()
	downloadSettings.OutputDir = dir
	cm := system.NewCleanupManager()
	downloaderProvider := util.NewStandardDownloaders(cm, downloadSettings)
	err := downloader.DownloadJob(context.Background(), submittedJob.Spec.Outputs, results, downloaderProvider, downloadSettings)
	return err
}

func InstructionToBacalhauCmd(cid, container, cmd string, memory int, gpu, network bool) string {
	gpuFlag := ""
	if gpu {
		gpuFlag = "--gpu 1 "
	}
	memoryFlag := ""
	if memory != 0 {
		memoryFlag = fmt.Sprintf("--memory %dgb ", memory)
	}
	networkFlag := ""
	if network {
		networkFlag = "--network full"
	}
	return fmt.Sprintf("bacalhau docker run --selector owner=labdao %s%s%s -i %s %s -- /bin/bash -c '%s'", gpuFlag, memoryFlag, networkFlag, cid, container, cmd)
}

func RunBacalhauCmd(cmdString string) {
	args := strings.Fields(cmdString)
	fmt.Println(args)
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Command failed: %s\n", err)
		return
	}
	fmt.Printf("Output: %s\n", out)
}
