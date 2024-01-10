package utils

import (
	"fmt"
	"time"

	"github.com/labdao/plex/internal/bacalhau"
	"github.com/labdao/plex/models"
	"gorm.io/gorm"
)

const maxQueueTime = 2 * time.Hour // Example maximum queue time

var db *gorm.DB // Assume db is initialized and available

func StartJobQueues() {
	fmt.Println("Starting job queues")

	go ProcessJobQueue("cpu")
	go ProcessJobQueue("gpu")
}

func ProcessJobQueue(jobType string) {
	for {
		// Fetch jobs from the database
		var jobs []models.Job
		result := db.Where("state = ? AND tool_id = ?", "queued", jobType).Find(&jobs)
		if result.Error != nil {
			// handle error
			time.Sleep(1 * time.Minute)
			continue
		}

		fmt.Printf("There are %d %s jobs in the queue\n", len(jobs), jobType)

		for _, job := range jobs {
			processJob(&job)
		}

		fmt.Printf("Finished processing queue, rehydrating %s queue with jobs\n", jobType)
		time.Sleep(5 * time.Minute) // Wait for some time before the next cycle
	}
}

func processJob(job *models.Job) {
	// Check if the job has timed out
	if time.Since(job.CreatedAt) > maxQueueTime {
		setJobStatus(job, "failed", "timed out in queue")
		return
	}

	// Submit job to Bacalhau (assuming function exists)
	submitBacalhauJob(job)

	// Check Bacalhau job status
	for {
		bacalhauJob, err := bacalhau.GetBacalhauJobState(job.BacalhauJobID)
		if err != nil {
			// handle error
			break
		}

		if bacalhauJob.State == "failed" {
			if time.Since(job.CreatedAt) > maxQueueTime {
				setJobStatus(job, "failed", "timed out in queue")
			} else {
				time.Sleep(1 * time.Second)
				submitBacalhauJob(job)
			}
		} else if bacalhauJob.State == "running" {
			setJobStatus(job, "running", "")
			break
		}

		time.Sleep(1 * time.Second) // Wait for a short period before checking the status again
	}
}

func setJobStatus(job *models.Job, status, errorMessage string) {
	job.State = status
	job.Error = errorMessage
	db.Save(job)
}

func submitBacalhauJob(job *models.Job) {
	// Function to submit the job to Bacalhau and update job.BacalhauJobID
	// This is a placeholder and needs to be replaced with actual submission logic
}
