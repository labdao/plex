package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"gorm.io/gorm"
)

func GetJobHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
		if !ok {
			http.Error(w, "User context not found", http.StatusUnauthorized)
			return
		}

		params := mux.Vars(r)
		jobID, err := strconv.Atoi(params["jobID"])
		if err != nil {
			http.Error(w, fmt.Sprintf("Job ID (%v) could not be converted to int", params["jobID"]), http.StatusNotFound)
		}

		var job models.Job
		query := db.Preload("OutputFiles.Tags").Preload("InputFiles.Tags").Where("id = ?", jobID)

		if result := query.First(&job); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Job not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Job: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		if !job.Public && job.WalletAddress != user.WalletAddress && !user.Admin {
			http.Error(w, "Job not found or not authorized", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(job); err != nil {
			http.Error(w, "Error encoding Job to JSON", http.StatusInternalServerError)
			return
		}
	}
}

type JobSummary struct {
	Count         int     `json:"count"`
	TotalCpu      float64 `json:"totalCpu"`
	TotalMemoryGb int     `json:"totalMemoryGb"`
	TotalGpu      int     `json:"totalGpu"`
}

type QueueSummary struct {
	Queued  JobSummary `json:"queued"`
	Running JobSummary `json:"running"`
}

type Summary struct {
	Ray QueueSummary `json:"ray"`
}

type AggregatedData struct {
	JobStatus models.JobState `gorm:"column:job_status"`
	Count     int
	JobType   models.JobType `gorm:"column:job_type"`
}

func GetJobsQueueSummaryHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var summary Summary
		var aggregatedResults []AggregatedData

		// Perform the query using GORM
		db = db.Debug()
		result := db.Table("jobs").
			Select("jobs.job_type, jobs.job_status, count(*) as count").
			Joins("left join models on models.id = jobs.model_id").
			Group("jobs.job_type, jobs.job_status").
			Find(&aggregatedResults)
		fmt.Printf("Aggregated Results: %+v\n", aggregatedResults)

		if result.Error != nil {
			http.Error(w, fmt.Sprintf("Error Querying Job Table (%v)", result.Error), http.StatusInternalServerError)
			return
		}

		// Compile results into summary
		for _, data := range aggregatedResults {
			jobSummary := JobSummary{
				Count: data.Count,
			}

			if data.JobStatus == models.JobStatePending {
				summary.Ray.Queued = jobSummary
			} else if data.JobStatus == models.JobStateRunning {
				summary.Ray.Running = jobSummary
			}

		}

		// Set content type and encode summary to JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(summary)
	}
}
