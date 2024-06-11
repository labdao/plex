package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"encoding/hex"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/internal/ipwl"
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
	errChan := make(chan error, 4) // Buffer for two types of jobs

	go ProcessJobQueue(models.QueueTypeRay, errChan, db)

	// Wait for error from any queue
	for i := 0; i < 4; i++ {
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
			//TODO_PR#970: uncomment this later
			// fmt.Printf("There are no Jobs Queued for %v, will recheck later\n", queueType)
			time.Sleep(5 * time.Second)
			continue
		} else if err != nil {
			errChan <- err
			return
		}

		if queueType == models.QueueTypeRay {
			if err := processRayJob(job.ID, db); err != nil {
				errChan <- err
				return
			}
		}
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

func processRayJob(jobID uint, db *gorm.DB) error {
	fmt.Printf("Processing Ray Job %v\n", jobID)
	var job models.Job
	err := fetchJobWithToolAndFlowData(&job, jobID, db)
	if err != nil {
		return err
	}

	var ToolJson ipwl.Tool
	if err := json.Unmarshal(job.Tool.ToolJson, &ToolJson); err != nil {
		return err
	}

	if job.JobID == "" {
		if err := submitRayJobAndUpdateID(&job, db); err != nil {
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

func submitRayJobAndUpdateID(job *models.Job, db *gorm.DB) error {
	log.Println("Preparing to submit job to Ray service")
	var jobInputs map[string]interface{}
	if err := json.Unmarshal(job.Inputs, &jobInputs); err != nil {
		return err
	}
	inputs := make(map[string]interface{})
	for key, value := range jobInputs {
		inputs[key] = []interface{}{value}
	}
	toolCID := job.Tool.CID

	log.Printf("Submitting to Ray with inputs: %+v\n", inputs)
	setJobStatus(job, models.JobStateRunning, "", db)
	log.Printf("setting job %v to running\n", job.ID)
	resp, err := ray.SubmitRayJob(toolCID, inputs)
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
		// Print the pretty JSON
		fmt.Printf("Parsed Ray job response:\n%s\n", prettyJSON)
		completeRayJobAndAddFiles(job, body, rayJobResponse, db)
	} else {
		//create a fn to handle this
		job.State = models.JobStateFailed
	}
	//for now only handling 200 and failed. implement retry and timeout later

	fmt.Printf("Job had id %v\n", job.ID)
	fmt.Printf("Finished Job with Ray id %v and status %v\n", job.JobID, job.State)
	return db.Save(job).Error
}

func setJobStatus(job *models.Job, state models.JobState, errorMessage string, db *gorm.DB) error {
	job.State = state
	job.StartedAt = time.Now().UTC()
	job.Error = errorMessage
	return db.Save(job).Error
}

func completeRayJobAndAddFiles(job *models.Job, body []byte, resultJSON models.RayJobResponse, db *gorm.DB) error {

	//TODO_PR#970: the files are getting uploaded with the path from responseJSON for now. Need to change this to use the CID
	//for now only handling 200 and failed. implement retry and timeout later
	job.JobID = resultJSON.UUID
	job.ResultJSON = body
	job.State = models.JobStateCompleted

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
		Timestamp:     time.Now(),
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
	region := "us-east-1"
	accessKeyID := os.Getenv("BUCKET_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("BUCKET_SECRET_ACCESS_KEY")
	bucket, key, err := ray.GetBucketAndKeyFromURI(URI)
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
