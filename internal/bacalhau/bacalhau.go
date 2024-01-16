package bacalhau

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bacalhau-project/bacalhau/cmd/util/parse"
	"github.com/bacalhau-project/bacalhau/pkg/config"
	"github.com/bacalhau-project/bacalhau/pkg/model"
	"github.com/bacalhau-project/bacalhau/pkg/publicapi/client"
	"github.com/labdao/plex/internal/ipfs"
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

func CreateBacalhauJob(inputs map[string]interface{}, container, selector string, maxTime, memory int, cpu float64, gpu, network bool, annotations []string) (job *model.Job, err error) {
	fmt.Println("CreatebacalhauJob")
	fmt.Println(inputs)
	job, err = model.NewJobWithSaneProductionDefaults()
	if err != nil {
		return nil, err
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

	// same data as IO.Inputs except it has local paths instead of CID paths for files
	localJobInputs := map[string]interface{}{}

	for key, input := range inputs {
		fmt.Println("going thru inputs with new code")
		fmt.Println(input)
		switch value := input.(type) {
		case string:
			if strings.HasPrefix(value, "Qm") {
				fmt.Println("found file input")
				// Split the string on the "/" character to separate the CID and filename
				parts := strings.Split(value, "/")
				if len(parts) != 2 {
					fmt.Println("here input file")
					fmt.Println(value)
					return nil, fmt.Errorf("not a valid cid path")
				}

				cid, filename := parts[0], parts[1]
				localDir := "/inputs/" + key
				localPath := localDir + "/" + filename

				var bacalhauPath string
				cidIsDir, err := ipfs.IsDirectory(cid)
				if err != nil {
					return nil, err
				}
				if cidIsDir {
					bacalhauPath = localDir
				} else {
					bacalhauPath = localPath
				}
				job.Spec.Inputs = append(job.Spec.Inputs,
					model.StorageSpec{
						StorageSource: model.StorageSourceIPFS,
						CID:           cid,
						Path:          bacalhauPath,
					})
				localJobInputs[key] = localPath
			} else {
				fmt.Println("input is a string but does not have 'Qm' prefix")
				localJobInputs[key] = value
			}
		case []interface{}: // Changed from []string to []interface{}
			fmt.Println("found slice, checking each for 'Qm' prefix")
			var stringArray []string
			allValid := true
			for _, elem := range value {
				str, ok := elem.(string)
				if !ok || !strings.HasPrefix(str, "Qm") {
					allValid = false
					fmt.Println("element is not a string or does not have 'Qm' prefix:", elem)
					break
				}
				stringArray = append(stringArray, str)
			}
			if allValid && len(stringArray) > 0 {
				fmt.Println("found file array")
				var localFilePaths []string
				for i, elem := range value {
					str, _ := elem.(string)
					// Split the string on the "/" character to separate the CID and filename
					parts := strings.Split(str, "/")
					cid, filename := parts[0], parts[1]

					// Construct the path with the key and index 'i'
					indexedDir := fmt.Sprintf("/inputs/%s/%d", key, i)
					indexedPath := indexedDir + "/" + filename
					cidIsDir, err := ipfs.IsDirectory(cid)
					var bacalhauPath string
					if err != nil {
						return nil, err
					}
					if cidIsDir {
						bacalhauPath = indexedDir
					} else {
						bacalhauPath = indexedPath
					}
					job.Spec.Inputs = append(job.Spec.Inputs,
						model.StorageSpec{
							StorageSource: model.StorageSourceIPFS,
							CID:           cid,
							Path:          bacalhauPath,
						})
					localFilePaths = append(localFilePaths, indexedPath)
				}
				localJobInputs[key] = localFilePaths
			} else {
				localJobInputs[key] = input
			}
		default:
			localJobInputs[key] = input
		}
	}
	jsonBytes, err := json.Marshal(localJobInputs)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	jsonString := string(jsonBytes)
	envVar := fmt.Sprintf("PLEX_JOB_INPUTS=%s", jsonString)

	job.Spec.EngineSpec = model.NewDockerEngineBuilder(container).WithEnvironmentVariables(envVar).Build()

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
