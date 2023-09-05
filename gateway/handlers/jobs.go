package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"

	"gorm.io/gorm"
)

func InitJobHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var requestData struct {
			ToolPath         string              `json:"toolPath"`
			ScatteringMethod string              `json:"scatteringMethod"`
			InputVectors     map[string][]string `json:"inputVectors"`
		}

		if err := utils.ReadRequestBody(r, &requestData); err != nil {
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
			return
		}

		ioList, err := ipwl.InitializeIo(requestData.ToolPath, requestData.ScatteringMethod, requestData.InputVectors)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error initializing IO: %v", err), http.StatusInternalServerError)
			return
		}

		// Iterate over the ioList
		for _, io := range ioList {
			// Convert the IO object to JSON
			ioJson, err := json.Marshal(io)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error converting IO to JSON: %v", err), http.StatusInternalServerError)
				return
			}

			// Write the IO JSON to a temporary file
			tempFile, err := utils.CreateAndWriteTempFile(bytes.NewReader(ioJson), "io.json")
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating temporary file: %v", err), http.StatusInternalServerError)
				return
			}
			defer os.Remove(tempFile.Name())

			// Pin the temporary file to IPFS to get the CID
			cid, err := ipfs.PinFile(tempFile.Name())
			if err != nil {
				http.Error(w, fmt.Sprintf("Error pinning file to IPFS: %v", err), http.StatusInternalServerError)
				return
			}

			// Create a new Job object and store it in the database
			job := models.Job{
				InitialIoCID:  cid,
				InitialIoJson: string(ioJson),
				Status:        "initialized",
			}

			if result := db.Create(&job); result.Error != nil {
				http.Error(w, fmt.Sprintf("Error creating job: %v", result.Error), http.StatusInternalServerError)
				return
			}
		}

		utils.SendJSONResponseWithCID(w, "Jobs initialized successfully")
	}
}

func RunJobHandler() {
	// log that this function is being hit
	fmt.Print("RunJobHandler hit")
}

func GetJobHandler() {
	// log that this function is being hit
	fmt.Print("GetJobHandler hit")
}

// func GetJobsHandler() {
// 	// log that this function is being hit
// 	fmt.Print("GetJobsHandler hit")
// }
