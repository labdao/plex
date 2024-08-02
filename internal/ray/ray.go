package ray

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/internal/ipwl"
	"gorm.io/gorm"
)

var rayClient *http.Client
var once sync.Once

func GetRayApiHost() string {
	// For colabfold local testing set this env var to http://colabfold-service:<PORT>
	rayApiHost, exists := os.LookupEnv("RAY_API_HOST")
	if exists {
		return rayApiHost
	} else {
		return "http://localhost:8000" // Default Ray API host
	}
}

func GetRayJobApiHost() string {
	// For colabfold local testing set this env var to http://colabfold-service:<PORT>
	rayApiHost, exists := os.LookupEnv("RAY_JOB_API_HOST")
	if exists {
		return rayApiHost
	} else {
		return "http://localhost:8265" // Default Ray API host
	}
}

// Prevents race conditions with Ray Client
func GetRayClient() *http.Client {
	once.Do(func() {
		rayClient = &http.Client{}
	})
	return rayClient
}

func handleSingleElementInput(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case float64, int, int64:
		// Convert numeric values to string
		return fmt.Sprintf("%v", v), nil
	case nil:
		return "", nil
	default:
		return "", fmt.Errorf("unsupported type: %T", v)
	}
}

func CreateRayJob(job *models.Job, modelPath string, rayJobID string, inputs map[string]interface{}, db *gorm.DB) (*http.Response, error) {
	model, _, err := ipwl.ReadModelConfig(modelPath, db)
	if err != nil {
		return nil, err
	}
	var jsonBytes []byte
	var rayServiceURL string

	// Validate input keys
	err = validateInputKeys(inputs, model.Inputs)
	if err != nil {
		return nil, err
	}

	adjustedInputs := make(map[string]interface{})
	for key, value := range inputs {
		switch v := value.(type) {
		case []interface{}:
			if len(v) == 1 {
				if v[0] == nil {
					adjustedInputs[key] = nil
				} else {
					adjustedValue, err := handleSingleElementInput(v[0])
					if err != nil {
						return nil, fmt.Errorf("invalid input for key %s: %v", key, err)
					}
					adjustedInputs[key] = adjustedValue
				}
			} else {
				return nil, fmt.Errorf("expected a single-element slice for key %s, got: %v", key, v)
			}
		case string, float64, int, nil:
			adjustedValue, err := handleSingleElementInput(value)
			if err != nil {
				return nil, fmt.Errorf("invalid input for key %s: %v", key, err)
			}
			adjustedInputs[key] = adjustedValue
		default:
			return nil, fmt.Errorf("unsupported type for key %s: %T", key, value)
		}
	}
	//add rayJobID to inputs
	fmt.Printf("adding rayJobID to the adjustedInputs: %s\n", rayJobID)
	adjustedInputs["uuid"] = rayJobID

	if job.JobType == models.JobTypeService {
		// Marshal the inputs to JSON
		jsonBytes, err := json.Marshal(adjustedInputs)
		if err != nil {
			return nil, err
		}

		log.Printf("Submitting Ray job with payload: %s\n", string(jsonBytes))

		rayServiceURL = GetRayApiHost() + model.RayEndpoint
		// Create the HTTP request

	} else if job.JobType == models.JobTypeJob {
		if pdbValue, ok := adjustedInputs["pdb"].(string); ok {
			adjustedInputs["pdb"] = "s3://labdao-testeks/test.pdb"
			log.Println("PDB value is: ", pdbValue)
		}

		inputsJSON, err := json.Marshal(adjustedInputs)
		if err != nil {
			return nil, err
		}
		log.Printf("Submitting Ray job with testeks inputs: %s\n", inputsJSON)

		rayServiceURL = GetRayJobApiHost() + model.RayEndpoint
		runtimeEnv := map[string]interface{}{
			"env_vars": map[string]string{
				"REQUEST_UUID":   rayJobID,
				"RAY_JOB_INPUTS": string(inputsJSON),
			},
		}

		// Create the request body for the Ray job
		reqBody := map[string]interface{}{
			"entrypoint":    model.RayJobEntrypoint,
			"submission_id": rayJobID,
			"runtime_env":   runtimeEnv,
		}
		jsonBytes, err = json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
		log.Printf("Submitting Ray job with payload: %s\n", string(jsonBytes))

	}
	req, err := http.NewRequest("POST", rayServiceURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request to the Ray service
	client := GetRayClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetRayJobStatus(rayJobID string) (string, error) {
	rayServiceURL := GetRayJobApiHost() + "/api/jobs/" + rayJobID
	req, err := http.NewRequest("GET", rayServiceURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request to the Ray service
	client := GetRayClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Printf("Ray job status response: %s\n", string(body))
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	status, ok := data["status"]
	if !ok {
		log.Fatal("Status field not found")
	}
	return status.(string), nil
}

func JobIsRunning(rayJobID string) bool {
	status, err := GetRayJobStatus(rayJobID)
	if err != nil {
		return false
	}
	return strings.ToLower(status) == string(models.JobStateRunning)
}

func JobIsPending(rayJobID string) bool {
	status, err := GetRayJobStatus(rayJobID)
	if err != nil {
		return false
	}
	return strings.ToLower(status) == string(models.JobStatePending)
}

func JobSucceeded(rayJobID string) bool {
	status, err := GetRayJobStatus(rayJobID)
	if err != nil {
		return false
	}
	return strings.ToLower(status) == string(models.JobStateSucceeded)
}

func JobFailed(rayJobID string) bool {
	status, err := GetRayJobStatus(rayJobID)
	if err != nil {
		return false
	}
	return strings.ToLower(status) == string(models.JobStateFailed)
}

func JobStopped(rayJobID string) bool {
	status, err := GetRayJobStatus(rayJobID)
	if err != nil {
		return false
	}
	return strings.ToLower(status) == string(models.JobStateStopped)
}

func validateInputKeys(inputVectors map[string]interface{}, modelInputs map[string]ipwl.ModelInput) error {
	for inputKey := range inputVectors {
		if _, exists := modelInputs[inputKey]; !exists {
			log.Printf("The argument %s is not in the model inputs.\n", inputKey)
			log.Printf("Available keys: %v\n", modelInputs)
			return fmt.Errorf("the argument %s is not in the model inputs", inputKey)
		}
	}
	return nil
}

func SubmitRayJob(job models.Job, modelPath string, rayJobID string, inputs map[string]interface{}, db *gorm.DB) (*http.Response, error) {
	log.Printf("Creating Ray job with modelPath: %s and inputs: %+v\n", modelPath, inputs)
	resp, err := CreateRayJob(&job, modelPath, rayJobID, inputs, db)
	if err != nil {
		log.Printf("Error creating Ray job: %v\n", err)
		return nil, err
	}

	if job.JobType == models.JobTypeService {
		log.Printf("Ray job finished with response status: %s\n", resp.Status)
		return resp, nil
	} else if job.JobType == models.JobTypeJob {
		log.Printf("Ray job submitted with response status: %s\n", resp.Status)
		return resp, nil
	}
	return nil, fmt.Errorf("unsupported job type: %s", job.JobType)
}
