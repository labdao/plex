package handlers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/bacalhau-project/bacalhau/pkg/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/bacalhau"
	"github.com/labdao/plex/internal/ipfs"
	"gorm.io/gorm"
)

func GetJobHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		// Get the ID from the URL
		params := mux.Vars(r)
		bacalhauJobID := params["bacalhauJobID"]

		var job models.Job
		if result := db.Preload("Outputs").Preload("Inputs").First(&job, "bacalhau_job_id = ?", bacalhauJobID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Job not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Job: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		log.Println("Fetched Job from DB: ", job)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(job); err != nil {
			http.Error(w, "Error encoding Job to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func UpdateJobHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			utils.SendJSONError(w, "Only PATCH method is supported", http.StatusBadRequest)
			return
		}

		// Get the ID from the URL
		params := mux.Vars(r)
		bacalhauJobID := params["bacalhauJobID"]

		var job models.Job
		if result := db.Preload("Outputs").Preload("Inputs").First(&job, "bacalhau_job_id = ?", bacalhauJobID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Job not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Job: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		log.Println("Updating job")
		updatedJob, err := bacalhau.GetBacalhauJobState(job.BacalhauJobID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating job %v", err), http.StatusInternalServerError)
		}

		// Update job state based on the external job state
		if updatedJob.State.State == model.JobStateCancelled {
			job.State = "failed"
		} else if updatedJob.State.State == model.JobStateError {
			job.State = "error"
		} else if updatedJob.State.State == model.JobStateQueued {
			job.State = "queued"
		} else if updatedJob.State.State == model.JobStateInProgress {
			job.State = "processing"
		} else if updatedJob.State.State == model.JobStateCompleted {
			job.State = "completed"
		} else if len(updatedJob.State.Executions) > 0 && updatedJob.State.Executions[0].State == model.ExecutionStateFailed {
			job.State = "failed"
		}
		log.Println("Updated job")

		// Save the updated job back to the database
		if err := db.Save(job).Error; err != nil {
			http.Error(w, fmt.Sprintf("Error saving Job: %v", err), http.StatusInternalServerError)
			return
		}

		if job.State == "completed" {
			// we always only do one execution at the moment
			fmt.Println("Job completed, getting output directory CID")
			outputDirCID := updatedJob.State.Executions[0].PublishedResult.CID
			outputFileEntries, err := ipfs.ListFilesInDirectory(outputDirCID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error listing files in directory: %v", err), http.StatusInternalServerError)
				return
			}

			// Create a map of existing output CIDs for quick lookup
			existingOutputs := make(map[string]struct{})
			for _, output := range job.Outputs {
				existingOutputs[output.CID] = struct{}{}
			}

			for _, fileEntry := range outputFileEntries {
				// Check if fileEntry is already in Job.Outputs
				if _, exists := existingOutputs[fileEntry["CID"]]; exists {
					// Skip if the file is already in the Job.Outputs
					continue
				}

				log.Println("Attempting to find DataFile in DB with CID: ", fileEntry["CID"])

				var dataFile models.DataFile
				result := db.Where("cid = ?", fileEntry["CID"]).First(&dataFile)
				if result.Error != nil {
					if errors.Is(result.Error, gorm.ErrRecordNotFound) {
						log.Println("DataFile not found in DB, creating new record with CID: ", fileEntry["CID"])

						dataFile.CID = fileEntry["CID"]
						dataFile.Filename = fileEntry["filename"]

						log.Println("Adding tags to new DataFile with CID:", dataFile.CID)
						if err := AddTagsToDataFile(db, dataFile.CID, []string{"generated"}); err != nil {
							http.Error(w, fmt.Sprintf("Error adding tags to datafile: %v", err), http.StatusInternalServerError)
							return
						}

						log.Println("Saving new DataFile to DB with CID:", dataFile.CID)
						if err := db.Create(&dataFile).Error; err != nil {
							http.Error(w, fmt.Sprintf("Error creating DataFile record: %v", err), http.StatusInternalServerError)
							return
						}

						log.Println("DataFile successfully saved to DB:", dataFile)
					} else {
						http.Error(w, fmt.Sprintf("Error querying DataFile table: %v", result.Error), http.StatusInternalServerError)
						return
					}
				}

				// Then add the DataFile to the Job.Outputs
				log.Println("Adding DataFile to Job.Outputs with CID:", dataFile.CID)
				job.Outputs = append(job.Outputs, dataFile)
				log.Println("Updated Job.Outputs:", job.Outputs)
			}

			// Update job in the database with new Outputs (this may need adjustment depending on your ORM)
			if err := db.Save(&job).Error; err != nil {
				http.Error(w, fmt.Sprintf("Error updating job: %v", err), http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(job); err != nil {
			http.Error(w, "Error encoding Job to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func StreamJobLogsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Streaming job logs")
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Check the origin of the request and return true if it's allowed
			// Here's a simple example that allows any origin:
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	params := mux.Vars(r)
	bacalhauJobID := params["bacalhauJobID"]
	bacalhauApiHost := bacalhau.GetBacalhauApiHost()
	cmd := exec.Command("bacalhau", "logs", "-f", "--api-host", bacalhauApiHost, bacalhauJobID)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Error creating stdout pipe:", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println("Error creating stderr pipe:", err)
		return
	}

	fmt.Println("Starting command", cmd)
	if err := cmd.Start(); err != nil {
		log.Println("Error starting command:", err)
		return
	}

	// Channel to gather output
	outputChan := make(chan string)

	// Read from stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			outputChan <- scanner.Text()
		}
	}()

	// Read from stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			outputChan <- scanner.Text()
		}
	}()

	// Write to WebSocket
	go func() {
		for {
			select {
			case outputLine, ok := <-outputChan:
				if !ok {
					// Handle closed channel, if necessary
					return
				}
				// Send to WebSocket
				if err := conn.WriteMessage(websocket.TextMessage, []byte(outputLine)); err != nil {
					log.Println("Error sending message through WebSocket:", err)
					return
				}
			}
		}
	}()

	// If you need to wait for cmd to finish
	cmd.Wait()
}
