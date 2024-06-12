package ray

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

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
		return "localhost:8000" // Default Ray API host
	}
}

// Prevents race conditions with Ray Client
// TODO_PR#970 - revisit timeout later
//
//	func GetRayClient(maxRunningTime int) *http.Client {
//		once.Do(func() {
//			rayClient = &http.Client{
//				Timeout: time.Second * time.Duration(maxRunningTime),
//			}
//		})
//		return rayClient
//	}
func GetRayClient() *http.Client {
	once.Do(func() {
		rayClient = &http.Client{}
	})
	return rayClient
}

func CreateRayJob(toolPath string, rayJobID string, inputs map[string]interface{}, db *gorm.DB) (*http.Response, error) {
	log.Printf("Creating Ray job with toolPath: %s and inputs: %+v\n", toolPath, inputs)
	tool, _, err := ipwl.ReadToolConfig(toolPath, db)
	if err != nil {
		return nil, err
	}

	// Validate input keys
	err = validateInputKeys(inputs, tool.Inputs)
	if err != nil {
		return nil, err
	}

	adjustedInputs := make(map[string]string)
	for key, value := range inputs {
		if valSlice, ok := value.([]interface{}); ok && len(valSlice) == 1 {
			if valStr, ok := valSlice[0].(string); ok {
				adjustedInputs[key] = valStr
			} else {
				return nil, fmt.Errorf("expected a string for key %s, got: %v", key, value)
			}
		} else {
			return nil, fmt.Errorf("expected a single-element slice for key %s, got: %v", key, value)
		}
	}

	//add rayJobID to inputs
	fmt.Printf("adding rayJobID to the adjustedInputs: %s\n", rayJobID)
	adjustedInputs["uuid"] = rayJobID

	// Marshal the inputs to JSON
	jsonBytes, err := json.Marshal(adjustedInputs)
	if err != nil {
		return nil, err
	}

	log.Printf("Submitting Ray job with payload: %s\n", string(jsonBytes))

	// construct from env var BUCKET ENDPOINT + tool.RayServiceEndpoint
	rayServiceURL := GetRayApiHost() + tool.RayServiceEndpoint
	// Create the HTTP request
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

func validateInputKeys(inputVectors map[string]interface{}, toolInputs map[string]ipwl.ToolInput) error {
	for inputKey := range inputVectors {
		if _, exists := toolInputs[inputKey]; !exists {
			log.Printf("The argument %s is not in the tool inputs.\n", inputKey)
			log.Printf("Available keys: %v\n", toolInputs)
			return fmt.Errorf("the argument %s is not in the tool inputs", inputKey)
		}
	}
	return nil
}

func SubmitRayJob(toolPath string, rayJobID string, inputs map[string]interface{}, db *gorm.DB) (*http.Response, error) {
	log.Printf("Creating Ray job with toolPath: %s and inputs: %+v\n", toolPath, inputs)
	resp, err := CreateRayJob(toolPath, rayJobID, inputs, db)
	if err != nil {
		log.Printf("Error creating Ray job: %v\n", err)
		return nil, err
	}

	log.Printf("Ray job finished with response status: %s\n", resp.Status)
	// if resp.StatusCode != http.StatusOK {
	// 	return "", fmt.Errorf("failed to create Ray job, status code: %d", resp.StatusCode)
	// }

	// var responseBody map[string]string
	// if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
	// 	return "", fmt.Errorf("failed to decode response body: %v", err)
	// }

	// jobID, exists := responseBody["job_id"]
	// if !exists {
	// 	return "", fmt.Errorf("job_id not found in response")
	// }

	return resp, nil
}
