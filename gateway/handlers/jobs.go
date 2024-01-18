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
		jobID, err := strconv.Atoi(params["jobID"])
		if err != nil {
			http.Error(w, fmt.Sprintf("Job ID (%v) could not be converted to int", params["jobID"]), http.StatusNotFound)
		}

		var job models.Job
		if result := db.Preload("OutputFiles.Tags").Preload("InputFiles.Tags").First(&job, "id = ?", jobID); result.Error != nil {
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

func StreamJobLogsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Streaming job logs")
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Check the origin of the request and return true if it's allowed
			// Here's a simple example that allows any origin:
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

	// Channel to gather output
	outputChan := make(chan string)

	// Read from stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			outputChan <- scanner.Text()
		}
	}()

	// Read from stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			outputChan <- scanner.Text()
		}
	}()

	// Write to WebSocket
	go func() {
		for {
			select {
			case outputLine, ok := <-outputChan:
				if !ok {
					// Handle closed channel, if necessary
					return
				}
				// Send to WebSocket
				if err := conn.WriteMessage(websocket.TextMessage, []byte(outputLine)); err != nil {
					log.Println("Error sending message through WebSocket:", err)
					return
				}
			}
		}
	}()

	// If you need to wait for cmd to finish
	cmd.Wait()
}
