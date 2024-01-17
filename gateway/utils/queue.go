package utils

import (
	"fmt"
	"time"

	"encoding/json"

	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/internal/bacalhau"
	"gorm.io/gorm"
)

const maxQueueTime = 5 * time.Minute

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
		var jobs []models.Job
		if err := fetchJobsWithToolData(&jobs, queueType, db); err != nil {
			errChan <- err
			return
		}

		fmt.Printf("There are %d jobs in the %s queue\n", len(jobs), queueType)

		for _, job := range jobs {
			// there should not be any errors from processJob
			if err := processJob(&job, db); err != nil {
				errChan <- err
				return
			}
		}

		fmt.Printf("Finished processing queue, rehydrating %s queue with jobs\n", queueType)
		time.Sleep(5 * time.Second) // Wait for some time before the next cycle
	}
}

func fetchJobsWithToolData(jobs *[]models.Job, queueType models.QueueType, db *gorm.DB) error {
	return db.Preload("Tool").Where("state = ? AND queue = ?", models.JobStateQueued, queueType).Find(jobs).Error
}

func processJob(job *models.Job, db *gorm.DB) error {
	fmt.Printf("Processing job %v\n", job.ID)
	if err := submitBacalhauJobAndUpdateID(job, db); err != nil {
		return err
	}

	for {
		fmt.Printf("Checking status for %v , %v\n", job.ID, job.BacalhauJobID)
		if time.Since(job.CreatedAt) > maxQueueTime {
			fmt.Printf("Job %v , %v timed out\n", job.ID, job.BacalhauJobID)
			return setJobStatus(job, models.JobStateFailed, "timed out in queue", db)
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
			time.Sleep(3 * time.Second) // Wait for a short period before checking the status again
			if err := submitBacalhauJobAndUpdateID(job, db); err != nil {
				return err
			}
			continue
		} else if bacalhau.JobIsRunning(bacalhauJob) {
			fmt.Printf("Job %v , %v is running\n", job.ID, job.BacalhauJobID)
			return setJobStatus(job, models.JobStateRunning, "", db)
		} else if bacalhau.JobFailed(bacalhauJob) {
			fmt.Printf("Job %v , %v failed\n", job.ID, job.BacalhauJobID)
			return setJobStatus(job, models.JobStateFailed, "", db)
		} else if bacalhau.JobCompleted(bacalhauJob) {
			fmt.Printf("Job %v , %v completed\n", job.ID, job.BacalhauJobID)
			return setJobStatus(job, models.JobStateCompleted, "", db)
		} else {
			fmt.Printf("Job %v , %v has state %v, will requery\n", job.ID, job.BacalhauJobID, bacalhauJob.State.State)
			time.Sleep(3 * time.Second) // Wait for a short period before checking the status again
		}
	}
}

func setJobStatus(job *models.Job, state models.JobState, errorMessage string, db *gorm.DB) error {
	job.State = state
	job.Error = errorMessage
	return db.Save(job).Error
}

func submitBacalhauJobAndUpdateID(job *models.Job, db *gorm.DB) error {
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

	selector := ""
	maxTime := 60 * 72

	bacalhauJob, err := bacalhau.CreateBacalhauJob(inputs, container, selector, maxTime, memory, cpu, gpu, network, annotations)
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
