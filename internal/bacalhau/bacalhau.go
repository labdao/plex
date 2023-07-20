package bacalhau

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bacalhau-project/bacalhau/pkg/downloader"
	"github.com/bacalhau-project/bacalhau/pkg/downloader/util"
	"github.com/bacalhau-project/bacalhau/pkg/model"
	"github.com/bacalhau-project/bacalhau/pkg/requester/publicapi"
	"github.com/bacalhau-project/bacalhau/pkg/system"
	"k8s.io/apimachinery/pkg/selection"
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

func CreateBacalhauJob(cid, container, cmd string, memory int, gpu, network bool, annotations []string) (job *model.Job, err error) {
	job, err = model.NewJobWithSaneProductionDefaults()
	if err != nil {
		return nil, err
	}
	job.Spec.Engine = model.EngineDocker
	job.Spec.Docker.Image = container
	job.Spec.Publisher = model.PublisherIpfs
	job.Spec.Docker.Entrypoint = []string{"/bin/bash", "-c", cmd}
	job.Spec.Annotations = annotations

	// had problems getting selector to work in bacalhau v0.28
	var selectorLabel string
	plexEnv, _ := os.LookupEnv("PLEX_ENV")
	if plexEnv == "stage" {
		selectorLabel = "labdaostage"
	} else {
		selectorLabel = "labdao"
	}
	selector := model.LabelSelectorRequirement{Key: "owner", Operator: selection.Equals, Values: []string{selectorLabel}}
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

func CreateBacalhauClient() *publicapi.RequesterAPIClient {
	system.InitConfig()
	apiPort := uint16(1234)
	apiHost := GetBacalhauApiHost()
	client := publicapi.NewRequesterAPIClient(apiHost, apiPort)
	return client
}

func SubmitBacalhauJob(job *model.Job) (submittedJob *model.Job, err error) {
	client := CreateBacalhauClient()
	submittedJob, err = client.Submit(context.Background(), job)
	return submittedJob, err
}

func GetBacalhauJobResults(submittedJob *model.Job, showAnimation bool) (results []model.PublishedResult, err error) {
	client := CreateBacalhauClient()
	maxTrys := 360 // 30 minutes divided by 5 seconds is 360 iterations
	animation := []string{"\U0001F331", "_", "_", "_", "_"}
	fmt.Println("Job running...")

	fmt.Printf("Bacalhau job id: %s \n", submittedJob.Metadata.ID)

	for i := 0; i < maxTrys; i++ {
		// mcmenemy check for status first
		updatedJob, _, err := client.Get(context.Background(), submittedJob.Metadata.ID)
		if err != nil {
			return results, err
		}
		if updatedJob.State.State == model.JobStateCancelled {
			return results, fmt.Errorf("bacalhau cancelled job")
		} else if updatedJob.State.State == model.JobStateError {
			return results, fmt.Errorf("bacalhau errored job")
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
		time.Sleep(2 * time.Second)
	}
	return results, err
}

func DownloadBacalhauResults(dir string, submittedJob *model.Job, results []model.PublishedResult) error {
	downloadSettings := util.NewDownloadSettings()
	downloadSettings.OutputDir = dir
	cm := system.NewCleanupManager()
	downloaderProvider := util.NewStandardDownloaders(cm, downloadSettings)
	err := downloader.DownloadResults(context.Background(), results, downloaderProvider, downloadSettings)
	return err
}
