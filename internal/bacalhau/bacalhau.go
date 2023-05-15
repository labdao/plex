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

func CreateBacalhauJob(cids []string, container, cmd string, memory int, gpu, network bool) (job *model.Job, err error) {
	job, err = model.NewJobWithSaneProductionDefaults()
	if err != nil {
		return nil, err
	}
	job.Spec.Engine = model.EngineDocker
	job.Spec.Docker.Image = container
	job.Spec.Publisher = model.PublisherIpfs
	job.Spec.Docker.Entrypoint = []string{"/bin/bash", "-c", cmd}

	// had problems getting selector to work in bacalhau v0.28
	// var selectorLabel string
	// plexEnv, _ := os.LookupEnv("PLEX_ENV")
	// if plexEnv == "stage" {
	// 	selectorLabel = "labdaostage"
	// } else {
	// 	selectorLabel = "labdao"
	// }
	// job.Spec.NodeSelectors = []model.LabelSelectorRequirement{selector}

	if memory > 0 {
		job.Spec.Resources.Memory = fmt.Sprintf("%dgb", memory)
	}
	if gpu {
		job.Spec.Resources.GPU = "1"
	}
	if network {
		job.Spec.Network = model.NetworkConfig{Type: model.NetworkFull}
	}

	// make a for loop that iterates through the cids and adds each to the job
	for index, cid := range cids {
		inputPath := fmt.Sprintf("/inputs/%d", index)
		job.Spec.Inputs = append(job.Spec.Inputs, model.StorageSpec{StorageSource: model.StorageSourceIPFS, CID: cid, Path: inputPath})
	}

	// job.Spec.Inputs = []model.StorageSpec{{StorageSource: model.StorageSourceIPFS, CID: cid, Path: "/inputs"}}
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

func GetBacalhauJobResults(submittedJob *model.Job) (results []model.PublishedResult, err error) {
	client := CreateBacalhauClient()
	maxTrys := 360 // 30 minutes divided by 5 seconds is 360 iterations
	animation := []string{"\U0001F331", "_", "_", "_", "_"}
	fmt.Println("Job running...")

	for i := 0; i < maxTrys; i++ {
		saplingIndex := i % 5

		results, err = client.GetResults(context.Background(), submittedJob.Metadata.ID)
		if err != nil {
			return results, err
		}
		if len(results) > 0 {
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
	err := downloader.DownloadResults(context.Background(), results, downloaderProvider, downloadSettings)
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
