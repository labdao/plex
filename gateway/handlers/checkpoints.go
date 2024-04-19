package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/models"

	"gorm.io/gorm"

	"encoding/csv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// listflowcheckpointshandler similar to the job one but deals with the flow level.
func ListFlowCheckpointsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		flowID := vars["flowID"]

		var flow models.Flow
		if err := db.First(&flow, "id = ?", flowID).Error; err != nil {
			http.Error(w, "Flow not found", http.StatusNotFound)
			return
		}
		flowUUID := flow.FlowUUID

		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		})
		if err != nil {
			http.Error(w, "Failed to create AWS session", http.StatusInternalServerError)
			return
		}

		svc := s3.New(sess)

		input := &s3.ListObjectsV2Input{
			Bucket: aws.String("app-checkpoint-bucket"),
			Prefix: aws.String("checkpoints/" + flowUUID + "/"),
		}

		result, err := svc.ListObjectsV2(input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var files []map[string]string
		for _, item := range result.Contents {
			trimmedKey := strings.TrimPrefix(*item.Key, "checkpoints/"+flowUUID+"/")

			if !strings.HasSuffix(trimmedKey, ".pdb") {
				continue
			}

			req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String("app-checkpoint-bucket"),
				Key:    item.Key,
			})
			urlStr, err := req.Presign(15 * time.Minute)
			if err != nil {
				http.Error(w, "Failed to sign request", http.StatusInternalServerError)
				return
			}

			files = append(files, map[string]string{
				"fileName": trimmedKey,
				"url":      urlStr,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(files); err != nil {
			http.Error(w, "Failed to encode checkpoints", http.StatusInternalServerError)
			return
		}
	}
}

func ListJobCheckpointsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		jobID := vars["jobID"]

		var job models.Job
		if err := db.Preload("Flow").First(&job, "id = ?", jobID).Error; err != nil {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		jobUUID := job.JobUUID
		flowUUID := job.Flow.FlowUUID

		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		})
		if err != nil {
			http.Error(w, "Failed to create AWS session", http.StatusInternalServerError)
			return
		}

		svc := s3.New(sess)

		input := &s3.ListObjectsV2Input{
			Bucket: aws.String("app-checkpoint-bucket"),
			Prefix: aws.String("checkpoints/" + flowUUID + "/" + jobUUID + "/"),
		}

		result, err := svc.ListObjectsV2(input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var files []map[string]string
		for _, item := range result.Contents {
			trimmedKey := strings.TrimPrefix(*item.Key, "checkpoints/"+flowUUID+"/"+jobUUID+"/")

			if !strings.HasSuffix(trimmedKey, ".pdb") {
				continue
			}

			req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String("app-checkpoint-bucket"),
				Key:    item.Key,
			})
			urlStr, err := req.Presign(15 * time.Minute)
			if err != nil {
				http.Error(w, "Failed to sign request", http.StatusInternalServerError)
				return
			}

			files = append(files, map[string]string{
				"fileName": trimmedKey,
				"url":      urlStr,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(files); err != nil {
			http.Error(w, "Failed to encode checkpoints", http.StatusInternalServerError)
			return
		}
	}
}

func AggregateCheckpointData(flowUUID string, db *gorm.DB) ([]models.ScatterPlotData, error) {
	var flows []models.Flow
	if err := db.Where("flow_uuid = ?", flowUUID).Preload("Jobs").Find(&flows).Error; err != nil {
		return nil, err
	}

	var allPlotData []models.ScatterPlotData
	for _, job := range flows[0].Jobs {
		plotData, err := AggregateJobCheckpointData(flowUUID, job.JobUUID)
		if err != nil {
			return nil, err
		}
		allPlotData = append(allPlotData, plotData...)
	}

	return allPlotData, nil
}

func AggregateJobCheckpointData(flowUUID string, jobUUID string) ([]models.ScatterPlotData, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	svc := s3.New(sess)

	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String("app-checkpoint-bucket"),
		Prefix: aws.String("checkpoints/" + flowUUID + "/" + jobUUID + "/"),
	}

	var plotData []models.ScatterPlotData
	err := svc.ListObjectsV2Pages(listInput, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, object := range page.Contents {
			if strings.HasSuffix(*object.Key, "summary.csv") {
				keyParts := strings.Split(*object.Key, "/")
				checkpointIndex := "0" // Default value in case parsing fails"
				for _, part := range keyParts {
					if strings.HasPrefix(part, "checkpoint_") {
						checkpointIndex = strings.TrimPrefix(part, "checkpoint_")
						break
					}
				}
				getObjectInput := &s3.GetObjectInput{
					Bucket: aws.String("app-checkpoint-bucket"),
					Key:    object.Key,
				}
				result, err := svc.GetObject(getObjectInput)
				if err != nil {
					return false
				}
				defer result.Body.Close()

				csvReader := csv.NewReader(result.Body)
				records, err := csvReader.ReadAll()
				if err != nil {
					return false
				}

				for _, record := range records[1:] {
					plddt, _ := strconv.ParseFloat(record[2], 64)
					i_pae, _ := strconv.ParseFloat(record[3], 64)
					pdbFileName := record[6]
					pdbPath := "checkpoints/" + flowUUID + "/" + jobUUID + "/checkpoint_" + checkpointIndex + "/" + pdbFileName
					req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
						Bucket: aws.String("app-checkpoint-bucket"),
						Key:    aws.String(pdbPath),
					})
					presignedURL, err := req.Presign(15 * time.Minute)
					if err != nil {
						return false
					}
					plotData = append(plotData, models.ScatterPlotData{Plddt: plddt, IPae: i_pae, Checkpoint: checkpointIndex, StructureFile: pdbFileName, PdbFilePath: presignedURL, JobUUID: jobUUID})
				}
			}
		}
		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return plotData, nil
}

func GetJobCheckpointDataHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		jobID := vars["jobID"]

		var job models.Job
		if err := db.Preload("Flow").First(&job, "id = ?", jobID).Error; err != nil {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}

		plotData, err := AggregateJobCheckpointData(job.Flow.FlowUUID, job.JobUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		if err := db.First(&flow, "id = ?", flowID).Error; err != nil {
			http.Error(w, "Flow not found", http.StatusNotFound)
			return
		}

		plotData, err := AggregateCheckpointData(flow.FlowUUID, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plotData)
	}
}

// StreamJobIDsForFlow streams the job IDs for a given flow to the client
func StreamJobIDsForFlow(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		flowID := vars["flowID"]

		var flow models.Flow
		if err := db.Preload("Jobs").First(&flow, "id = ?", flowID).Error; err != nil {
			http.Error(w, "Flow not found", http.StatusNotFound)
			return
		}

		encoder := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		// Stream each job ID to the client
		for _, job := range flow.Jobs {
			if err := encoder.Encode(job.ID); err != nil {
				http.Error(w, "Failed to stream job ID", http.StatusInternalServerError)
				return
			}
			flusher.Flush() // Send each job ID as it's processed
		}
	}
}
