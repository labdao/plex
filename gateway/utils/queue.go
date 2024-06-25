package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"encoding/hex"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/labdao/plex/internal/ray"
	s3client "github.com/labdao/plex/internal/s3"
	"gorm.io/gorm"
)

type RayQueue struct {
	db         *gorm.DB
	maxWorkers int
	jobChan    chan models.Job
	errChan    chan error
}

func NewRayQueue(db *gorm.DB, maxWorkers int) *RayQueue {
	return &RayQueue{
		db:         db,
		maxWorkers: maxWorkers,
		jobChan:    make(chan models.Job, maxWorkers),
		errChan:    make(chan error, maxWorkers),
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

func StartJobQueues(db *gorm.DB, maxWorkers int) error {
	rq := NewRayQueue(db, maxWorkers)
	rq.StartWorkers()

	go rq.ProcessJobQueue(models.QueueTypeRay)

	for err := range rq.errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (rq *RayQueue) StartWorkers() {
	for i := 0; i < rq.maxWorkers; i++ {
		go rq.worker()
	}
}

func (rq *RayQueue) worker() {
	for job := range rq.jobChan {
		err := processRayJob(job.ID, rq.db)
		if err != nil {
			rq.errChan <- err
		}
	}
}

func (rq *RayQueue) ProcessJobQueue(queueType models.QueueType) {
	for {
		var job models.Job
		err := fetchOldestQueuedJob(&job, queueType, rq.db)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			time.Sleep(5 * time.Second)
			continue
		} else if err != nil {
			rq.errChan <- err
			return
		}

		rq.jobChan <- job
	}
}

func fetchOldestQueuedJob(job *models.Job, queueType models.QueueType, db *gorm.DB) error {
	return db.Where("state = ? AND queue = ?", models.JobStateQueued, queueType).Order("created_at ASC").First(job).Error
}

// func fetchJobWithModelData(job *models.Job, id uint, db *gorm.DB) error {
// 	return db.Preload("Model").First(&job, id).Error
// }

func fetchJobWithModelAndExperimentData(job *models.Job, id uint, db *gorm.DB) error {
	return db.Preload("Model").Preload("Experiment").First(&job, id).Error
}

func fetchLatestRequestTracker(requestTracker *models.RequestTracker, jobID uint, db *gorm.DB) error {
	return db.Where("job_id = ?", jobID).Order("created_at DESC").First(requestTracker).Error
}

func processRayJob(jobID uint, db *gorm.DB) error {
	fmt.Printf("Processing Ray Job %v\n", jobID)
	var job models.Job
	err := fetchJobWithModelAndExperimentData(&job, jobID, db)
	if err != nil {
		return err
	}

	var requestTracker models.RequestTracker
	err = fetchLatestRequestTracker(&requestTracker, jobID, db)
	if err != nil {
		return err
	}

	var ModelJson ipwl.Model
	if err := json.Unmarshal(job.Model.ModelJson, &ModelJson); err != nil {
		return err
	}

	if job.RayJobID == "" {
		if err := submitRayJobAndUpdateID(&job, &requestTracker, db); err != nil {
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

	responseData := rawData["response"].(map[string]interface{})

	response.UUID = responseData["uuid"].(string)
	if pdbData, ok := responseData["pdb"].(map[string]interface{}); ok {
		response.PDB = models.FileDetail{
			URI: pdbData["uri"].(string),
		}
	}

	response.Scores = make(map[string]float64)
	response.Files = make(map[string]models.FileDetail)

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
		if key == "uuid" || key == "pdb" {
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

func submitRayJobAndUpdateID(job *models.Job, requestTracker *models.RequestTracker, db *gorm.DB) error {
	log.Println("Preparing to submit job to Ray service")
	var jobInputs map[string]interface{}
	if err := json.Unmarshal(job.Inputs, &jobInputs); err != nil {
		return err
	}
	inputs := make(map[string]interface{})
	for key, value := range jobInputs {
		inputs[key] = []interface{}{value}
	}
	modelCID := job.Model.CID

	rayJobID := uuid.New().String()
	log.Printf("Submitting to Ray with inputs: %+v\n", inputs)
	log.Printf("Here is the UUID for the job: %v\n", rayJobID)
	setJobStatusAndID(requestTracker, job, models.JobStateRunning, rayJobID, "", db)
	log.Printf("setting job %v to running\n", job.ID)
	resp, err := ray.SubmitRayJob(modelCID, rayJobID, inputs, db)
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
		completeRayJobAndAddFiles(requestTracker, job, body, rayJobResponse, db)
	} else {
		requestTracker.State = models.JobStateFailed
		requestTracker.ResponseCode = resp.StatusCode
		requestTracker.ErrorMessage = string(body)
		requestTracker.CompletedAt = time.Now().UTC()
		err = db.Save(requestTracker).Error
		if err != nil {
			return err
		}
		// 504 - no retry, 404 - immediate retry once, 500 - exponential back off and retry twice
		// TBD retry count and delay
		if resp.StatusCode == http.StatusGatewayTimeout {
			fmt.Printf("Job %v had timeout. Not retrying again. Marking as failed\n", job.ID)
			job.State = models.JobStateFailed
			err = db.Save(job).Error
			if err != nil {
				return err
			}
			return nil
		} else if (resp.StatusCode == http.StatusNotFound) && job.RetryCount < maxRetryCountFor404 {
			job.RetryCount = job.RetryCount + 1
			job.State = models.JobStateQueued
			err = db.Save(job).Error
			if err != nil {
				return err
			}
			newRequestTracker := models.RequestTracker{
				JobID:      job.ID,
				State:      models.JobStateQueued,
				RetryCount: job.RetryCount,
				CreatedAt:  time.Now().UTC(),
			}
			err = db.Save(&newRequestTracker).Error
			if err != nil {
				return err
			}

			return submitRayJobAndUpdateID(job, &newRequestTracker, db)
		} else if (resp.StatusCode == http.StatusInternalServerError) && job.RetryCount < maxRetryCountFor500 {
			job.RetryCount = job.RetryCount + 1
			job.State = models.JobStateQueued
			err = db.Save(job).Error
			if err != nil {
				return err
			}
			newRequestTracker := models.RequestTracker{
				JobID:      job.ID,
				State:      models.JobStateQueued,
				RetryCount: job.RetryCount,
				CreatedAt:  time.Now().UTC(),
			}
			err = db.Save(&newRequestTracker).Error
			if err != nil {
				return err
			}
			// Calculate delay with exponential backoff and randomness
			delay := getRandomDelay(job.RetryCount, 2*time.Second, 1.2, 500*time.Millisecond)

			log.Printf("Retry after %v due to server error (500)", delay)
			time.Sleep(delay)

			return submitRayJobAndUpdateID(job, &newRequestTracker, db)
		} else {
			var latestRequestTracker models.RequestTracker
			err = fetchLatestRequestTracker(&latestRequestTracker, job.ID, db)
			if err != nil {
				return err
			}
			fmt.Printf("The latest try of job %v had error %v. Retried %v times already. Marking as failed\n", job.ID, latestRequestTracker.ResponseCode, job.RetryCount)
			job.State = models.JobStateFailed
			err = db.Save(job).Error
			if err != nil {
				return err
			}
		}

	}
	fmt.Printf("Job had id %v\n", job.ID)
	fmt.Printf("Finished Job with Ray id %v and status %v\n", job.RayJobID, job.State)
	err = db.Save(job).Error
	if err != nil {
		return err
	}
	err = db.Save(requestTracker).Error
	if err != nil {
		return err
	}
	return nil
}

func setJobStatusAndID(requestTracker *models.RequestTracker, job *models.Job, state models.JobState, rayJobID string, errorMessage string, db *gorm.DB) error {
	requestTracker.RayJobID = rayJobID
	requestTracker.State = state
	requestTracker.StartedAt = time.Now().UTC()
	err := db.Save(requestTracker).Error
	if err != nil {
		return err
	}

	job.State = state
	job.StartedAt = time.Now().UTC()
	job.Error = errorMessage
	job.RayJobID = rayJobID
	err = db.Save(job).Error
	if err != nil {
		return err
	}
	return nil
}

func completeRayJobAndAddFiles(requestTracker *models.RequestTracker, job *models.Job, body []byte, resultJSON models.RayJobResponse, db *gorm.DB) error {

	//TODO_PR#970: the files are getting uploaded with the path from responseJSON for now. Need to change this to use the CID
	requestTracker.CompletedAt = time.Now().UTC()
	requestTracker.State = models.JobStateCompleted
	requestTracker.ResponseCode = http.StatusOK
	requestTracker.JobResponse = body

	job.ResultJSON = body
	job.State = models.JobStateCompleted
	job.CompletedAt = time.Now().UTC()

	// Iterate over all files in the RayJobResponse
	for key, fileDetail := range resultJSON.Files {
		if err := addFileToDB(job, fileDetail, key, db); err != nil {
			return fmt.Errorf("failed to add file (%s) to database: %v", key, err)
		}
	}

	// Special handling for PDB as it's a common file across many jobs
	if err := addFileToDB(job, resultJSON.PDB, "pdb", db); err != nil {
		return fmt.Errorf("failed to add PDB file to database: %v", err)
	}

	return nil
}

func addFileToDB(job *models.Job, fileDetail models.FileDetail, fileType string, db *gorm.DB) error {
	fmt.Printf("Processing file: %s, Type: %s\n", fileDetail.URI, fileType)

	// Check if the file already exists
	var dataFile models.DataFile
	result := db.Where("s3_uri = ?", fileDetail.URI).First(&dataFile)
	if result.Error == nil {
		fmt.Println("File already exists in DB:", fileDetail.URI)
		return nil // File already processed
	}

	// Handle not found error specifically
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error querying DataFile table: %v", result.Error)
	}

	// Create new DataFile entry
	tags := []models.Tag{
		{Name: fileType, Type: "filetype"},
		{Name: "generated", Type: "autogenerated"},
	}

	hash, err := hashS3Object(fileDetail.URI)
	if err != nil {
		return fmt.Errorf("error hashing S3 object: %v", err)
	}

	dataFile = models.DataFile{
		CID:           hash,
		WalletAddress: job.WalletAddress,
		Filename:      filepath.Base(fileDetail.URI),
		Tags:          tags,
		Timestamp:     time.Now().UTC(),
		S3URI:         fileDetail.URI,
	}

	if err := db.Create(&dataFile).Error; err != nil {
		return fmt.Errorf("error creating DataFile record: %v", err)
	}

	// Associate file with job output
	job.OutputFiles = append(job.OutputFiles, dataFile)
	if err := db.Save(&job).Error; err != nil {
		return fmt.Errorf("error updating job with new output file: %v", err)
	}

	return nil
}

func hashS3Object(URI string) (string, error) {
	// Load the AWS configuration
	bucketEndpoint := os.Getenv("BUCKET_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	accessKeyID := os.Getenv("BUCKET_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("BUCKET_SECRET_ACCESS_KEY")
	s3client, err := s3client.NewS3Client()
	if err != nil {
		return "", fmt.Errorf("error creating S3 client: %w", err)
	}
	bucket, key, err := s3client.GetBucketAndKeyFromURI(URI)
	if err != nil {
		return "", fmt.Errorf("error parsing S3 URI: %w", err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(bucketEndpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return "", fmt.Errorf("error creating AWS session: %w", err)
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// Get the object
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get S3 object: %w", err)
	}
	defer resp.Body.Close()

	// Initialize hasher and hash the key
	hasher := sha256.New()
	if _, err := hasher.Write([]byte(key)); err != nil {
		return "", fmt.Errorf("failed to hash key: %w", err)
	}

	// Read and hash the contents of the object
	if _, err := io.Copy(hasher, resp.Body); err != nil {
		return "", fmt.Errorf("failed to hash file contents: %w", err)
	}

	// Compute the final hash
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}
