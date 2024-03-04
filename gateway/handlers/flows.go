package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
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

		user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
		if !ok {
			utils.SendJSONError(w, "User not found in context", http.StatusUnauthorized)
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

		flow := models.Flow{
			CID:           ioListCid,
			WalletAddress: user.WalletAddress,
			Name:          name,
			StartTime:     time.Now(),
		}

		log.Println("Creating Flow entry")
		result = db.Create(&flow)
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
			jobUUID := uuid.New().String()

			job := models.Job{
				ToolID:        ioItem.Tool.IPFS,
				FlowID:        flow.ID,
				WalletAddress: user.WalletAddress,
				Inputs:        datatypes.JSON(inputsJSON),
				Queue:         queue,
				CreatedAt:     time.Now(),
				JobUUID:       jobUUID,
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
					// result := db.First(&dataFile, "cid = ? and wallet_address = ? ", cid, user.WalletAddress)
					if result.Error != nil {
						if errors.Is(result.Error, gorm.ErrRecordNotFound) {
							// http.Error(w, fmt.Sprintf("DataFile with CID %v and WalletAddress %v not found", cid, user.WalletAddress), http.StatusInternalServerError)
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
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(flow); err != nil {
			http.Error(w, "Error encoding Flow to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func GetFlowHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
		if !ok {
			utils.SendJSONError(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		params := mux.Vars(r)
		flowID, err := strconv.Atoi(params["flowID"])
		if err != nil {
			http.Error(w, fmt.Sprintf("Flow ID (%v) could not be converted to int", params["flowID"]), http.StatusNotFound)
			return
		}

		var flow models.Flow
		if result := db.Preload("Jobs.Tool").First(&flow, "id = ? AND wallet_address = ?", flowID, user.WalletAddress); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Flow not found or not authorized", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Flow: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

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

		user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
		if !ok {
			utils.SendJSONError(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		query := db.Model(&models.Flow{}).Where("wallet_address = ?", user.WalletAddress)

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
