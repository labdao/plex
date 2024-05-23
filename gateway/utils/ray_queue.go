package utils

import (
	"log"
	"net/http"
	"sync"

	"github.com/labdao/plex/internal/ray"
	"gorm.io/gorm"
)

type RayQueue struct {
	db         *gorm.DB
	maxWorkers int
	maxRetry   int
	running    int
	mutex      sync.Mutex
}

func NewRayQueue(db *gorm.DB, maxWorkers, maxRetry int) *RayQueue {
	return &RayQueue{
		db:         db,
		maxWorkers: maxWorkers,
		maxRetry:   maxRetry,
		running:    0,
	}
}

// TODO: move ray part of the queue from queue.go to ray_queue.go later
// func (rq *RayQueue) AddToQueue(toolPath string, inputs map[string]interface{}) error {
// 	// Create a new RayJob
// 	rayJob := &ray.RayJob{
// 		ToolPath: toolPath,
// 		Inputs:   inputs,
// 		Status:   ray.Queued,
// 		Retry:    0,
// 		JobUUID:  uuid.New().String(),
// 	}

// 	// Add the RayJob to the database
// 	err := rq.db.Create(rayJob).Error
// 	if err != nil {
// 		return err
// 	}

// 	// Process the queue
// 	rq.processTasks()

// 	return nil
// }

// func (rq *RayQueue) processTasks() {
// 	rq.mutex.Lock()
// 	defer rq.mutex.Unlock()

// 	for rq.running < rq.maxWorkers {
// 		var rayJob ray.RayJob
// 		err := rq.db.Where("status = ?", ray.Queued).First(&rayJob).Error
// 		if err != nil {
// 			log.Println("No jobs in the queue")
// 			return
// 		}

// 		// Set the status of the job to running
// 		rayJob.Status = ray.Running
// 		rq.db.Save(&rayJob)
// 		rq.running++

// 		// Submit the job to the Ray API
// 		go rq.runTask(&rayJob)
// 	}
// }

// func (rq *RayQueue) runTask(rayJob *ray.RayJob) {
// 	defer rq.completeTask()

// 	// Submit the job to the Ray API
// 	resp, err := rq.submitRayJob(rayJob.ToolPath, rayJob.Inputs)
// 	if err != nil {
// 		rq.handleRetry(rayJob, -1)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == http.StatusOK {
// 		rayJob.Status = ray.Complete
// 		rq.db.Save(rayJob)
// 	} else {
// 		rq.handleRetry(rayJob, resp.StatusCode)
// 	}
// }

// func (rq *RayQueue) completeTask() {
// 	rq.mutex.Lock()
// 	defer rq.mutex.Unlock()
// 	rq.running--
// 	rq.processTasks()
// }

// func (rq *RayQueue) handleRetry(rayJob *ray.RayJob, responseCode int) {
// 	rq.mutex.Lock()
// 	defer rq.mutex.Unlock()

// 	rayJob.Retry++
// 	if rayJob.Retry < rq.maxRetry {
// 		rayJob.Status = ray.Queued
// 		waitTime := time.Duration(rand.ExpFloat64()*float64(rayJob.Retry)) * time.Second
// 		if responseCode == http.StatusInternalServerError {
// 			time.Sleep(waitTime)
// 		}
// 		rq.db.Save(rayJob)
// 		rq.processTasks()
// 	} else {
// 		rayJob.Status = ray.Failed
// 		rq.db.Save(rayJob)
// 	}
// }

func (rq *RayQueue) submitRayJob(toolPath string, inputs map[string]interface{}) (*http.Response, error) {
	log.Printf("Creating Ray job with toolPath: %s and inputs: %+v\n", toolPath, inputs)
	resp, err := ray.CreateRayJob(toolPath, inputs)
	if err != nil {
		log.Printf("Error creating Ray job: %v\n", err)
		return nil, err
	}
	log.Printf("Ray job created with response status: %s\n", resp.Status)
	return resp, nil
}

// func (rq *RayQueue) validateInputKeys(inputVectors map[string]interface{}, toolInputs map[string]ray.ToolInput) error {
// 	for inputKey := range inputVectors {
// 		if _, exists := toolInputs[inputKey]; !exists {
// 			log.Printf("The argument %s is not in the tool inputs.\n", inputKey)
// 			log.Printf("Available keys: %v\n", toolInputs)
// 			return fmt.Errorf("the argument %s is not in the tool inputs", inputKey)
// 		}
// 	}
// 	return nil
// }
