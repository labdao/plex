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
	"github.com/aws/aws-sdk-go/service/s3"
	s3client "github.com/labdao/plex/internal/s3"
)

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

func fetchJobScatterPlotData(experimentListCheckpointsResult ExperimentListCheckpointsResult, db *gorm.DB) ([]models.ScatterPlotData, error) {
	bucketName := os.Getenv("BUCKET_NAME")
	var ModelJson ipwl.Model
	if err := json.Unmarshal(experimentListCheckpointsResult.ModelJson, &ModelJson); err != nil {
		return nil, err
	}

	xAxis := ModelJson.XAxis
	yAxis := ModelJson.YAxis
	// Unmarshal the result JSON from the job only when job.ResultJSON is not empty
	if string(experimentListCheckpointsResult.JobResultJson) != "" {
		resultJSON, err := UnmarshalRayJobResponse([]byte(experimentListCheckpointsResult.JobResultJson))
		if err != nil {
			fmt.Println("Error unmarshalling result JSON:", err)
			return nil, err
		}

		xAxisValue, xAxisExists := resultJSON.Scores[xAxis]
		yAxisValue, yAxisExists := resultJSON.Scores[yAxis]

		if !xAxisExists || !yAxisExists {
			return nil, fmt.Errorf("xAxis or yAxis value not found in the result JSON")
		}

		s3client, err := s3client.NewS3Client(true)
		if err != nil {
			return nil, err
		}

		_, key, err := s3client.GetBucketAndKeyFromURI(resultJSON.PDB.URI)
		if err != nil {
			return nil, err
		}
		pdbFileName := filepath.Base(key)

		req, _ := s3client.Client.GetObjectRequest(&s3.GetObjectInput{
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

		// select jobs.id as job_id, models.model_json as model_json, inference_events.output_json as result_json
		// from jobs
		// JOIN inference_events ON inference_events.job_id = jobs.id and inference_events.event_type = 'file_processed'
		// JOIN experiments ON experiments.id = jobs.experiment_id
		// JOIN models ON models.id = jobs.model_id
		// where experiments.id = 5;
		if err := db.Table("jobs").
			Select("jobs.id as job_id, models.model_json as model_json, inference_events.output_json as result_json").
			Joins("JOIN inference_events ON inference_events.job_id = jobs.id AND inference_events.file_name is not null AND inference_events.event_type = ?", models.EventTypeFileProcessed).
			Joins("JOIN experiments ON experiments.id = jobs.experiment_id").
			Joins("JOIN models ON models.id = jobs.model_id").
			Where("experiments.id = ?", experimentID).
			Scan(&experimentListCheckpointsResult).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "No jobs found for the given experiment", http.StatusNotFound)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
		// var experimentListCheckpointsResult []ExperimentListCheckpointsResult

		// if err := db.Table("jobs").
		// 	Select("jobs.id as job_id, models.model_json as model_json, inference_events.result_json as result_json").
		// 	Joins("join (select job_id, max(event_time) as max_created_at from inference_events group by job_id) as latest_events on latest_events.job_id = inference_events.job_id and latest_events.max_created_at = inference_events.created_at").
		// 	Joins("JOIN inference_events ON inference_events.job_id = latest_events.job_id AND inference_events.event_time = latest_events.max_created_at").
		// 	Joins("join experiments on experiments.id = jobs.experiment_id").
		// 	Joins("join models on models.id = jobs.model_id").
		// 	Where("experiments.id = ?", experimentID).
		// 	Scan(&experimentListCheckpointsResult).Error; err != nil {
		// 	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 		http.Error(w, "No jobs found for the given experiment", http.StatusNotFound)
		// 	} else {
		// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
		// 	}
		// 	return
		// }

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
