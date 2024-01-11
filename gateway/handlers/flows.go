package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bacalhau-project/bacalhau/pkg/model"
	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/bacalhau"
	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func pinIoList(ios []ipwl.IO) (string, error) {
	data, err := json.Marshal(ios)
	if err != nil {
		return "", fmt.Errorf("failed to marshal IO slice: %v", err)
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	if err != nil {
		return "", fmt.Errorf("cannot create temporary file: %v", err)
	}

	if _, err = tmpFile.Write(data); err != nil {
		return "", fmt.Errorf("failed to write to temporary file: %v", err)
	}

	cid, err := ipfs.PinFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to pin file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close the file: %v", err)
	}

	return cid, nil
}

func AddFlowHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO McMenemy, update so that jobs are just created to DB and not Bacalhau
		log.Println("Received Post request at /flows")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		log.Println("Request body: ", string(body))

		requestData := make(map[string]json.RawMessage)
		err = json.Unmarshal(body, &requestData)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var walletAddress string
		err = json.Unmarshal(requestData["walletAddress"], &walletAddress)
		if err != nil || walletAddress == "" {
			http.Error(w, "Invalid or missing walletAddress", http.StatusBadRequest)
			return
		}

		var toolCid string
		err = json.Unmarshal(requestData["toolCid"], &toolCid)
		if err != nil || toolCid == "" {
			http.Error(w, "Invalid or missing Tool CID", http.StatusBadRequest)
			return
		}

		var tool models.Tool
		result := db.Where("cid = ?", toolCid).First(&tool)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Tool not found", http.StatusNotFound)
			} else {
				http.Error(w, "Error fetching Tool", http.StatusInternalServerError)
			}
			return
		}

		var scatteringMethod string
		err = json.Unmarshal(requestData["scatteringMethod"], &scatteringMethod)
		if err != nil || scatteringMethod == "" {
			http.Error(w, "Invalid or missing Scattering Method", http.StatusBadRequest)
			return
		}

		var name string
		err = json.Unmarshal(requestData["name"], &name)
		if err != nil || name == "" {
			http.Error(w, "Invalid or missing Name", http.StatusBadRequest)
			return
		}

		kwargsRaw, ok := requestData["kwargs"]
		if !ok {
			http.Error(w, "missing kwargs in the request", http.StatusBadRequest)
			return
		}

		var kwargs map[string][]interface{}
		err = json.Unmarshal(kwargsRaw, &kwargs)
		if err != nil {
			log.Printf("Error unmarshalling kwargs: %v; Raw data: %s\n", err, string(kwargsRaw))
			http.Error(w, "Invalid structure for kwargs", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(requestData["kwargs"], &kwargs)
		if err != nil {
			http.Error(w, "Invalid or missing kwargs", http.StatusBadRequest)
			return
		}

		ioList, err := ipwl.InitializeIo(toolCid, scatteringMethod, kwargs)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while transforming validated JSON: %v", err), http.StatusInternalServerError)
			return
		}
		log.Println("Initialized IO List")

		// save ioList to IPFS
		ioListCid, err := pinIoList(ioList)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error pinning IO List: %v", err), http.StatusInternalServerError)
			return
		}

		flowEntry := models.Flow{
			CID:           ioListCid,
			WalletAddress: walletAddress,
			Name:          name,
		}

		log.Println("Creating Flow entry")
		result = db.Create(&flowEntry)
		if result.Error != nil {
			http.Error(w, fmt.Sprintf("Error creating Flow entity: %v", result.Error), http.StatusInternalServerError)
			return
		}

		for _, ioItem := range ioList {
			log.Println("Creating job entry")
			inputsJSON, err := json.Marshal(ioItem.Inputs)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error transforming job inputs: %v", err), http.StatusInternalServerError)
				return
			}
			var queue models.QueueType
			if tool.Gpu == 0 {
				queue = models.QueueTypeCPU
			} else {
				queue = models.QueueTypeGPU
			}
			job := models.Job{
				ToolID:        ioItem.Tool.IPFS,
				FlowID:        flowEntry.CID,
				WalletAddress: walletAddress,
				Inputs:        datatypes.JSON(inputsJSON),
				Queue:         queue,
			}
			result := db.Create(&job)
			if result.Error != nil {
				http.Error(w, fmt.Sprintf("Error creating Job entity: %v", result.Error), http.StatusInternalServerError)
				return
			}

			for _, input := range ioItem.Inputs {
				var cidsToAdd []string
				switch v := input.(type) {
				case string:
					strInput, ok := input.(string)
					if !ok {
						continue
					}
					if strings.HasPrefix(strInput, "Qm") && strings.Contains(strInput, "/") {
						split := strings.SplitN(strInput, "/", 2)
						cid := split[0]
						cidsToAdd = append(cidsToAdd, cid)
					}
				case []interface{}:
					fmt.Println("found slice, checking each for 'Qm' prefix")
					for _, elem := range v {
						strInput, ok := elem.(string)
						if !ok {
							continue
						}
						if strings.HasPrefix(strInput, "Qm") && strings.Contains(strInput, "/") {
							split := strings.SplitN(strInput, "/", 2)
							cid := split[0]
							cidsToAdd = append(cidsToAdd, cid)
						}
					}
				default:
					continue
				}
				for _, cid := range cidsToAdd {
					var dataFile models.DataFile
					result := db.First(&dataFile, "cid = ?", cid)
					if result.Error != nil {
						if errors.Is(result.Error, gorm.ErrRecordNotFound) {
							http.Error(w, fmt.Sprintf("DataFile with CID %v not found", cid), http.StatusInternalServerError)
							return
						} else {
							http.Error(w, fmt.Sprintf("Error looking up DataFile: %v", result.Error), http.StatusInternalServerError)
							return
						}
					}
					job.InputFiles = append(job.InputFiles, dataFile)
				}
			}
			result = db.Save(&job)
			if result.Error != nil {
				http.Error(w, fmt.Sprintf("Error updating Job entity with input data: %v", result.Error), http.StatusInternalServerError)
				return
			}
		}
		utils.SendJSONResponseWithCID(w, ioListCid)
	}
}

func GetFlowHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		params := mux.Vars(r)
		cid := params["cid"]

		var flow models.Flow
		if result := db.Preload("Jobs.Tool").First(&flow, "cid = ?", cid); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Flow not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Flow: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		log.Println("Fetched flow from DB: ", flow)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(flow); err != nil {
			http.Error(w, "Error encoding Flow to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func UpdateFlowHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received Patch request at /flows")
		if r.Method != http.MethodPatch {
			utils.SendJSONError(w, "Only PATCH method is supported", http.StatusBadRequest)
			return
		}
		log.Println("Received Patch request at /flows")

		params := mux.Vars(r)
		cid := params["cid"]
		log.Println("CID: ", cid)

		var flow models.Flow
		if result := db.Preload("Jobs.Tool").First(&flow, "cid = ?", cid); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Flow not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Flow: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		log.Println("Fetched flow from DB")
		for index, job := range flow.Jobs {
			log.Println("Updating job: ", index)
			updatedJob, err := bacalhau.GetBacalhauJobState(job.BacalhauJobID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error updating job %v", err), http.StatusInternalServerError)
			}
			if updatedJob.State.State == model.JobStateCancelled {
				flow.Jobs[index].State = "failed"
			} else if updatedJob.State.State == model.JobStateError {
				flow.Jobs[index].State = "error"
			} else if updatedJob.State.State == model.JobStateQueued {
				flow.Jobs[index].State = "queued"
			} else if updatedJob.State.State == model.JobStateInProgress {
				flow.Jobs[index].State = "processing"
			} else if updatedJob.State.State == model.JobStateCompleted {
				flow.Jobs[index].State = "completed"
			} else if len(updatedJob.State.Executions) > 0 && updatedJob.State.Executions[0].State == model.ExecutionStateFailed {
				flow.Jobs[index].State = "failed"
			}

			log.Println("Updated job")
			if err := db.Save(&flow.Jobs[index]).Error; err != nil {
				http.Error(w, fmt.Sprintf("Error saving job: %v", err), http.StatusInternalServerError)
				return
			}
		}

		log.Println("Updated flow from DB: ", flow)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(flow); err != nil {
			http.Error(w, "Error encoding Flow to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func ListFlowsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		query := db.Model(&models.Flow{})

		if cid := r.URL.Query().Get("cid"); cid != "" {
			query = query.Where("cid = ?", cid)
		}

		if name := r.URL.Query().Get("name"); name != "" {
			query = query.Where("name = ?", name)
		}

		if walletAddress := r.URL.Query().Get("walletAddress"); walletAddress != "" {
			query = query.Where("wallet_address = ?", walletAddress)
		}

		var flows []models.Flow
		if result := query.Preload("Jobs").Find(&flows); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching Flows: %v", result.Error), http.StatusInternalServerError)
			return
		}

		log.Println("Fetched flows from DB: ", flows)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(flows); err != nil {
			http.Error(w, "Error encoding Flows to JSON", http.StatusInternalServerError)
			return
		}
	}
}
