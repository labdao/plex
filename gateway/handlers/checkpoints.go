package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/internal/ipwl"

	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	s3client "github.com/labdao/plex/internal/s3"
)

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

type ExperimentListCheckpointsResult struct {
	JobID         int    `gorm:"column:job_id"`
	JobResultJson []byte `gorm:"column:result_json"`
	ModelJson     []byte `gorm:"column:model_json"`
}

func fetchJobCheckpoints(job models.Job) ([]map[string]string, error) {
	bucketEndpoint := os.Getenv("BUCKET_ENDPOINT")
	bucketName := os.Getenv("BUCKET_NAME")
	region := os.Getenv("AWS_REGION")
	accessKeyID := os.Getenv("BUCKET_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("BUCKET_SECRET_ACCESS_KEY")

	var presignedURLEndpoint string

	// Check if the bucket endpoint is the local development endpoint
	if bucketEndpoint == "http://object-store:9000" {
		presignedURLEndpoint = "http://localhost:9000"
	} else {
		presignedURLEndpoint = bucketEndpoint
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(presignedURLEndpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)
	var files []map[string]string

	var resultJSON models.RayJobResponse

	resultJSON, err = UnmarshalRayJobResponse([]byte(job.ResultJSON))
	if err != nil {
		fmt.Println("Error unmarshalling result JSON:", err)
		return nil, err
	}

	pdbKey := resultJSON.PDB.URI
	pdbFileName := filepath.Base(pdbKey)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(pdbKey),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return nil, err
	}

	files = append(files, map[string]string{
		"fileName": pdbFileName,
		"url":      urlStr,
	})

	return files, nil
}

func ListJobCheckpointsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		jobID := vars["jobID"]

		var job models.Job
		if err := db.Preload("Experiment").Preload("Model").First(&job, "id = ?", jobID).Error; err != nil {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}

		files, err := fetchJobCheckpoints(job)
		if err != nil {
			http.Error(w, "Failed to fetch checkpoints", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(files); err != nil {
			http.Error(w, "Failed to encode checkpoints", http.StatusInternalServerError)
			return
		}
	}
}

func ListExperimentCheckpointsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		experimentID := vars["experimentID"]

		var experiment models.Experiment
		if err := db.Preload("Jobs").First(&experiment, "id = ?", experimentID).Error; err != nil {
			http.Error(w, "Experiment not found", http.StatusNotFound)
			return
		}

		var allFiles []map[string]string

		for _, job := range experiment.Jobs {
			files, err := fetchJobCheckpoints(job)
			if err != nil {
				http.Error(w, "Failed to fetch checkpoints", http.StatusInternalServerError)
				return
			}
			allFiles = append(allFiles, files...)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(allFiles); err != nil {
			http.Error(w, "Failed to encode checkpoints", http.StatusInternalServerError)
			return
		}
	}
}

func fetchJobScatterPlotData(experimentListCheckpointsResult ExperimentListCheckpointsResult, db *gorm.DB) ([]models.ScatterPlotData, error) {
	bucketEndpoint := os.Getenv("BUCKET_ENDPOINT")
	bucketName := os.Getenv("BUCKET_NAME")
	region := os.Getenv("AWS_REGION")
	accessKeyID := os.Getenv("BUCKET_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("BUCKET_SECRET_ACCESS_KEY")

	var presignedURLEndpoint string

	// Check if the bucket endpoint is the local development endpoint
	if bucketEndpoint == "http://object-store:9000" {
		presignedURLEndpoint = "http://localhost:9000"
	} else {
		presignedURLEndpoint = bucketEndpoint
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(presignedURLEndpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)

	var ModelJson ipwl.Model
	if err := json.Unmarshal(experimentListCheckpointsResult.ModelJson, &ModelJson); err != nil {
		return nil, err
	}

	xAxis := ModelJson.XAxis
	yAxis := ModelJson.YAxis

	var resultJSON models.RayJobResponse
	// Unmarshal the result JSON from the job only when job.ResultJSON is not empty
	if string(experimentListCheckpointsResult.JobResultJson) != "" {
		resultJSON, err = UnmarshalRayJobResponse([]byte(experimentListCheckpointsResult.JobResultJson))
		if err != nil {
			fmt.Println("Error unmarshalling result JSON:", err)
			return nil, err
		}

		xAxisValue, xAxisExists := resultJSON.Scores[xAxis]
		yAxisValue, yAxisExists := resultJSON.Scores[yAxis]

		if !xAxisExists || !yAxisExists {
			return nil, fmt.Errorf("xAxis or yAxis value not found in the result JSON")
		}

		s3client, err := s3client.NewS3Client()
		if err != nil {
			return nil, err
		}

		_, key, err := s3client.GetBucketAndKeyFromURI(resultJSON.PDB.URI)
		if err != nil {
			return nil, err
		}
		pdbFileName := filepath.Base(key)

		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		urlStr, err := req.Presign(15 * time.Minute)
		if err != nil {
			return nil, err
		}

		plotData := []models.ScatterPlotData{}
		plotData = append(plotData, models.ScatterPlotData{
			Plddt:         xAxisValue,
			IPae:          yAxisValue,
			Checkpoint:    "0", // Default checkpoint
			StructureFile: pdbFileName,
			PdbFilePath:   urlStr,
			RayJobID:      resultJSON.UUID,
		})

		return plotData, nil
	} else {
		// If the job.ResultJSON is empty, return an empty array
		return []models.ScatterPlotData{}, nil
	}
}

func GetExperimentCheckpointDataHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		experimentID := vars["experimentID"]

		var experimentListCheckpointsResult []ExperimentListCheckpointsResult
		if err := db.Table("jobs").
			Select("jobs.id as job_id, models.model_json as model_json, jobs.result_json as result_json").
			Joins("join experiments on experiments.id = jobs.experiment_id").
			Joins("join models on models.cid = jobs.model_id").
			Where("experiments.id = ?", experimentID).
			Scan(&experimentListCheckpointsResult).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "No jobs found for the given experiment", http.StatusNotFound)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		var allPlotData []models.ScatterPlotData

		for _, job := range experimentListCheckpointsResult {
			plotData, err := fetchJobScatterPlotData(job, db) // Adjust function to accept the new job structure
			if err != nil {
				http.Error(w, "Failed to fetch scatter plot data for a job", http.StatusInternalServerError)
				return
			}
			allPlotData = append(allPlotData, plotData...)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allPlotData)
	}
}

//not used in the UI. but may be useful for API calls
// func fetchJobCheckpoints(job models.Job) ([]map[string]string, error) {
// 	bucketEndpoint := os.Getenv("BUCKET_ENDPOINT")
// 	bucketName := os.Getenv("BUCKET_NAME")
// 	region := "us-east-1"
// 	accessKeyID := os.Getenv("BUCKET_ACCESS_KEY_ID")
// 	secretAccessKey := os.Getenv("BUCKET_SECRET_ACCESS_KEY")

// 	var presignedURLEndpoint string

// 	// Check if the bucket endpoint is the local development endpoint
// 	if bucketEndpoint == "http://object-store:9000" {
// 		presignedURLEndpoint = "http://localhost:9000"
// 	} else {
// 		presignedURLEndpoint = bucketEndpoint
// 	}

// 	sess, err := session.NewSession(&aws.Config{
// 		Region:           aws.String(region),
// 		Endpoint:         aws.String(presignedURLEndpoint),
// 		S3ForcePathStyle: aws.Bool(true),
// 		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	svc := s3.New(sess)
// 	var files []map[string]string

// 	var resultJSON models.RayJobResponse

// 	resultJSON, err = UnmarshalRayJobResponse([]byte(job.ResultJSON))
// 	if err != nil {
// 		fmt.Println("Error unmarshalling result JSON:", err)
// 		return nil, err
// 	}

// 	pdbKey := resultJSON.PDB.URI
// 	pdbFileName := filepath.Base(pdbKey)

// 	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(pdbKey),
// 	})
// 	urlStr, err := req.Presign(15 * time.Minute)
// 	if err != nil {
// 		return nil, err
// 	}

// 	files = append(files, map[string]string{
// 		"fileName": pdbFileName,
// 		"url":      urlStr,
// 	})

// 	return files, nil
// }

//not used in the UI. but may be useful for API calls
// func ListJobCheckpointsHandler(db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		jobID := vars["jobID"]

// 		var job models.Job
// 		// if err := db.Joins("JOIN experiments ON experiments.id = jobs.experiment_id").
// 		// 			 Joins("JOIN models ON models.cid = jobs.model_id").
// 		// 			 Where("id = ?", jobID).
// 		// 			 First(&job).Error; err != nil {
// 		// 				if errors.Is(err, gorm.ErrRecordNotFound) {
// 		// 					http.Error(w, "Job not found", http.StatusNotFound)
// 		// 				} else {
// 		// 					http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		// 				}
// 		// 				return
// 		// }
// 		if err := db.First(&job, "id = ?", jobID).Error; err != nil {
// 			http.Error(w, "Job not found", http.StatusNotFound)
// 			return
// 		}

// 		files, err := fetchJobCheckpoints(job)
// 		if err != nil {
// 			http.Error(w, "Failed to fetch checkpoints", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)

// 		if err := json.NewEncoder(w).Encode(files); err != nil {
// 			http.Error(w, "Failed to encode checkpoints", http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

//not used in the UI. but may be useful for API calls
// func ListExperimentCheckpointsHandler(db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		experimentID := vars["experimentID"]

// 		var experiment models.Experiment
// 		if err := db.Preload("Jobs").First(&experiment, "id = ?", experimentID).Error; err != nil {
// 			http.Error(w, "Experiment not found", http.StatusNotFound)
// 			return
// 		}

// 		var allFiles []map[string]string

// 		for _, job := range experiment.Jobs {
// 			files, err := fetchJobCheckpoints(job)
// 			if err != nil {
// 				http.Error(w, "Failed to fetch checkpoints", http.StatusInternalServerError)
// 				return
// 			}
// 			allFiles = append(allFiles, files...)
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)

// 		if err := json.NewEncoder(w).Encode(allFiles); err != nil {
// 			http.Error(w, "Failed to encode checkpoints", http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

//not used in the UI. but may be useful for API calls
// func GetJobCheckpointDataHandler(db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		jobID := vars["jobID"]

// 		var job models.Job
// 		if err := db.Preload("Model").First(&job, "id = ?", jobID).Error; err != nil {
// 			http.Error(w, "Job not found", http.StatusNotFound)
// 			return
// 		}

// 		plotData, err := fetchJobScatterPlotData(job, db)
// 		if err != nil {
// 			http.Error(w, "Failed to fetch scatter plot data", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(plotData)
// 	}
// }
