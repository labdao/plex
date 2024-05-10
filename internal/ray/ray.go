package ray

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var rayClient *http.Client
var once sync.Once

func GetRayApiHost() string {
	rayApiHost, exists := os.LookupEnv("RAY_API_HOST")
	if exists {
		return rayApiHost
	} else {
		return "localhost:8000" // Default Ray API host
	}
}

// Prevents race conditions with Ray Client
func GetRayClient() *http.Client {
	once.Do(func() {
		rayClient = &http.Client{}
	})
	return rayClient
}

func CreateRayJob(inputs map[string]interface{}, endpoint string) (*http.Response, error) {
	// Process the inputs to match the expected structure for the Ray service
	processedInputs := make(map[string]interface{})
	for key, value := range inputs {
		switch key {
		case "pdb":
			// Assuming 'pdb' requires special handling
			pdbData, ok := value.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid type for pdb, expected map[string]interface{}")
			}
			// Validate and process pdbData if necessary
			// For example, check if 'key' and 'location' exist
			if _, ok := pdbData["key"]; !ok {
				return nil, fmt.Errorf("missing 'key' in pdb data")
			}
			if _, ok := pdbData["location"]; !ok {
				return nil, fmt.Errorf("missing 'location' in pdb data")
			}
			// Add processed pdb data to processedInputs
			processedInputs["pdb"] = pdbData
		default:
			// For all other keys, use the value as is
			processedInputs[key] = value
		}
	}

	// Marshal the processed inputs to JSON
	jsonBytes, err := json.Marshal(processedInputs)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Print the JSON request body
	fmt.Printf("HTTP Request Body: %s\n", string(jsonBytes))

	apiHost := GetRayApiHost()
	url := fmt.Sprintf("http://%s%s", apiHost, endpoint)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Comment out the actual sending of the request
	// client := GetRayClient()
	// resp, err := client.Do(req)
	// if err != nil {
	//     return nil, err
	// }

	// Return nil response and no error for now
	return nil, nil
}

// TODO: when we are ready to submit requests

// func CreateRayJob(inputs map[string]interface{}, endpoint string) (*http.Response, error) {
// 	jsonBytes, err := json.Marshal(inputs)
// 	if err != nil {
// 		log.Fatalf("Error marshaling JSON: %v", err)
// 	}

// 	apiHost := GetRayApiHost()
// 	url := fmt.Sprintf("http://%s%s", apiHost, endpoint)

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	client := GetRayClient()
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

func SubmitRayJob(inputs map[string]interface{}, endpoint string) (*http.Response, error) {
	return CreateRayJob(inputs, endpoint)
}
