package utils

import (
	"fmt"
	"time"

	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/internal/bacalhau"
	"github.com/labdao/plex/models"
	"gorm.io/gorm"
)

const maxQueueTime = 2 * time.Hour // Example maximum queue time

var db *gorm.DB // Assume db is initialized and available

func StartJobQueues() error {
	errChan := make(chan error, 2) // Buffer for two types of jobs

	go ProcessJobQueue(models.QueueTypeCPU, errChan)
	go ProcessJobQueue(models.QueueTypeGPU, errChan)

	// Wait for error from any queue
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

func ProcessJobQueue(queueType models.QueueType, errChan chan<- error) {
	for {
		var jobs []models.Job
		if err := fetchJobs(&jobs, queueType); err != nil {
			errChan <- err
			return
		}

		fmt.Printf("There are %d jobs in the %s queue\n", len(jobs), queueType)

		for _, job := range jobs {
			// there should not be any errors from processJob
			if err := processJob(&job); err != nil {
				errChan <- err
				return
			}
		}

		fmt.Printf("Finished processing queue, rehydrating %s queue with jobs\n", jobType)
		time.Sleep(1 * time.Second) // Wait for some time before the next cycle
	}
}

func fetchJobs(jobs *[]models.Job, queueType models.QueueType) error {
	return db.Where("state = ? AND queue = ?", models.JobStateQueued, queueType).Find(jobs).Error
}

func processJob(job *models.Job) error {
	if time.Since(job.CreatedAt) > maxQueueTime {
		return setJobStatus(job, models.JobStateFailed, "timed out in queue")
	}

	if err := submitBacalhauJob(job); err != nil {
		return err
	}

	for {
		bacalhauJob, err := bacalhau.GetBacalhauJobState(job.BacalhauJobID)
		if err != nil {
			return err
		}

		if bacalhauJob.State == "failed" {
			if time.Since(job.CreatedAt) > maxQueueTime {
				return setJobStatus(job, "failed", "timed out in queue")
			}
			time.Sleep(1 * time.Second)
			if err := submitBacalhauJob(job); err != nil {
				return err
			}
		} else if bacalhauJob.State == "running" {
			return setJobStatus(job, "running", "")
		}

		time.Sleep(1 * time.Second) // Wait for a short period before checking the status again
	}
}

func setJobStatus(job *models.Job, state models.JobState, errorMessage string) error {
	job.State = state
	job.Error = errorMessage
	return db.Save(job).Error
}

func submitBacalhauJob(job *models.Job) error {
	// Define inputs and other parameters required for CreateBacalhauJob
	// This is a placeholder and needs actual implementation based on job details and requirements

	createdJob, err := bacalhau.CreateBacalhauJob(inputs, container, selector, maxTime, memory, cpu, gpu, network, annotations)
	if err != nil {
		return err
	}

	job.BacalhauJobID = createdJob.ID
	return db.Save(job).Error
}
