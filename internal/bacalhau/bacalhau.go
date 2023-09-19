package bacalhau

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bacalhau-project/bacalhau/pkg/downloader"
	"github.com/bacalhau-project/bacalhau/pkg/downloader/util"
	node "github.com/bacalhau-project/bacalhau/pkg/job"
	"github.com/bacalhau-project/bacalhau/pkg/model"
	"github.com/bacalhau-project/bacalhau/pkg/requester/publicapi"
	"github.com/bacalhau-project/bacalhau/pkg/system"
)

func GetBacalhauApiHost() string {
	bacalApiHost, exists := os.LookupEnv("BACALHAU_API_HOST")
	plexEnv, _ := os.LookupEnv("PLEX_ENV")
	if exists {
		return bacalApiHost
	} else if plexEnv == "stage" {
		return "44.198.42.30"
	} else {
		return "54.210.19.52"
	}
}

func CreateBacalhauJob(cid, container, cmd, selector string, maxTime, memory int, gpu, network bool, annotations []string) (job *model.Job, err error) {
	job, err = model.NewJobWithSaneProductionDefaults()
	if err != nil {
		return nil, err
	}
	job.Spec.Engine = model.EngineDocker
	job.Spec.Docker.Image = container
	job.Spec.Publisher = model.PublisherIpfs
	job.Spec.Docker.Entrypoint = []string{"/bin/bash", "-c", cmd}
	job.Spec.Annotations = annotations
	job.Spec.Timeout = float64(maxTime * 60)

	plexEnv, _ := os.LookupEnv("PLEX_ENV")
	if selector == "" && plexEnv == "stage" {
		selector = "owner=labdaostage"
	} else if selector == "" && plexEnv == "prod" {
		selector = "owner=labdao"
	}
	nodeSelectorRequirements, err := node.ParseNodeSelector(selector)
	if err != nil {
		return nil, err
	}
	job.Spec.NodeSelectors = nodeSelectorRequirements

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

func CreateBacalhauClient() *publicapi.RequesterAPIClient {
	system.InitConfig()
	apiHost := GetBacalhauApiHost()
	apiPort := uint16(1234)
	client := publicapi.NewRequesterAPIClient(apiHost, apiPort)
	return client
}

func SubmitBacalhauJob(job *model.Job) (submittedJob *model.Job, err error) {
	client := CreateBacalhauClient()
	submittedJob, err = client.Submit(context.Background(), job)
	return submittedJob, err
}

func GetBacalhauJobResults(submittedJob *model.Job, showAnimation bool, maxTime int) (results []model.PublishedResult, err error) {
	client := CreateBacalhauClient()

	sleepConstant := 2
	maxTrys := maxTime * 60 / sleepConstant

	animation := []string{"\U0001F331", "_", "_", "_", "_"}
	fmt.Println("Job running...")

	fmt.Printf("Bacalhau job id: %s \n", submittedJob.Metadata.ID)

	for i := 0; i < maxTrys; i++ {
		updatedJob, _, err := client.Get(context.Background(), submittedJob.Metadata.ID)
		if err != nil {
			return results, err
		}
		if i == maxTrys-1 {
			return results, fmt.Errorf("bacalhau job did not finish within the expected time (~%d min); please check the job status manually with `bacalhau describe %s`", maxTime, submittedJob.Metadata.ID)
		}
		if updatedJob.State.State == model.JobStateCancelled {
			return results, fmt.Errorf("bacalhau cancelled job; please run `bacalhau describe %s` for more details", submittedJob.Metadata.ID)
		} else if updatedJob.State.State == model.JobStateError {
			return results, fmt.Errorf("bacalhau errored job; please run `bacalhau describe %s` for more details", submittedJob.Metadata.ID)
		} else if updatedJob.State.State == model.JobStateCompleted {
			results, err = client.GetResults(context.Background(), submittedJob.Metadata.ID)
			if err != nil {
				return results, err
			}
			if len(results) > 0 {
				return results, err
			} else {
				return results, fmt.Errorf("bacalhau job completed but no results found")
			}
		}
		if showAnimation {
			saplingIndex := i % 5
			animation[saplingIndex] = "\U0001F331"
			fmt.Printf("////%s////\r", strings.Join(animation, ""))
			animation[saplingIndex] = "_"
		}
		time.Sleep(time.Duration(sleepConstant) * time.Second)
	}
	return results, err
}

func DownloadBacalhauResults(dir string, submittedJob *model.Job, results []model.PublishedResult) error {
	cm := system.NewCleanupManager()
	downloadSettings := &model.DownloaderSettings{
		Timeout:   50 * time.Second,
		OutputDir: dir,
	}
	downloadSettings.OutputDir = dir
	downloaderProvider := util.NewStandardDownloaders(cm, downloadSettings)
	err := downloader.DownloadResults(context.Background(), results, downloaderProvider, downloadSettings)
	return err
}
