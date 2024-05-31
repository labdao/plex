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

func CreateRayJob(toolPath string, inputs map[string]interface{}) (*http.Response, error) {
	tool, _, err := ipwl.ReadToolConfig(toolPath)
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
		//TODO_PR#970: work on input file handling after convexity side key, location has been moved to uri
		// until then, job submission wont work with input files
		// if key == "pdb" {
		// 	// // Package both the path and the location into a single JSON string or a structured object
		// 	// pdbInfo := map[string]string{
		// 	// 	"key":      value.(string),
		// 	// 	"location": "convexity",
		// 	// }
		// 	// pdbJSON, err := json.Marshal(pdbInfo)
		// 	// if err != nil {
		// 	// 	return nil, fmt.Errorf("error marshaling pdb info: %v", err)
		// 	// }
		// 	// adjustedInputs[key] = string(pdbJSON)
		// 	var jsonBytes []byte
		// 	var err error

		// 	// Check if the value is directly a map[string]interface{}
		// 	if pdbMap, ok := value.(map[string]interface{}); ok {
		// 		jsonBytes, err = json.Marshal(pdbMap)
		// 		if err != nil {
		// 			return nil, fmt.Errorf("error marshaling pdb map to JSON: %v", err)
		// 		}
		// 	} else if valSlice, ok := value.([]interface{}); ok && len(valSlice) == 1 {
		// 		// If not a map, check if it's a slice with a single map element
		// 		if pdbMap, ok := valSlice[0].(map[string]interface{}); ok {
		// 			jsonBytes, err = json.Marshal(pdbMap)
		// 			if err != nil {
		// 				return nil, fmt.Errorf("error marshaling pdb map within a slice to JSON: %v", err)
		// 			}
		// 		} else if pdbMap, ok := valSlice[0].(map[string]string); ok {
		// 			// Handle map[string]string if present in the slice
		// 			tempMap := make(map[string]interface{})
		// 			for k, v := range pdbMap {
		// 				tempMap[k] = v
		// 			}
		// 			jsonBytes, err = json.Marshal(tempMap)
		// 			if err != nil {
		// 				return nil, fmt.Errorf("error marshaling pdb map (string map) within a slice to JSON: %v", err)
		// 			}
		// 		} else {
		// 			return nil, fmt.Errorf("expected a map for key 'pdb' within a slice, got: %v", valSlice[0])
		// 		}
		// 	} else {
		// 		return nil, fmt.Errorf("expected a map or a single-element slice with a map for key 'pdb', got: %v", value)
		// 	}

		// 	adjustedInputs[key] = string(jsonBytes)
		// } else {
		if valSlice, ok := value.([]interface{}); ok && len(valSlice) == 1 {
			if valStr, ok := valSlice[0].(string); ok {
				adjustedInputs[key] = valStr
			} else {
				return nil, fmt.Errorf("expected a string for key %s, got: %v", key, value)
			}
		} else {
			return nil, fmt.Errorf("expected a single-element slice for key %s, got: %v", key, value)
		}
		// }
	}

	// Marshal the inputs to JSON
	jsonBytes, err := json.Marshal(adjustedInputs)
	if err != nil {
		return nil, err
	}

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

func SubmitRayJob(toolPath string, inputs map[string]interface{}) (*http.Response, error) {
	log.Printf("Creating Ray job with toolPath: %s and inputs: %+v\n", toolPath, inputs)
	resp, err := CreateRayJob(toolPath, inputs)
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
