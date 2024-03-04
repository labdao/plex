package handlers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/labdao/plex/gateway/middleware"
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
		if result := db.Preload("OutputFiles.Tags").Preload("InputFiles.Tags").Where("wallet_address = ?", user.WalletAddress).First(&job, "id = ?", jobID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Job not found or not authorized", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Job: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(job); err != nil {
			http.Error(w, "Error encoding Job to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func StreamJobLogsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Streaming job logs")
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	params := mux.Vars(r)
	bacalhauJobID := params["bacalhauJobID"]
	bacalhauApiHost := bacalhau.GetBacalhauApiHost()
	cmd := exec.Command("bacalhau", "logs", "-f", "--api-host", bacalhauApiHost, bacalhauJobID)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Error creating stdout pipe:", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println("Error creating stderr pipe:", err)
		return
	}

	fmt.Println("Starting command", cmd)
	if err := cmd.Start(); err != nil {
		log.Println("Error starting command:", err)
		return
	}

	outputChan := make(chan string)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			outputChan <- scanner.Text()
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			outputChan <- scanner.Text()
		}
	}()

	go func() {
		for {
			select {
			case outputLine, ok := <-outputChan:
				if !ok {
					return
				}
				if err := conn.WriteMessage(websocket.TextMessage, []byte(outputLine)); err != nil {
					log.Println("Error sending message through WebSocket:", err)
					return
				}
			}
		}
	}()

	cmd.Wait()
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
	CPU QueueSummary `json:"cpu"`
	GPU QueueSummary `json:"gpu"`
}

type AggregatedData struct {
	QueueType     models.QueueType `gorm:"column:queue"`
	State         models.JobState  `gorm:"column:state"`
	TotalCpu      float64
	TotalMemoryGb int
	TotalGpu      int
	Count         int
}

func GetJobsQueueSummaryHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var summary Summary
		var aggregatedResults []AggregatedData

		// Perform the query using GORM
		db = db.Debug()
		result := db.Table("jobs").
			Select("queue, state, sum(tools.cpu) as total_cpu, sum(tools.memory) as total_memory_gb, sum(tools.gpu) as total_gpu, count(*) as count").
			Joins("left join tools on tools.cid = jobs.tool_id").
			Group("queue, state").
			Find(&aggregatedResults)
		fmt.Printf("Aggregated Results: %+v\n", aggregatedResults)

		if result.Error != nil {
			http.Error(w, fmt.Sprintf("Error Querying Job Table (%v)", result.Error), http.StatusInternalServerError)
			return
		}

		// Compile results into summary
		for _, data := range aggregatedResults {
			jobSummary := JobSummary{
				Count:         data.Count,
				TotalCpu:      data.TotalCpu,
				TotalMemoryGb: data.TotalMemoryGb,
				TotalGpu:      data.TotalGpu,
			}

			switch data.QueueType {
			case models.QueueTypeCPU:
				if data.State == models.JobStateQueued {
					summary.CPU.Queued = jobSummary
				} else if data.State == models.JobStateRunning {
					summary.CPU.Running = jobSummary
				}
			case models.QueueTypeGPU:
				if data.State == models.JobStateQueued {
					summary.GPU.Queued = jobSummary
				} else if data.State == models.JobStateRunning {
					summary.GPU.Running = jobSummary
				}
			}
		}

		// Set content type and encode summary to JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(summary)
	}
}
