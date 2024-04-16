package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/internal/bacalhau"
	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"
	"gorm.io/gorm"
)

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid format for %s. Using default value. \n", name)
		return defaultValue
	}
	return value
}

// Consider allowing Job creation payloads to configure lower times
var maxQueueTime = time.Duration(getEnvAsInt("MAX_QUEUE_TIME_SECONDS", 259200)) * time.Second
var maxComputeTime = getEnvAsInt("MAX_COMPUTE_TIME_SECONDS", 259200)
var retryJobSleepTime = time.Duration(10) * time.Second

func StartJobQueues(db *gorm.DB) error {
	errChan := make(chan error, 2) // Buffer for two types of jobs

	go ProcessJobQueue(models.QueueTypeCPU, errChan, db)
	go ProcessJobQueue(models.QueueTypeGPU, errChan, db)

	// Wait for error from any queue
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

func ProcessJobQueue(queueType models.QueueType, errChan chan<- error, db *gorm.DB) {
	for {
		var job models.Job
		err := fetchOldestQueuedJob(&job, queueType, db)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("There are no Jobs Queued for %v, will recheck later\n", queueType)
			time.Sleep(5 * time.Second)
			continue
		} else if err != nil {
			errChan <- err
			return
		}

		if err := processJob(job.ID, db); err != nil {
			errChan <- err
			return
		}
	}
}

func MonitorRunningJobs(db *gorm.DB) error {
	for {
		var jobs []models.Job
		if err := fetchRunningJobsWithToolData(&jobs, db); err != nil {
			return err
		}
		fmt.Printf("There are %d running jobs\n", len(jobs))
		for _, job := range jobs {
			var ToolJson ipwl.Tool
			var checkpointCompatible string
			var maxRunningTime time.Duration
			if err := json.Unmarshal(job.Tool.ToolJson, &ToolJson); err != nil {
				return err
			}
			if ToolJson.CheckpointCompatible {
				checkpointCompatible = "True"
			} else {
				checkpointCompatible = "False"
			}
			if ToolJson.MaxRunningTime != 0 {
				maxRunningTime = time.Duration(ToolJson.MaxRunningTime) * time.Second
			} else {
				maxRunningTime = 2700 * time.Second
			}
			elapsed := time.Since(job.StartedAt)
			log.Printf("Job %d running for %v\n", job.ID, elapsed)
			if elapsed > maxRunningTime {
				fmt.Printf("Job %d has exceeded the maximum running time of %v, retrying job\n", job.ID, maxRunningTime)
				if job.RetryCount < 1 {
					job.RetryCount++
					job.State = models.JobStateQueued
					if err := updateJobRetryState(&job, db); err != nil {
						return err
					}
					fmt.Printf("retrying job %d\n", job.ID)
					if err := submitBacalhauJobAndUpdateID(&job, db, checkpointCompatible); err != nil {
						fmt.Printf("Error retrying job %d: %v\n", job.ID, err)
						return err
					}
					time.Sleep(retryJobSleepTime)
					continue
				} else {
					fmt.Printf("Job %d has already been retried, failing job\n", job.ID)
					job.State = models.JobStateFailed
					if err := db.Save(&job).Error; err != nil {
						fmt.Printf("Error updating job to failed: %v\n", err)
						return err
					}
				}
			}
			if err := checkRunningJob(job.ID, db); err != nil {
				return err
			}
		}
		fmt.Printf("Finished watching all running jobs, rehydrating watcher with jobs\n")
		time.Sleep(1 * time.Second) // Wait for some time before the next cycle
	}
}

func updateJobRetryState(job *models.Job, db *gorm.DB) error {
	return db.Model(job).Updates(models.Job{RetryCount: job.RetryCount, State: job.State}).Error
}

func fetchRunningJobsWithToolData(jobs *[]models.Job, db *gorm.DB) error {
	return db.Preload("Tool").Where("state = ?", models.JobStateRunning).Find(jobs).Error
}

func fetchOldestQueuedJob(job *models.Job, queueType models.QueueType, db *gorm.DB) error {
	return db.Where("state = ? AND queue = ?", models.JobStateQueued, queueType).Order("created_at ASC").First(job).Error
}

func fetchJobWithToolData(job *models.Job, id uint, db *gorm.DB) error {
	return db.Preload("Tool").First(&job, id).Error
}

func fetchJobWithToolAndFlowData(job *models.Job, id uint, db *gorm.DB) error {
	return db.Preload("Tool").Preload("Flow").First(&job, id).Error
}

func checkDBForJobStateCompleted(jobID uint, db *gorm.DB) (bool, error) {
	var job models.Job
	if result := db.First(&job, "id = ?", jobID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, result.Error
		}
	}

	if job.State == models.JobStateCompleted {
		return true, nil
	} else {
		return false, nil
	}
}

func processJob(jobID uint, db *gorm.DB) error {
	fmt.Printf("Processing job %v\n", jobID)
	var job models.Job
	err := fetchJobWithToolAndFlowData(&job, jobID, db)
	if err != nil {
		return err
	}

	var ToolJson ipwl.Tool
	if err := json.Unmarshal(job.Tool.ToolJson, &ToolJson); err != nil {
		return err
	}

	var checkpointCompatible string
	if ToolJson.CheckpointCompatible == true {
		checkpointCompatible = "True"
	} else {
		checkpointCompatible = "False"
	}

	// TODO may be: If checkpoint is false, do not pass job and flow UUIDs

	if job.BacalhauJobID == "" {
		if err := submitBacalhauJobAndUpdateID(&job, db, checkpointCompatible); err != nil {
			return err
		}
	}

	for {
		// we refresh the Job data from the database on every loop to make sure the state is accurate
		var job models.Job
		err := fetchJobWithToolAndFlowData(&job, jobID, db)
		if err != nil {
			return err
		}

		// Safety check against DB to make sure the job hasn't already been completed
		isCompleted, err := checkDBForJobStateCompleted(jobID, db)
		if err != nil {
			return fmt.Errorf("Error checking DB for Job state: %v", err)
		}
		if isCompleted {
			fmt.Printf("Job %v , %v already completed\n", job.ID, job.BacalhauJobID)
			return nil
		}

		fmt.Printf("Checking status for %v , %v\n", job.ID, job.BacalhauJobID)
		if time.Since(job.CreatedAt) > maxQueueTime {
			fmt.Printf("Job %v , %v timed out\n", job.ID, job.BacalhauJobID)
			return setJobStatus(&job, models.JobStateFailed, "timed out in queue", db)
		}

		bacalhauJob, err := bacalhau.GetBacalhauJobState(job.BacalhauJobID)
		if err != nil {
			return err
		}

		// keep retrying job if there is a capacity error until job times out
		// ideally replace with a query if the Bacalhau nodes have capacity
		fmt.Printf("Job %v , %v has state %v\n", job.ID, job.BacalhauJobID, bacalhauJob.State.State)
		if bacalhau.JobFailedWithCapacityError(bacalhauJob) {
			fmt.Printf("Job %v , %v failed with capacity full, will try again\n", job.ID, job.BacalhauJobID)
			if err := submitBacalhauJobAndUpdateID(&job, db, checkpointCompatible); err != nil {
				return err
			}
			time.Sleep(retryJobSleepTime) // Wait for a short period before checking the status again
			continue
		} else if bacalhau.JobIsRunning(bacalhauJob) {
			fmt.Printf("Job %v , %v is running\n", job.ID, job.BacalhauJobID)
			return setJobStatus(&job, models.JobStateRunning, "", db)
		} else if bacalhau.JobFailed(bacalhauJob) {
			fmt.Printf("Job %v , %v failed\n", job.ID, job.BacalhauJobID)
			return setJobStatus(&job, models.JobStateFailed, fmt.Sprintf("Error running with Bacalhau ID %v", job.BacalhauJobID), db)
		} else if bacalhau.JobCompleted(bacalhauJob) {
			fmt.Printf("Job %v , %v completed\n", job.ID, job.BacalhauJobID)
			if len(bacalhauJob.State.Executions) > 0 {
				return completeJobAndAddOutputFiles(&job, models.JobStateCompleted, bacalhauJob.State.Executions[0].PublishedResult.CID, db)
			}
			return setJobStatus(&job, models.JobStateFailed, fmt.Sprintf("Output execution data lost for %v", job.BacalhauJobID), db)
		} else if bacalhau.JobCancelled(bacalhauJob) {
			fmt.Printf("Job %v , %v cancelled\n", job.ID, job.BacalhauJobID)
			return setJobStatus(&job, models.JobStateFailed, fmt.Sprintf("Job %v cancelled", job.BacalhauJobID), db)
		} else {
			fmt.Printf("Job %v , %v has state %v, will requery\n", job.ID, job.BacalhauJobID, bacalhauJob.State.State)
			time.Sleep(retryJobSleepTime) // Wait for a short period before checking the status again
			continue
		}
	}
}

func checkRunningJob(jobID uint, db *gorm.DB) error {
	var job models.Job
	err := fetchJobWithToolAndFlowData(&job, jobID, db)
	if err != nil {
		return err
	}
	bacalhauJob, err := bacalhau.GetBacalhauJobState(job.BacalhauJobID)
	if err != nil && strings.Contains(err.Error(), "Job not found") {
		fmt.Printf("Job %v , %v has missing Bacalhau Job, failing Job\n", job.ID, job.BacalhauJobID)
		return setJobStatus(&job, models.JobStateFailed, fmt.Sprintf("Bacalhau job %v not found", job.BacalhauJobID), db)
	} else if err != nil {
		return err
	}

	if bacalhau.JobIsRunning(bacalhauJob) {
		fmt.Printf("Job %v , %v is still running nothing to do\n", job.ID, job.BacalhauJobID)
		return nil
	} else if bacalhau.JobBidAccepted(bacalhauJob) {
		fmt.Printf("Job %v , %v is still in bid accepted nothing to do\n", job.ID, job.BacalhauJobID)
		return nil
	} else if bacalhau.JobFailed(bacalhauJob) {
		fmt.Printf("Job %v , %v failed, updating status and adding output files\n", job.ID, job.BacalhauJobID)
		return setJobStatus(&job, models.JobStateFailed, fmt.Sprintf("Bacalhau job %v failed", job.BacalhauJobID), db)
	} else if bacalhau.JobCompleted(bacalhauJob) {
		fmt.Printf("Job %v , %v completed, updating status and adding output files\n", job.ID, job.BacalhauJobID)
		if len(bacalhauJob.State.Executions) > 0 {
			return completeJobAndAddOutputFiles(&job, models.JobStateCompleted, bacalhauJob.State.Executions[0].PublishedResult.CID, db)
		}
		return setJobStatus(&job, models.JobStateFailed, fmt.Sprintf("Output execution data lost for %v", job.BacalhauJobID), db)
	} else {
		fmt.Printf("Job %v , %v had unexpected Bacalhau state %v, marking as failed\n", job.ID, job.BacalhauJobID, bacalhauJob.State.State)
		return setJobStatus(&job, models.JobStateFailed, fmt.Sprintf("unexpected Bacalhau state %v", bacalhauJob.State.State), db)
	}
}

func setJobStatus(job *models.Job, state models.JobState, errorMessage string, db *gorm.DB) error {
	job.State = state
	job.StartedAt = time.Now().UTC()
	job.Error = errorMessage
	return db.Save(job).Error
}

func completeJobAndAddOutputFiles(job *models.Job, state models.JobState, outputDirCID string, db *gorm.DB) error {
	job.State = state
	outputFileEntries, err := ipfs.ListFilesInDirectory(outputDirCID)
	if err != nil {
		fmt.Printf("Error listing files in directory: %v", err)
		return err
	}

	var flow models.Flow
	if result := db.First(&flow, "id = ?", job.FlowID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("Flow not found for given FlowID:", job.FlowID)
		} else {
			fmt.Printf("Error fetching Flow: %v", result.Error)
			return result.Error
		}
	}
	log.Println("Flow name:", flow.Name)

	// Create a map of existing output CIDs for quick lookup
	existingOutputs := make(map[string]struct{})
	for _, output := range job.OutputFiles {
		existingOutputs[output.CID] = struct{}{}
	}
	fmt.Printf("Number of output files: %d\n", len(outputFileEntries))
	for _, fileEntry := range outputFileEntries {
		fmt.Printf("Processing fileEntry with CID: %s, Filename: %s\n", fileEntry["CID"], fileEntry["filename"])
		if _, exists := existingOutputs[fileEntry["CID"]]; exists {
			continue
		}
		var dataFile models.DataFile
		result := db.Where("cid = ?", fileEntry["CID"]).First(&dataFile)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				fmt.Println("DataFile not found in DB, creating new record with CID: ", fileEntry["CID"])

				var generatedTag models.Tag
				if err := db.Where("name = ?", "generated").First(&generatedTag).Error; err != nil {
					fmt.Printf("Error finding generated tag: %v\n", err)
					return err
				}
				experimentTagName := "experiment:" + flow.Name
				var experimentTag models.Tag
				result := db.Where("name = ?", experimentTagName).First(&experimentTag)
				if result.Error != nil {
					if errors.Is(result.Error, gorm.ErrRecordNotFound) {
						experimentTag = models.Tag{Name: experimentTagName, Type: "autogenerated"}
						if err := db.Create(&experimentTag).Error; err != nil {
							fmt.Printf("Error creating experiment tag: %v\n", err)
							return err
						}
					} else {
						fmt.Printf("Error querying Tag table: %v\n", result.Error)
						return result.Error
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
							fmt.Printf("Error creating extension tag: %v\n", err)
							return err
						}
					} else {
						fmt.Printf("Error querying Tag table: %v\n", result.Error)
						return err
					}
				}

				fmt.Println("Saving generated DataFile to DB with CID:", fileEntry["CID"])

				dataFile = models.DataFile{
					CID:           fileEntry["CID"],
					WalletAddress: job.WalletAddress,
					Filename:      fileName,
					Tags:          []models.Tag{generatedTag, experimentTag, extensionTag},
					Timestamp:     time.Now(),
				}

				if err := db.Create(&dataFile).Error; err != nil {
					fmt.Printf("Error creating DataFile record: %v\n", err)
					return err
				}
			} else {
				fmt.Printf("Error querying DataFile table: %v\n", result.Error)
				return err
			}
		}
		var user models.User
		if err := db.Where("wallet_address = ?", job.WalletAddress).First(&user).Error; err != nil {
			fmt.Printf("Error finding user with wallet address %v\n: ", job.WalletAddress)
			return err
		}
		// Then add the DataFile to the Job.OutputFiles
		fmt.Println("Adding DataFile to Job.Outputs with CID:", dataFile.CID)
		job.OutputFiles = append(job.OutputFiles, dataFile)
		fmt.Println("Updated Job.Outputs:", job.OutputFiles)
		if err := db.Model(&user).Association("UserDatafiles").Append(&dataFile); err != nil {
			fmt.Printf("Error associating DataFile with User's UserDatafiles: %v\n", err)
			return err
		}
		fmt.Println("Updated User.UserDatafiles")
	}
	// Update job in the database with new OutputFiles
	if err := db.Save(&job).Error; err != nil {
		fmt.Printf("Error updating job: %v\n", err)
		return err
	}
	return nil
}

func submitBacalhauJobAndUpdateID(job *models.Job, db *gorm.DB, checkpointCompatible string) error {
	var inputs map[string]interface{}
	if err := json.Unmarshal(job.Inputs, &inputs); err != nil {
		return err
	}
	annotations := []string{}
	container := job.Tool.Container
	memory := job.Tool.Memory
	cpu := job.Tool.Cpu
	gpu := job.Tool.Gpu == 1
	network := job.Tool.Network
	flowuuid := job.Flow.FlowUUID
	jobuuid := job.JobUUID

	selector := ""

	bacalhauJob, err := bacalhau.CreateBacalhauJob(inputs, container, selector, maxComputeTime, memory, cpu, gpu, network, annotations, flowuuid, jobuuid, checkpointCompatible)
	if err != nil {
		return err
	}

	submittedBacalhauJob, err := bacalhau.SubmitBacalhauJob(bacalhauJob)
	if err != nil {
		return err
	}

	job.BacalhauJobID = submittedBacalhauJob.Metadata.ID
	fmt.Printf("Job had id %v\n", job.ID)
	fmt.Printf("Creating Job with bacalhau id %v\n", submittedBacalhauJob.Metadata.ID)
	return db.Save(job).Error
}
