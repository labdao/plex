package handlers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

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
		jobID, err := strconv.Atoi(params["jobID"])
		if err != nil {
			http.Error(w, fmt.Sprintf("Job ID (%v) could not be converted to int", params["jobID"]), http.StatusNotFound)
		}

		var job models.Job
		if result := db.Preload("OutputFiles.Tags").Preload("InputFiles.Tags").First(&job, "id = ?", jobID); result.Error != nil {
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
		jobID, err := strconv.Atoi(params["jobID"])
		if err != nil {
			http.Error(w, fmt.Sprintf("Job ID (%v) could not be converted to int", params["jobID"]), http.StatusNotFound)
		}

		var job models.Job
		if result := db.Preload("OutputFiles").Preload("InputFiles").First(&job, "id = ?", jobID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Job not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Job: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		var flow models.Flow
		if result := db.First(&flow, "cid = ?", job.FlowID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				log.Println("Flow not found for given FlowID:", job.FlowID)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Flow: %v", result.Error), http.StatusInternalServerError)
				return
			}
		}
		log.Println("Flow name:", flow.Name)

		log.Println("Updating job")
		updatedJob, err := bacalhau.GetBacalhauJobState(job.BacalhauJobID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating job %v", err), http.StatusInternalServerError)
		}

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
			for _, output := range job.OutputFiles {
				existingOutputs[output.CID] = struct{}{}
			}

			log.Printf("Number of output files: %d\n", len(outputFileEntries))
			fileNames := make([]string, 0, len(outputFileEntries))
			for _, fileEntry := range outputFileEntries {
				fileNames = append(fileNames, fileEntry["filename"])
			}
			log.Printf("Output file names: %v\n", fileNames)

			for _, fileEntry := range outputFileEntries {
				log.Printf("Processing fileEntry with CID: %s, Filename: %s", fileEntry["CID"], fileEntry["filename"])

				if _, exists := existingOutputs[fileEntry["CID"]]; exists {
					continue
				}

				log.Println("Attempting to find DataFile in DB with CID: ", fileEntry["CID"])

				var dataFile models.DataFile
				result := db.Where("cid = ?", fileEntry["CID"]).First(&dataFile)
				if result.Error != nil {
					if errors.Is(result.Error, gorm.ErrRecordNotFound) {
						log.Println("DataFile not found in DB, creating new record with CID: ", fileEntry["CID"])

						var generatedTag models.Tag
						if err := db.Where("name = ?", "generated").First(&generatedTag).Error; err != nil {
							http.Error(w, fmt.Sprintf("Error finding generated tag: %v", err), http.StatusInternalServerError)
							return
						}

						experimentTagName := "experiment:" + flow.Name
						var experimentTag models.Tag
						result := db.Where("name = ?", experimentTagName).First(&experimentTag)
						if result.Error != nil {
							if errors.Is(result.Error, gorm.ErrRecordNotFound) {
								experimentTag = models.Tag{Name: experimentTagName, Type: "autogenerated"}
								if err := db.Create(&experimentTag).Error; err != nil {
									http.Error(w, fmt.Sprintf("Error creating experiment tag: %v", err), http.StatusInternalServerError)
									return
								}
							} else {
								http.Error(w, fmt.Sprintf("Error querying Tag table: %v", result.Error), http.StatusInternalServerError)
								return
							}
						}

						fileName := fileEntry["filename"]
						dotIndex := strings.LastIndex(fileName, ".")
						var extension string
						if dotIndex != -1 && dotIndex < len(fileName)-1 {
							extension = fileName[dotIndex+1:]
						} else {
							extension = "utility"
						}
						extensionTagName := "file_extension:" + extension

						var extensionTag models.Tag
						result = db.Where("name = ?", extensionTagName).First(&extensionTag)
						if result.Error != nil {
							if errors.Is(result.Error, gorm.ErrRecordNotFound) {
								extensionTag = models.Tag{Name: extensionTagName, Type: "autogenerated"}
								if err := db.Create(&extensionTag).Error; err != nil {
									http.Error(w, fmt.Sprintf("Error creating extension tag: %v", err), http.StatusInternalServerError)
									return
								}
							} else {
								http.Error(w, fmt.Sprintf("Error querying Tag table: %v", result.Error), http.StatusInternalServerError)
								return
							}
						}

						log.Println("Saving generated DataFile to DB with CID:", fileEntry["CID"])

						dataFile = models.DataFile{
							CID:       fileEntry["CID"],
							Filename:  fileName,
							Tags:      []models.Tag{generatedTag, experimentTag, extensionTag},
							Timestamp: time.Now(),
						}

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
				job.OutputFiles = append(job.OutputFiles, dataFile)
				log.Println("Updated Job.Outputs:", job.OutputFiles)
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
