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

func CreateRayJob(modelPath string, rayJobID string, inputs map[string]interface{}, db *gorm.DB) (*http.Response, error) {
	log.Printf("Creating Ray job with modelPath: %s and inputs: %+v\n", modelPath, inputs)
	model, _, err := ipwl.ReadModelConfig(modelPath, db)
	if err != nil {
		return nil, err
	}

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

	// Marshal the inputs to JSON
	jsonBytes, err := json.Marshal(adjustedInputs)
	if err != nil {
		return nil, err
	}

	log.Printf("Submitting Ray job with payload: %s\n", string(jsonBytes))

	// construct from env var BUCKET ENDPOINT + model.RayServiceEndpoint
	rayServiceURL := GetRayApiHost() + model.RayServiceEndpoint
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

func SubmitRayJob(modelPath string, rayJobID string, inputs map[string]interface{}, db *gorm.DB) (*http.Response, error) {
	log.Printf("Creating Ray job with modelPath: %s and inputs: %+v\n", modelPath, inputs)
	resp, err := CreateRayJob(modelPath, rayJobID, inputs, db)
	if err != nil {
		log.Printf("Error creating Ray job: %v\n", err)
		return nil, err
	}

	log.Printf("Ray job finished with response status: %s\n", resp.Status)
	return resp, nil
}
