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

	// Marshal the inputs to JSON
	jsonBytes, err := json.Marshal(inputs)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", tool.RayServiceURL, bytes.NewBuffer(jsonBytes))
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
	return CreateRayJob(toolPath, inputs)
}
