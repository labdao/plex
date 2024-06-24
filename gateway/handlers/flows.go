package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/ipwl"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func AddFlowHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received Post request at /flows")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.SendJSONError(w, "Bad request", http.StatusBadRequest)
			return
		}

		log.Println("Request body: ", string(body))

		requestData := make(map[string]json.RawMessage)
		err = json.Unmarshal(body, &requestData)
		if err != nil {
			utils.SendJSONError(w, "Invalid JSON", http.StatusBadRequest)
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
			utils.SendJSONError(w, "Invalid or missing Tool CID", http.StatusBadRequest)
			return
		}

		var tool models.Tool
		result := db.Where("cid = ?", toolCid).First(&tool)
		if result.Error != nil {
			log.Printf("Error fetching Tool: %v\n", result.Error)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				utils.SendJSONError(w, "Tool not found", http.StatusNotFound)
			} else {
				utils.SendJSONError(w, "Error fetching Tool", http.StatusInternalServerError)
			}
			return
		}

		var scatteringMethod string
		err = json.Unmarshal(requestData["scatteringMethod"], &scatteringMethod)
		if err != nil || scatteringMethod == "" {
			utils.SendJSONError(w, "Invalid or missing Scattering Method", http.StatusBadRequest)
			return
		}

		var name string
		err = json.Unmarshal(requestData["name"], &name)
		if err != nil || name == "" {
			utils.SendJSONError(w, "Invalid or missing Name", http.StatusBadRequest)
			return
		}

		kwargsRaw, ok := requestData["kwargs"]
		if !ok {
			utils.SendJSONError(w, "missing kwargs in the request", http.StatusBadRequest)
			return
		}

		var kwargs map[string][]interface{}
		err = json.Unmarshal(kwargsRaw, &kwargs)
		if err != nil {
			log.Printf("Error unmarshalling kwargs: %v; Raw data: %s\n", err, string(kwargsRaw))
			utils.SendJSONError(w, "Invalid structure for kwargs", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(requestData["kwargs"], &kwargs)
		if err != nil {
			utils.SendJSONError(w, "Invalid or missing kwargs", http.StatusBadRequest)
			return
		}

		ioList, err := ipwl.InitializeIo(toolCid, scatteringMethod, kwargs, db)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error while transforming validated JSON: %v", err), http.StatusInternalServerError)
			return
		}
		log.Println("Initialized IO List")

		// save ioList to IPFS
		// ioListCid, err := pinIoList(ioList)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error pinning IO List: %v", err), http.StatusInternalServerError)
			return
		}

		flowUUID := uuid.New().String()

		flow := models.Flow{
			WalletAddress: user.WalletAddress,
			Name:          name,
			StartTime:     time.Now(),
			FlowUUID:      flowUUID,
			Public:        false,
		}

		log.Println("Creating Flow entry")
		result = db.Create(&flow)
		if result.Error != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error creating Flow entity: %v", result.Error), http.StatusInternalServerError)
			return
		}

		for _, ioItem := range ioList {
			log.Println("Creating job entry")
			inputsJSON, err := json.Marshal(ioItem.Inputs)
			if err != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error transforming job inputs: %v", err), http.StatusInternalServerError)
				return
			}
			var queue models.QueueType
			if tool.ToolType == "ray" {
				queue = models.QueueTypeRay
			}
			// TODO: consolidate below with the above checks.
			var jobType models.JobType
			if tool.ToolType == "ray" {
				jobType = models.JobTypeRay
			} else {
				jobType = models.JobTypeBacalhau
			}
			job := models.Job{
				ToolID:        ioItem.Tool.S3,
				FlowID:        flow.ID,
				WalletAddress: user.WalletAddress,
				Inputs:        datatypes.JSON(inputsJSON),
				Queue:         queue,
				CreatedAt:     time.Now(),
				Public:        false,
				JobType:       jobType,
			}

			result := db.Create(&job)
			if result.Error != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error creating Job entity: %v", result.Error), http.StatusInternalServerError)
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
							utils.SendJSONError(w, fmt.Sprintf("DataFile with CID %v not found", cid), http.StatusInternalServerError)
							return
						} else {
							utils.SendJSONError(w, fmt.Sprintf("Error looking up DataFile: %v", result.Error), http.StatusInternalServerError)
							return
						}
					}
					job.InputFiles = append(job.InputFiles, dataFile)
				}
			}
			result = db.Save(&job)
			if result.Error != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error updating Job entity with input data: %v", result.Error), http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(flow); err != nil {
			utils.SendJSONError(w, "Error encoding Flow to JSON", http.StatusInternalServerError)
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
		query := db.Preload("Jobs.Tool").Where("id = ?", flowID)

		if result := query.First(&flow); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Flow not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Flow: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		if !flow.Public && flow.WalletAddress != user.WalletAddress && !user.Admin {
			http.Error(w, "Flow not found or not authorized", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(flow); err != nil {
			http.Error(w, "Error encoding Flow to JSON", http.StatusInternalServerError)
			return
		}
	}
}
func ListFlowNamesHandler(db *gorm.DB) http.HandlerFunc {
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

		var flowNames []struct {
			ID        int
			Name      string
			StartTime time.Time
		}
		if result := db.Model(&models.Flow{}).
			Where("wallet_address = ?", user.WalletAddress).
			Select("name, id, start_time").
			Find(&flowNames); result.Error != nil { // Directly scan the names into the slice of strings
			http.Error(w, fmt.Sprintf("Error fetching flow names: %v", result.Error), http.StatusInternalServerError)
			return
		}

		log.Println("Fetched flow names from DB: ", flowNames)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(flowNames); err != nil { // Encode the slice of flow names directly
			http.Error(w, "Error encoding flow names to JSON", http.StatusInternalServerError)
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
		query = query.Order("start_time DESC")

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

func UpdateFlowHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.SendJSONError(w, "Only PUT method is supported", http.StatusBadRequest)
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
		if result := db.Where("id = ?", flowID).First(&flow); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Flow not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Flow: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		if flow.WalletAddress != user.WalletAddress {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var requestData struct {
			Name   *string `json:"name,omitempty"`
			Public *bool   `json:"public,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		newPublicFlag := false
		newName := ""
		if requestData.Name != nil {
			newName = *requestData.Name
		}
		if requestData.Public != nil && *requestData.Public != flow.Public {
			if flow.Public {
				http.Error(w, "Flow is already public and cannot be made private", http.StatusBadRequest)
				return
			}
			newPublicFlag = *requestData.Public
		}
		if newName != "" {
			flow.Name = newName
			if result := db.Model(&flow).Updates(models.Flow{Name: flow.Name}); result.Error != nil {
				http.Error(w, fmt.Sprintf("Error updating Flow: %v", result.Error), http.StatusInternalServerError)
				return
			}
		}
		if newPublicFlag {
			flow.Public = true

			if result := db.Model(&flow).Updates(models.Flow{Public: flow.Public}); result.Error != nil {
				http.Error(w, fmt.Sprintf("Error updating Flow: %v", result.Error), http.StatusInternalServerError)
				return
			}

			var jobs []models.Job
			if result := db.Where("flow_id = ?", flow.ID).Find(&jobs); result.Error != nil {
				http.Error(w, fmt.Sprintf("Error fetching Jobs: %v", result.Error), http.StatusInternalServerError)
				return
			}

			for _, job := range jobs {
				if result := db.Model(&job).Updates(models.Job{Public: flow.Public}); result.Error != nil {
					http.Error(w, fmt.Sprintf("Error updating Job: %v", result.Error), http.StatusInternalServerError)
					return
				}

				if result := db.Model(&models.DataFile{}).Where("cid IN (?)", db.Table("job_input_files").Select("data_file_c_id").Where("job_id = ?", job.ID)).Updates(models.DataFile{Public: flow.Public}); result.Error != nil {
					http.Error(w, fmt.Sprintf("Error updating input DataFiles: %v", result.Error), http.StatusInternalServerError)
					return
				}

				if result := db.Model(&models.DataFile{}).Where("cid IN (?)", db.Table("job_output_files").Select("data_file_c_id").Where("job_id = ?", job.ID)).Updates(models.DataFile{Public: flow.Public}); result.Error != nil {
					http.Error(w, fmt.Sprintf("Error updating output DataFiles: %v", result.Error), http.StatusInternalServerError)
					return
				}
			}

			log.Println("Generating and storing RecordCID...")
			metadataCID, err := utils.GenerateAndStoreRecordCID(db, &flow)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error generating and storing RecordCID: %v", err), http.StatusInternalServerError)
				return
			}
			log.Printf("Generated and stored RecordCID: %s", metadataCID)

			log.Println("Minting NFT...")
			if err := utils.MintNFT(db, &flow, metadataCID); err != nil {
				http.Error(w, fmt.Sprintf("Error minting NFT: %v", err), http.StatusInternalServerError)
				return
			}
			log.Println("NFT minted")
		}

		log.Printf("Updated Flow: %+v", flow)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(flow); err != nil {
			http.Error(w, "Error encoding Flow to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func AddJobToFlowHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received Post request to add job to a flow")
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
		params := mux.Vars(r)
		flowID, err := strconv.Atoi(params["flowID"])
		if err != nil {
			http.Error(w, fmt.Sprintf("Flow ID (%v) could not be converted to int", params["flowID"]), http.StatusNotFound)
			return
		}

		var flow models.Flow
		if result := db.Preload("Jobs").Where("id = ?", flowID).First(&flow); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Flow not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Flow: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		if flow.WalletAddress != user.WalletAddress {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		//TODO: think about moving toolID to flow level instead of job level
		var toolId = flow.Jobs[0].ToolID

		var tool models.Tool
		result := db.Where("cid = ?", toolId).First(&tool)
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

		ioList, err := ipwl.InitializeIo(tool.CID, scatteringMethod, kwargs, db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while transforming validated JSON: %v", err), http.StatusInternalServerError)
			return
		}
		log.Println("Initialized IO List")

		for _, ioItem := range ioList {
			log.Println("Creating job entry")
			inputsJSON, err := json.Marshal(ioItem.Inputs)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error transforming job inputs: %v", err), http.StatusInternalServerError)
				return
			}
			var queue models.QueueType
			if tool.ToolType == "ray" {
				queue = models.QueueTypeRay
			}

			job := models.Job{
				ToolID:        ioItem.Tool.S3,
				FlowID:        flow.ID,
				WalletAddress: user.WalletAddress,
				Inputs:        datatypes.JSON(inputsJSON),
				Queue:         queue,
				CreatedAt:     time.Now(),
				Public:        false,
			}

			result = db.Create(&job)
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

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(flow); err != nil {
			http.Error(w, "Error encoding Flow to JSON", http.StatusInternalServerError)
			return
		}
	}
}
