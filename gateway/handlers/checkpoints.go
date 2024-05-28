package handlers

import (
	"encoding/json"
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
)

func fetchJobCheckpoints(job models.Job) ([]map[string]string, error) {
	bucketEndpoint := os.Getenv("BUCKET_ENDPOINT")
	bucketName := os.Getenv("BUCKET_NAME")
	region := "us-east-1"
	accessKeyID := os.Getenv("BUCKET_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("BUCKET_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(bucketEndpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)
	var files []map[string]string

	var resultJSON models.RayJobResponse
	if err := json.Unmarshal([]byte(job.ResultJSON), &resultJSON); err != nil {
		return nil, err
	}

	pdbKey := resultJSON.PDB.Key
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
		if err := db.Preload("Flow").Preload("Tool").First(&job, "id = ?", jobID).Error; err != nil {
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

func ListFlowCheckpointsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		flowID := vars["flowID"]

		var flow models.Flow
		if err := db.Preload("Jobs").First(&flow, "id = ?", flowID).Error; err != nil {
			http.Error(w, "Flow not found", http.StatusNotFound)
			return
		}

		var allFiles []map[string]string

		for _, job := range flow.Jobs {
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

func GetDynamicFieldValue(response models.RayJobResponse, field string) (interface{}, bool) {
	value, exists := response.DynamicFields[field]
	return value, exists
}

func fetchJobScatterPlotData(job models.Job, db *gorm.DB) ([]models.ScatterPlotData, error) {
	bucketEndpoint := os.Getenv("BUCKET_ENDPOINT")
	bucketName := os.Getenv("BUCKET_NAME")
	region := "us-east-1"
	accessKeyID := os.Getenv("BUCKET_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("BUCKET_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(bucketEndpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)

	var resultMap map[string]interface{}
	if err := json.Unmarshal([]byte(job.ResultJSON), &resultMap); err != nil {
		return nil, err
	}

	// Create a RayJobResponse and populate known fields
	var resultJSON models.RayJobResponse
	knownFieldsData, err := json.Marshal(resultMap)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(knownFieldsData, &resultJSON); err != nil {
		return nil, err
	}

	// Remove known fields from the map
	delete(resultMap, "uuid")
	delete(resultMap, "pdb")
	delete(resultMap, "structure_metrics")
	delete(resultMap, "plots")
	delete(resultMap, "msa")

	// Assign the remaining fields to DynamicFields
	resultJSON.DynamicFields = resultMap

	var ToolJson ipwl.Tool
	if err := json.Unmarshal(job.Tool.ToolJson, &ToolJson); err != nil {
		return nil, err
	}

	xAxis := ToolJson.XAxis
	yAxis := ToolJson.YAxis

	xAxisValue, xAxisExists := GetDynamicFieldValue(resultJSON, xAxis)
	yAxisValue, yAxisExists := GetDynamicFieldValue(resultJSON, yAxis)

	if !xAxisExists || !yAxisExists {
		return nil, fmt.Errorf("xAxis or yAxis value not found in the result JSON")
	}

	// Default to 0 if xAxis or yAxis value is nil
	var xAxisFloat, yAxisFloat float64
	var xAxisIsFloat, yAxisIsFloat bool

	if xAxisValue == nil {
		xAxisFloat = 0
	} else {
		xAxisFloat, xAxisIsFloat = xAxisValue.(float64)
		if !xAxisIsFloat {
			return nil, fmt.Errorf("xAxis value is not a float64")
		}
	}

	if yAxisValue == nil {
		yAxisFloat = 0
	} else {
		yAxisFloat, yAxisIsFloat = yAxisValue.(float64)
		if !yAxisIsFloat {
			return nil, fmt.Errorf("yAxis value is not a float64")
		}
	}

	pdbKey := resultJSON.PDB.Key
	pdbFileName := filepath.Base(pdbKey)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(pdbKey),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return nil, err
	}

	plotData := []models.ScatterPlotData{}
	plotData = append(plotData, models.ScatterPlotData{
		Plddt:         xAxisFloat,
		IPae:          yAxisFloat,
		Checkpoint:    "0", // Default checkpoint
		StructureFile: pdbFileName,
		PdbFilePath:   urlStr,
		JobUUID:       resultJSON.UUID,
	})

	return plotData, nil
}

func GetJobCheckpointDataHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		jobID := vars["jobID"]

		var job models.Job
		if err := db.Preload("Flow").Preload("Tool").First(&job, "id = ?", jobID).Error; err != nil {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}

		plotData, err := fetchJobScatterPlotData(job, db)
		if err != nil {
			http.Error(w, "Failed to fetch scatter plot data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plotData)
	}
}

func GetFlowCheckpointDataHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		flowID := vars["flowID"]

		var flow models.Flow
		if err := db.Preload("Jobs.Tool").First(&flow, "id = ?", flowID).Error; err != nil {
			http.Error(w, "Flow not found", http.StatusNotFound)
			return
		}

		var allPlotData []models.ScatterPlotData

		for _, job := range flow.Jobs {
			plotData, err := fetchJobScatterPlotData(job, db)
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
