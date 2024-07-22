package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"encoding/json"

	"github.com/google/uuid"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/labdao/plex/internal/ray"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RayQueue struct {
	db         *gorm.DB
	maxWorkers int
}

type WorkerState struct {
	ID         int
	Busy       bool
	CurrentJob *uint
	LastActive time.Time
}

var workerStates []*WorkerState

func NewRayQueue(db *gorm.DB, maxWorkers int) *RayQueue {
	return &RayQueue{
		db:         db,
		maxWorkers: maxWorkers,
	}
}

var maxRetryCountFor500 = GetEnvAsInt("MAX_RETRY_COUNT_FOR_500", 2)
var maxRetryCountFor404 = GetEnvAsInt("MAX_RETRY_COUNT_FOR_404", 1)

// Helper function to get a normally distributed random delay
// TODO_PR979: Make these env variables
func getRandomDelay(retryCount int, base time.Duration, factor float64, stdDev time.Duration) time.Duration {
	mean := float64(base.Nanoseconds()) * math.Pow(factor, float64(retryCount))
	delay := mean + rand.NormFloat64()*float64(stdDev.Nanoseconds())
	return time.Duration(delay)
}

var once sync.Once

func StartJobQueues(db *gorm.DB, maxWorkers int) error {
	once.Do(func() {
		rq := NewRayQueue(db, maxWorkers)
		rq.StartWorkers()
	})
	return nil
}

func (rq *RayQueue) StartWorkers() {
	workerStates = make([]*WorkerState, rq.maxWorkers) // initialize the slice based on maxWorkers
	for i := 0; i < rq.maxWorkers; i++ {
		workerStates[i] = &WorkerState{ID: i} // initialize each worker state
		go func(workerID int) {
			fmt.Printf("Starting worker %d\n", workerID)
			rq.worker(workerID)
		}(i)
	}
}

func (rq *RayQueue) worker(workerID int) {
	state := workerStates[workerID]
	for {
		// TEMP FIX: Marking dangling jobs to stopped after an hr. This wont be needed once we move to ray jobs
		err := fetchAndMarkOldestRunningJobAsStopped(rq.db)
		if err != nil {
			fmt.Printf("Error marking dangling jobs to stopped: %v\n", err)
		}
		var job models.Job
		err = fetchAndMarkOldestQueuedJobAsPending(&job, models.QueueTypeRay, rq.db)
		if err != nil {
			state.Busy = false
			state.CurrentJob = nil
			if errors.Is(err, gorm.ErrRecordNotFound) {
				time.Sleep(10 * time.Second)
				continue
			}
			fmt.Printf("Error fetching job: %v\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		state.Busy = true
		state.CurrentJob = &job.ID
		state.LastActive = time.Now()

		// Process the job
		if err = processRayJob(job.ID, rq.db); err != nil {
			fmt.Printf("Error processing job: %v\n", err)
		}

		state.Busy = false
		state.CurrentJob = nil
		time.Sleep(5 * time.Second) // Even after processing a job, sleep for a bit
	}
}

func GetWorkerSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workerStates)
}

func fetchAndMarkOldestRunningJobAsStopped(db *gorm.DB) error {
	var job models.Job
	err := db.Where("job_status = ?", models.JobStateRunning).Order("started_at ASC").First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	if time.Since(job.StartedAt) > 1*time.Hour {
		job.JobStatus = models.JobStateStopped
		return db.Save(&job).Error
	}

	return nil
}

func fetchAndMarkOldestQueuedJobAsPending(job *models.Job, queueType models.QueueType, db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("job_status = ?", models.JobStateQueued).Order("created_at ASC").First(job).Error; err != nil {
			return err
		}
		job.JobStatus = models.JobStatePending
		if err := tx.Save(job).Error; err != nil {
			return err
		}
		return nil
	})
}

// func fetchJobWithModelData(job *models.Job, id uint, db *gorm.DB) error {
// 	return db.Preload("Model").First(&job, id).Error
// }

func fetchJobWithModelAndExperimentData(job *models.Job, id uint, db *gorm.DB) error {
	return db.Preload("Model").Preload("Experiment").First(&job, id).Error
}

func fetchLatestInferenceEvent(inferenceEvent *models.InferenceEvent, jobID uint, db *gorm.DB) error {
	return db.Where("job_id = ?", jobID).Order("event_time DESC").First(inferenceEvent).Error
}

func processRayJob(jobID uint, db *gorm.DB) error {
	fmt.Printf("Processing Ray Job %v\n", jobID)
	var job models.Job
	err := fetchJobWithModelAndExperimentData(&job, jobID, db)
	if err != nil {
		return err
	}

	var inferenceEvent models.InferenceEvent
	err = fetchLatestInferenceEvent(&inferenceEvent, jobID, db)
	if err != nil {
		return err
	}

	var ModelJson ipwl.Model
	if err := json.Unmarshal(job.Model.ModelJson, &ModelJson); err != nil {
		return err
	}

	if job.RayJobID == "" {
		if err := submitRayJobAndUpdateID(&job, &inferenceEvent, db); err != nil {
			return err
		}
	}

	return nil
}

func UnmarshalRayJobResponse(data []byte) (models.RayJobResponse, error) {
	var response models.RayJobResponse
	var rawData map[string]interface{}

	if err := json.Unmarshal(data, &rawData); err != nil {
		return response, err
	}

	var responseData map[string]interface{}
	if _, ok := rawData["response"]; !ok {
		responseData = rawData
	} else {
		responseData = rawData["response"].(map[string]interface{})
	}

	response.UUID = responseData["uuid"].(string)
	if pdbData, ok := responseData["pdb"].(map[string]interface{}); ok {
		response.PDB = models.FileDetail{
			URI: pdbData["uri"].(string),
		}
	}

	response.Scores = make(map[string]float64)
	response.Files = make(map[string]models.FileDetail)

	if points, ok := responseData["points"].(float64); ok {
		response.Points = int(points)
	}

	// Function to recursively process map entries
	var processMap func(string, interface{})
	processMap = func(prefix string, value interface{}) {
		switch v := value.(type) {
		case map[string]interface{}:
			// Check if it's a file detail
			if uri, uriOk := v["uri"].(string); uriOk {
				response.Files[prefix] = models.FileDetail{URI: uri}
				return
			}
			// Otherwise, recursively process each field in the map
			for k, val := range v {
				newPrefix := k
				if prefix != "" {
					newPrefix = prefix + "." + k // To trace nested keys
				}
				processMap(newPrefix, val)
			}
		case []interface{}:
			// Process each item in the array
			for i, arrVal := range v {
				arrPrefix := fmt.Sprintf("%s[%d]", prefix, i)
				processMap(arrPrefix, arrVal)
			}
		case float64:
			// Handle scores which are float64
			response.Scores[prefix] = v
		}
	}

	// Initialize the recursive processing with an empty prefix
	for key, value := range responseData {
		if key == "uuid" || key == "pdb" || key == "points" {
			continue // Skip already processed or special handled fields
		}
		processMap(key, value)
	}

	return response, nil
}

func PrettyPrintRayJobResponse(response models.RayJobResponse) (string, error) {
	result := map[string]interface{}{
		"uuid":   response.UUID,
		"pdb":    response.PDB,
		"files":  response.Files,
		"scores": response.Scores,
	}

	prettyJSON, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return "", err
	}

	return string(prettyJSON), nil
}

func submitRayJobAndUpdateID(job *models.Job, inferenceEvent *models.InferenceEvent, db *gorm.DB) error {
	log.Println("Preparing to submit job to Ray service")
	var jobInputs map[string]interface{}
	if err := json.Unmarshal(job.Inputs, &jobInputs); err != nil {
		return err
	}
	inputs := make(map[string]interface{})
	for key, value := range jobInputs {
		inputs[key] = []interface{}{value}
	}
	modelPath := job.Model.S3URI

	rayJobID := uuid.New().String()
	log.Printf("Submitting to Ray with inputs: %+v\n", inputs)
	log.Printf("Here is the UUID for the job: %v\n", rayJobID)
	setJobStatusAndID(inferenceEvent, job, models.JobStateRunning, rayJobID, "", db)
	log.Printf("setting job %v to running\n", job.ID)
	resp, err := ray.SubmitRayJob(modelPath, rayJobID, inputs, db)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		var rayJobResponse models.RayJobResponse
		rayJobResponse, err = UnmarshalRayJobResponse([]byte(body))
		if err != nil {
			fmt.Println("Error unmarshalling result JSON:", err)
			return err
		}

		fmt.Printf("Parsed Ray job response: %+v\n", rayJobResponse)
		prettyJSON, err := PrettyPrintRayJobResponse(rayJobResponse)
		if err != nil {
			log.Fatalf("Error generating pretty JSON: %v", err)
		}

		fmt.Printf("Parsed Ray job response:\n%s\n", prettyJSON)
		completeRayJobAndAddFiles(inferenceEvent, job, body, rayJobResponse, db)
		fmt.Printf("Job %v completed and added files to DB\n", job.ID)
	} else {
		inferenceEvent.JobStatus = models.JobStateFailed
		inferenceEvent.ResponseCode = resp.StatusCode
		inferenceEvent.EventType = models.EventTypeJobFailed
		inferenceEvent.EventMessage = string(body)
		inferenceEvent.EventTime = time.Now().UTC()
		err = db.Save(inferenceEvent).Error
		if err != nil {
			return err
		}
		// 504 - no retry, 404 - immediate retry once, 500 - exponential back off and retry twice
		// TBD retry count and delay
		if resp.StatusCode == http.StatusGatewayTimeout {
			fmt.Printf("Job %v had timeout. Not retrying again. Marking as failed\n", job.ID)
			job.JobStatus = models.JobStateFailed
			err = db.Save(job).Error
			if err != nil {
				return err
			}
			return nil
		} else if (resp.StatusCode == http.StatusNotFound) && job.RetryCount < maxRetryCountFor404 {
			job.RetryCount = job.RetryCount + 1
			job.JobStatus = models.JobStateQueued
			err = db.Save(job).Error
			if err != nil {
				return err
			}
			newInferenceEvent := models.InferenceEvent{
				JobID:      job.ID,
				JobStatus:  models.JobStateQueued,
				RetryCount: job.RetryCount,
				EventTime:  time.Now().UTC(),
				EventType:  models.EventTypeJobCreated,
			}
			err = db.Save(&newInferenceEvent).Error
			if err != nil {
				return err
			}

			return submitRayJobAndUpdateID(job, &newInferenceEvent, db)
		} else if (resp.StatusCode == http.StatusInternalServerError) && job.RetryCount < maxRetryCountFor500 {
			job.RetryCount = job.RetryCount + 1
			job.JobStatus = models.JobStateQueued
			err = db.Save(job).Error
			if err != nil {
				return err
			}
			newInferenceEvent := models.InferenceEvent{
				JobID:      job.ID,
				JobStatus:  models.JobStateQueued,
				RetryCount: job.RetryCount,
				EventTime:  time.Now().UTC(),
				EventType:  models.EventTypeJobCreated,
			}
			err = db.Save(&newInferenceEvent).Error
			if err != nil {
				return err
			}
			// Calculate delay with exponential backoff and randomness
			delay := getRandomDelay(job.RetryCount, 2*time.Second, 1.2, 500*time.Millisecond)

			log.Printf("Retry after %v due to server error (500)", delay)
			time.Sleep(delay)

			return submitRayJobAndUpdateID(job, &newInferenceEvent, db)
		} else {
			var latestInferenceEvent models.InferenceEvent
			err = fetchLatestInferenceEvent(&latestInferenceEvent, job.ID, db)
			if err != nil {
				return err
			}
			fmt.Printf("The latest try of job %v had error %v. Retried %v times already. Marking as failed\n", job.ID, latestInferenceEvent.ResponseCode, job.RetryCount)
			job.JobStatus = models.JobStateFailed
			err = db.Save(job).Error
			if err != nil {
				return err
			}
		}

	}
	fmt.Printf("Job had id %v\n", job.ID)
	fmt.Printf("Finished Job with Ray id %v and status %v\n", job.RayJobID, job.JobStatus)
	err = db.Save(&job).Error
	if err != nil {
		return err
	}
	err = db.Save(&inferenceEvent).Error
	if err != nil {
		return err
	}
	return nil
}

func setJobStatusAndID(inferenceEvent *models.InferenceEvent, job *models.Job, state models.JobState, rayJobID string, errorMessage string, db *gorm.DB) error {
	inferenceEvent.RayJobID = rayJobID
	inferenceEvent.JobStatus = state
	inferenceEvent.EventTime = time.Now().UTC()
	inferenceEvent.EventType = models.EventTypeJobStarted
	err := db.Save(inferenceEvent).Error
	if err != nil {
		return err
	}

	job.JobStatus = state
	job.StartedAt = time.Now().UTC()
	job.Error = errorMessage
	job.RayJobID = rayJobID
	err = db.Save(job).Error
	if err != nil {
		return err
	}
	return nil
}

func completeRayJobAndAddFiles(inferenceEvent *models.InferenceEvent, job *models.Job, body []byte, resultJSON models.RayJobResponse, db *gorm.DB) error {

	//TODO_PR#970: the files are getting uploaded with the path from responseJSON for now. Need to change this to use the CID
	inferenceEvent.EventTime = time.Now().UTC()
	inferenceEvent.EventType = models.EventTypeJobCompleted
	inferenceEvent.JobStatus = models.JobStateCompleted
	inferenceEvent.ResponseCode = http.StatusOK
	inferenceEvent.OutputJson = body

	job.JobStatus = models.JobStateCompleted
	job.CompletedAt = time.Now().UTC()

	fmt.Printf("Looping through files in RayJobResponse\n %v\n", resultJSON.Files)
	// Iterate over all files in the RayJobResponse
	for key, fileDetail := range resultJSON.Files {
		fmt.Printf("AddFileToDB for file: %s, Key: %s\n", fileDetail.URI, key)
		if err := addFileToDB(job, fileDetail, key, db); err != nil {
			return fmt.Errorf("failed to add file (%s) to database: %v", key, err)
		}
	}

	fmt.Printf("Adding PDB file to DB\n %v\n", resultJSON.PDB)
	// Special handling for PDB as it's a common file across many jobs
	if err := addFileToDB(job, resultJSON.PDB, "pdb", db); err != nil {
		return fmt.Errorf("failed to add PDB file to database: %v", err)
	}

	return nil
}

func addFileToDB(job *models.Job, fileDetail models.FileDetail, fileType string, db *gorm.DB) error {
	fmt.Printf("Processing file: %s, Type: %s\n", fileDetail.URI, fileType)

	// Check if the file already exists
	var file models.File
	result := db.Where("s3_uri = ?", fileDetail.URI).First(&file)
	if result.Error == nil {
		fmt.Println("File already exists in DB:", fileDetail.URI)
		return nil // File already processed
	}

	// Handle not found error specifically
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error querying File table: %v", result.Error)
	}

	// Create new File entry
	tags := []models.Tag{
		{Name: fileType, Type: "filetype"},
		{Name: "generated", Type: "autogenerated"},
	}

	fmt.Printf("Creating new File record for %s \n", fileDetail.URI)
	file = models.File{
		WalletAddress: job.WalletAddress,
		Filename:      filepath.Base(fileDetail.URI),
		Tags:          tags,
		CreatedAt:     time.Now().UTC(),
		S3URI:         fileDetail.URI,
	}

	if err := db.Create(&file).Error; err != nil {
		return fmt.Errorf("error creating File record: %v", err)
	}

	// Associate file with job output
	job.OutputFiles = append(job.OutputFiles, file)
	if err := db.Save(&job).Error; err != nil {
		return fmt.Errorf("error updating job with new output file: %v", err)
	}

	return nil
}
