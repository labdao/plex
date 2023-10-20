package bacalhau

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bacalhau-project/bacalhau/cmd/util/parse"
	"github.com/bacalhau-project/bacalhau/pkg/config"
	"github.com/bacalhau-project/bacalhau/pkg/downloader"
	"github.com/bacalhau-project/bacalhau/pkg/downloader/util"
	"github.com/bacalhau-project/bacalhau/pkg/model"
	"github.com/bacalhau-project/bacalhau/pkg/publicapi/client"
	"github.com/bacalhau-project/bacalhau/pkg/system"
)

func GetBacalhauApiHost() string {
	bacalApiHost, exists := os.LookupEnv("BACALHAU_API_HOST")
	plexEnv, _ := os.LookupEnv("PLEX_ENV")
	if exists {
		return bacalApiHost
	} else if plexEnv == "stage" {
		return "bacalhau.staging.labdao.xyz"
	} else {
		return "bacalhau.labdao.xyz"
	}
}

func CreateBacalhauJobV2(inputs map[string]string, container, selector string, cmd []string, maxTime, memory int, gpu, network bool, annotations []string) (job *model.Job, err error) {
	log.Println("Creating job inside v2 function")
	job, err = model.NewJobWithSaneProductionDefaults()
	if err != nil {
		return nil, err
	}
	fmt.Println("container cmd", cmd)
	job.Spec.EngineSpec = model.NewDockerEngineBuilder(container).
		WithEntrypoint(cmd...).Build()
	job.Spec.PublisherSpec = model.PublisherSpec{
		Type: model.PublisherIpfs,
	}
	job.Spec.Annotations = annotations
	job.Spec.Timeout = int64(maxTime * 60)

	nodeSelectorRequirements, err := parse.NodeSelector(selector)
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
	job.Spec.Inputs = []model.StorageSpec{}
	for key, cid := range inputs {
		job.Spec.Inputs = append(job.Spec.Inputs,
			// ToDo for arrays split by comma and put inside a dir that is key/index
			model.StorageSpec{
				StorageSource: model.StorageSourceIPFS,
				CID:           cid,
				Path:          "/inputs/" + key,
			})
	}
	job.Spec.Outputs = []model.StorageSpec{{Name: "outputs", StorageSource: model.StorageSourceIPFS, Path: "/outputs"}}
	log.Println("returning job")
	return job, err
}

func CreateBacalhauJob(cid, container, cmd, selector string, maxTime, memory int, gpu, network bool, annotations []string) (job *model.Job, err error) {
	job, err = model.NewJobWithSaneProductionDefaults()
	if err != nil {
		return nil, err
	}
	job.Spec.Engine = model.EngineDocker
	job.Spec.Docker.Image = container
	job.Spec.PublisherSpec = model.PublisherSpec{
		Type: model.PublisherIpfs,
	}
	job.Spec.Docker.Entrypoint = []string{"/bin/bash", "-c", cmd}
	job.Spec.Annotations = annotations
	job.Spec.Timeout = int64(maxTime * 60)

	plexEnv, _ := os.LookupEnv("PLEX_ENV")
	if selector == "" && plexEnv == "stage" {
		selector = "owner=labdaostage"
	} else if selector == "" && plexEnv == "prod" {
		selector = "owner=labdao"
	}
	nodeSelectorRequirements, err := parse.NodeSelector(selector)
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

func CreateBacalhauClient() (*client.APIClient, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	bacalhauConfigDirPath := filepath.Join(home, ".bacalhau")
	config.SetUserKey(filepath.Join(bacalhauConfigDirPath, "user_id.pem"))
	config.SetLibp2pKey(filepath.Join(bacalhauConfigDirPath, "libp2p_private_key"))
	defaultConfig := config.ForEnvironment()
	if os.Getenv("BACALHAU_IPFS_SWARM_ADDRESSES") != "" {
		defaultConfig.Node.IPFS.SwarmAddresses = []string{os.Getenv("BACALHAU_IPFS_SWARM_ADDRESSES")}
	}
	config.Set(defaultConfig)
	_, err = config.Init(defaultConfig, filepath.Join(home, ".bacalhau"), "config", "yaml")
	if err != nil {
		return nil, err
	}
	apiHost := GetBacalhauApiHost()
	apiPort := uint16(1234)
	client := client.NewAPIClient(apiHost, apiPort)
	return client, err
}

func SubmitBacalhauJob(job *model.Job) (submittedJob *model.Job, err error) {
	client, err := CreateBacalhauClient()
	if err != nil {
		return nil, err
	}
	submittedJob, err = client.Submit(context.Background(), job)
	return submittedJob, err
}

func GetBacalhauJobState(jobId string) (*model.JobWithInfo, error) {
	client, err := CreateBacalhauClient()
	if err != nil {
		return nil, err
	}
	updatedJob, _, err := client.Get(context.Background(), jobId)
	return updatedJob, err
}

func GetBacalhauJobResults(submittedJob *model.Job, showAnimation bool, maxTime int) (results []model.PublishedResult, err error) {
	client, err := CreateBacalhauClient()
	if err != nil {
		return nil, err
	}

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
		Timeout:   model.DefaultDownloadTimeout,
		OutputDir: dir,
	}
	downloadSettings.OutputDir = dir
	downloaderProvider := util.NewStandardDownloaders(cm, downloadSettings)
	err := downloader.DownloadResults(context.Background(), results, downloaderProvider, downloadSettings)
	return err
}
