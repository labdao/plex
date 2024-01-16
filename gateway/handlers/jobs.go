package handlers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
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

		params := mux.Vars(r)
		bacalhauJobID := params["bacalhauJobID"]

		var job models.Job
		if result := db.Preload("Outputs.Tags").Preload("Inputs.Tags").First(&job, "bacalhau_job_id = ?", bacalhauJobID); result.Error != nil {
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

func addTag(db *gorm.DB, name, tagType string) error {
	if name == "" || tagType == "" {
		return fmt.Errorf("Tag name and type are required")
	}

	tag := models.Tag{
		Name: name,
		Type: tagType,
	}

	result := db.Create(&tag)
	if result.Error != nil {
		if utils.IsDuplicateKeyError(result.Error) {
			log.Printf("A tag with the same name already exists: %s", name)
			return nil
		} else {
			return fmt.Errorf("Error creating tag: %v", result.Error)
		}
	}

	return nil
}

func UpdateJobHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			utils.SendJSONError(w, "Only PATCH method is supported", http.StatusBadRequest)
			return
		}

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

		log.Println("Getting job history...")
		events, err := bacalhau.GetBacalhauJobEvents(job.BacalhauJobID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting job history: %v", err), http.StatusInternalServerError)
			return
		}

		for _, event := range events {
			if event.ExecutionState != nil {
				switch event.ExecutionState.New {
				case model.ExecutionStateBidAccepted:
					log.Printf("Bid Accepted Time: %s", event.Time)
					job.StartedAt = event.Time
				case model.ExecutionStateCompleted:
					log.Printf("Job Completed Time: %s", event.Time)
					job.CompletedAt = event.Time
				}
			}
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
			fmt.Println("Job completed, getting output directory CID...")
			outputDirCID := updatedJob.State.Executions[0].PublishedResult.CID
			outputFileEntries, err := ipfs.ListFilesInDirectory(outputDirCID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error listing files in directory: %v", err), http.StatusInternalServerError)
				return
			}

			existingOutputs := make(map[string]struct{})
			for _, output := range job.Outputs {
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
							CID:           fileEntry["CID"],
							WalletAddress: job.WalletAddress,
							Filename:      fileName,
							Tags:          []models.Tag{generatedTag, experimentTag, extensionTag},
							Timestamp:     time.Now(),
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

				log.Println("Adding DataFile to Job.Outputs with CID:", dataFile.CID)
				job.Outputs = append(job.Outputs, dataFile)
				log.Println("Updated Job.Outputs:", job.Outputs)
			}

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

	outputChan := make(chan string)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			outputChan <- scanner.Text()
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			outputChan <- scanner.Text()
		}
	}()

	go func() {
		for {
			select {
			case outputLine, ok := <-outputChan:
				if !ok {
					return
				}
				if err := conn.WriteMessage(websocket.TextMessage, []byte(outputLine)); err != nil {
					log.Println("Error sending message through WebSocket:", err)
					return
				}
			}
		}
	}()

	cmd.Wait()
}
