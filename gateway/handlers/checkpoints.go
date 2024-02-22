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

func ListCheckpointsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		jobID := vars["jobID"]

		var job models.Job
		if err := db.First(&job, "id = ?", jobID).Error; err != nil {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		jobUUID := job.JobUUID

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
			Prefix: aws.String("checkpoints/" + jobUUID + "/"),
		}

		result, err := svc.ListObjectsV2(input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var files []map[string]string
		for _, item := range result.Contents {
			trimmedKey := strings.TrimPrefix(*item.Key, "checkpoints/"+jobUUID+"/")

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

func AggregateCheckpointData(jobUUID string) ([]models.ScatterPlotData, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	svc := s3.New(sess)

	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String("app-checkpoint-bucket"),
		Prefix: aws.String("checkpoints/" + jobUUID + "/"),
	}

	var plotData []models.ScatterPlotData
	err := svc.ListObjectsV2Pages(listInput, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, object := range page.Contents {
			if strings.HasSuffix(*object.Key, "event.csv") {
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

				for index, record := range records[1:] {
					factor1, _ := strconv.ParseFloat(record[2], 64)
					factor2, _ := strconv.ParseFloat(record[3], 64)
					pdbFileName := record[6]
					pdbPath := "checkpoints/" + jobUUID + "/checkpoint_" + strconv.Itoa(index) + "/" + pdbFileName
					req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
						Bucket: aws.String("app-checkpoint-bucket"),
						Key:    aws.String(pdbPath),
					})
					presignedURL, err := req.Presign(15 * time.Minute)
					if err != nil {
						return false
					}
					plotData = append(plotData, models.ScatterPlotData{Factor1: factor1, Factor2: factor2, PdbFilePath: presignedURL})
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

func GetCheckpointDataHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		jobID := vars["jobID"]

		var job models.Job
		if err := db.First(&job, "id = ?", jobID).Error; err != nil {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}

		plotData, err := AggregateCheckpointData(job.JobUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plotData)
	}
}
