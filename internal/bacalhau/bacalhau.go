package bacalhau

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bacalhau-project/bacalhau/cmd/util/parse"
	"github.com/bacalhau-project/bacalhau/pkg/config"
	"github.com/bacalhau-project/bacalhau/pkg/model"
	"github.com/bacalhau-project/bacalhau/pkg/publicapi/client"
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

func CreateBacalhauJob(fileInputs map[string]string, fileArrayInputs map[string][]string, container, selector string, baseCmd, params []string, maxTime, memory int, cpu float64, gpu, network bool, annotations []string) (job *model.Job, err error) {
	log.Println("Creating job inside v2 function")
	job, err = model.NewJobWithSaneProductionDefaults()
	if err != nil {
		return nil, err
	}
	fmt.Println("container baseCmd", baseCmd)
	fmt.Println("container params", params)
	if len(baseCmd) > 0 && len(params) > 0 {
		job.Spec.EngineSpec = model.NewDockerEngineBuilder(container).WithEntrypoint(baseCmd...).WithParameters(params...).Build()
	} else if len(baseCmd) > 0 {
		job.Spec.EngineSpec = model.NewDockerEngineBuilder(container).WithEntrypoint(baseCmd...).Build()
	} else if len(params) > 0 {
		job.Spec.EngineSpec = model.NewDockerEngineBuilder(container).WithParameters(params...).Build()
	} else {
		job.Spec.EngineSpec = model.NewDockerEngineBuilder(container).Build()
	}

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
	if cpu > 0 {
		job.Spec.Resources.CPU = fmt.Sprintf("%f", cpu)
	}
	if gpu {
		job.Spec.Resources.GPU = "1"
	}
	if network {
		job.Spec.Network = model.NetworkConfig{Type: model.NetworkFull}
	}
	job.Spec.Inputs = []model.StorageSpec{}
	for key, input := range fileInputs {
		// Split the string on the "/" character to separate the CID and filename
		parts := strings.Split(input, "/")
		if len(parts) != 2 {
			fmt.Println("here input file")
			fmt.Println(input)
			return nil, fmt.Errorf("not a valid cid path")
		}

		cid, _ := parts[0], parts[1]

		job.Spec.Inputs = append(job.Spec.Inputs,
			model.StorageSpec{
				StorageSource: model.StorageSourceIPFS,
				CID:           cid,
				Path:          "/inputs/" + key,
			})
	}

	for key, inputs := range fileArrayInputs {
		for i, input := range inputs {
			// Split the string on the "/" character to separate the CID and filename
			parts := strings.Split(input, "/")
			if len(parts) != 2 {
				fmt.Println("here input file array")
				fmt.Println(i)
				fmt.Println(input)
				return nil, fmt.Errorf("not a valid cid path")
			}

			cid, _ := parts[0], parts[1]

			// Construct the path with the key and index 'i'
			indexedPath := fmt.Sprintf("/inputs/%s/%d", key, i)

			job.Spec.Inputs = append(job.Spec.Inputs,
				model.StorageSpec{
					StorageSource: model.StorageSourceIPFS,
					CID:           cid,
					Path:          indexedPath,
				})
		}
	}

	job.Spec.Outputs = []model.StorageSpec{{Name: "outputs", StorageSource: model.StorageSourceIPFS, Path: "/outputs"}}
	log.Println("returning job")
	return job, err
}

func CreateBacalhauClient() (*client.APIClient, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	bacalhauConfigDirPath := filepath.Join(home, ".bacalhau")
	bacalhauConfig, err := config.Load(bacalhauConfigDirPath, "config", "yaml")
	if err != nil {
		return nil, err
	}
	if os.Getenv("BACALHAU_IPFS_SWARM_ADDRESSES") != "" {
		swarmAddresses := []string{os.Getenv("BACALHAU_IPFS_SWARM_ADDRESSES")}
		bacalhauConfig.Node.IPFS.SwarmAddresses = swarmAddresses
	}
	config.SetUserKey(filepath.Join(bacalhauConfigDirPath, "user_id.pem"))
	config.SetLibp2pKey(filepath.Join(bacalhauConfigDirPath, "libp2p_private_key"))
	config.Set(bacalhauConfig)

	_, err = config.Init(bacalhauConfig, bacalhauConfigDirPath, "config", "yaml")
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
