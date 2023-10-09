package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/bacalhau-project/bacalhau/pkg/model"
	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/bacalhau"
	"gorm.io/gorm"
)

func GetJobHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		// Get the ID from the URL
		params := mux.Vars(r)
		bacalhauJobID := params["bacalhauJobID"]

		var job models.Job
		if result := db.Preload("Inputs").First(&job, "bacalhau_job_id = ?", bacalhauJobID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Job not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Job: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		log.Println("Fetched Job from DB: ", job)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(job); err != nil {
			http.Error(w, "Error encoding Job to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func UpdateJobHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			utils.SendJSONError(w, "Only PATCH method is supported", http.StatusBadRequest)
			return
		}

		// Get the ID from the URL
		params := mux.Vars(r)
		bacalhauJobID := params["bacalhauJobID"]

		var job models.Job
		if result := db.First(&job, "bacalhau_job_id = ?", bacalhauJobID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Job not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Job: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		log.Println("Updating job")
		updatedJob, err := bacalhau.GetBacalhauJobState(job.BacalhauJobID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating job %v", err), http.StatusInternalServerError)
		}

		// Update job state based on the external job state
		if updatedJob.State.State == model.JobStateCancelled {
			job.State = "failed"
		} else if updatedJob.State.State == model.JobStateError {
			job.State = "error"
		} else if updatedJob.State.State == model.JobStateQueued {
			job.State = "queued"
		} else if updatedJob.State.State == model.JobStateInProgress {
			job.State = "processing"
		} else if updatedJob.State.State == model.JobStateCompleted {
			job.State = "completed"
		}
		log.Println("Updated job")

		// Save the updated job back to the database
		if err := db.Save(job).Error; err != nil {
			http.Error(w, fmt.Sprintf("Error saving Job: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(job); err != nil {
			http.Error(w, "Error encoding Job to JSON", http.StatusInternalServerError)
			return
		}
	}
}
