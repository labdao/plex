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
	"github.com/labdao/plex/internal/ipwl"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func AddExperimentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received Post request at /experiments")
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
			utils.SendJSONError(w, "Invalid or missing Model CID", http.StatusBadRequest)
			return
		}

		var model models.Model
		result := db.Where("cid = ?", toolCid).First(&model)
		if result.Error != nil {
			log.Printf("Error fetching Model: %v\n", result.Error)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				utils.SendJSONError(w, "Model not found", http.StatusNotFound)
			} else {
				utils.SendJSONError(w, "Error fetching Model", http.StatusInternalServerError)
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

		experimentUUID := uuid.New().String()

		experiment := models.Experiment{
			WalletAddress:  user.WalletAddress,
			Name:           name,
			StartTime:      time.Now().UTC(),
			ExperimentUUID: experimentUUID,
			Public:         false,
		}

		log.Println("Creating Experiment entry")
		result = db.Create(&experiment)
		if result.Error != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error creating Experiment entity: %v", result.Error), http.StatusInternalServerError)
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
			if model.ToolType == "ray" {
				queue = models.QueueTypeRay
			}
			// TODO: consolidate below with the above checks.
			var jobType models.JobType
			if model.ToolType == "ray" {
				jobType = models.JobTypeRay
			} else {
				jobType = models.JobTypeBacalhau
			}
			job := models.Job{
				ToolID:        ioItem.Model.S3,
				ExperimentID:  experiment.ID,
				WalletAddress: user.WalletAddress,
				Inputs:        datatypes.JSON(inputsJSON),
				Queue:         queue,
				CreatedAt:     time.Now().UTC(),
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
			requestTracker := models.RequestTracker{
				JobID:      job.ID,
				RetryCount: 0,
				CreatedAt:  time.Now().UTC(),
			}

			result = db.Save(&requestTracker)
			if result.Error != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error creating RequestTracker entity: %v", result.Error), http.StatusInternalServerError)
				return
			}

			user.ComputeTally += model.ComputeCost
			result = db.Save(user)
			if result.Error != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error updating user compute tally: %v", result.Error), http.StatusInternalServerError)
				return
			}

			thresholdStr := os.Getenv("TIER_THRESHOLD")
			if thresholdStr == "" {
				utils.SendJSONError(w, "TIER_THRESHOLD environment variable is not set", http.StatusInternalServerError)
				return
			}

			threshold, err := strconv.Atoi(thresholdStr)
			if err != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error converting TIER_THRESHOLD to integer: %v", err), http.StatusInternalServerError)
				return
			}

			err = UpdateUserTier(db, user.WalletAddress, threshold)
			if err != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error updating user tier: %v", err), http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(experiment); err != nil {
			utils.SendJSONError(w, "Error encoding Experiment to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func GetExperimentHandler(db *gorm.DB) http.HandlerFunc {
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
		experimentID, err := strconv.Atoi(params["experimentID"])
		if err != nil {
			http.Error(w, fmt.Sprintf("Experiment ID (%v) could not be converted to int", params["experimentID"]), http.StatusNotFound)
			return
		}

		var experiment models.Experiment
		query := db.Preload("Jobs.Model").Where("id = ?", experimentID)

		if result := query.First(&experiment); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Experiment not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Experiment: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		if !experiment.Public && experiment.WalletAddress != user.WalletAddress && !user.Admin {
			http.Error(w, "Experiment not found or not authorized", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(experiment); err != nil {
			http.Error(w, "Error encoding Experiment to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func ListExperimentsHandler(db *gorm.DB) http.HandlerFunc {
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

		query := db.Model(&models.Experiment{}).Where("wallet_address = ?", user.WalletAddress)

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

		var experiments []models.Experiment
		if result := query.Preload("Jobs").Find(&experiments); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching Experiments: %v", result.Error), http.StatusInternalServerError)
			return
		}

		log.Println("Fetched experiments from DB: ", experiments)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(experiments); err != nil {
			http.Error(w, "Error encoding Experiments to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func UpdateExperimentHandler(db *gorm.DB) http.HandlerFunc {
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
		experimentID, err := strconv.Atoi(params["experimentID"])
		if err != nil {
			http.Error(w, fmt.Sprintf("Experiment ID (%v) could not be converted to int", params["experimentID"]), http.StatusNotFound)
			return
		}

		var experiment models.Experiment
		if result := db.Where("id = ?", experimentID).First(&experiment); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Experiment not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Experiment: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		if experiment.WalletAddress != user.WalletAddress {
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
		if requestData.Public != nil && *requestData.Public != experiment.Public {
			if experiment.Public {
				http.Error(w, "Experiment is already public and cannot be made private", http.StatusBadRequest)
				return
			}
			newPublicFlag = *requestData.Public
		}
		if newName != "" {
			experiment.Name = newName
			if result := db.Model(&experiment).Updates(models.Experiment{Name: experiment.Name}); result.Error != nil {
				http.Error(w, fmt.Sprintf("Error updating Experiment: %v", result.Error), http.StatusInternalServerError)
				return
			}
		}
		if newPublicFlag {
			experiment.Public = true

			if result := db.Model(&experiment).Updates(models.Experiment{Public: experiment.Public}); result.Error != nil {
				http.Error(w, fmt.Sprintf("Error updating Experiment: %v", result.Error), http.StatusInternalServerError)
				return
			}

			var jobs []models.Job
			if result := db.Where("experiment_id = ?", experiment.ID).Find(&jobs); result.Error != nil {
				http.Error(w, fmt.Sprintf("Error fetching Jobs: %v", result.Error), http.StatusInternalServerError)
				return
			}

			for _, job := range jobs {
				if result := db.Model(&job).Updates(models.Job{Public: experiment.Public}); result.Error != nil {
					http.Error(w, fmt.Sprintf("Error updating Job: %v", result.Error), http.StatusInternalServerError)
					return
				}

				if result := db.Model(&models.DataFile{}).Where("cid IN (?)", db.Table("job_input_files").Select("data_file_c_id").Where("job_id = ?", job.ID)).Updates(models.DataFile{Public: experiment.Public}); result.Error != nil {
					http.Error(w, fmt.Sprintf("Error updating input DataFiles: %v", result.Error), http.StatusInternalServerError)
					return
				}

				if result := db.Model(&models.DataFile{}).Where("cid IN (?)", db.Table("job_output_files").Select("data_file_c_id").Where("job_id = ?", job.ID)).Updates(models.DataFile{Public: experiment.Public}); result.Error != nil {
					http.Error(w, fmt.Sprintf("Error updating output DataFiles: %v", result.Error), http.StatusInternalServerError)
					return
				}
			}

			log.Println("Generating and storing RecordCID...")
			metadataCID, err := utils.GenerateAndStoreRecordCID(db, &experiment)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error generating and storing RecordCID: %v", err), http.StatusInternalServerError)
				return
			}
			log.Printf("Generated and stored RecordCID: %s", metadataCID)

			log.Println("Minting NFT...")
			if err := utils.MintNFT(db, &experiment, metadataCID); err != nil {
				http.Error(w, fmt.Sprintf("Error minting NFT: %v", err), http.StatusInternalServerError)
				return
			}
			log.Println("NFT minted")
		}

		log.Printf("Updated Experiment: %+v", experiment)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(experiment); err != nil {
			http.Error(w, "Error encoding Experiment to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func AddJobToExperimentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received Post request to add job to a experiment")
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
		experimentID, err := strconv.Atoi(params["experimentID"])
		if err != nil {
			http.Error(w, fmt.Sprintf("Experiment ID (%v) could not be converted to int", params["experimentID"]), http.StatusNotFound)
			return
		}

		var experiment models.Experiment
		if result := db.Preload("Jobs").Where("id = ?", experimentID).First(&experiment); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Experiment not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error fetching Experiment: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		if experiment.WalletAddress != user.WalletAddress {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		//TODO: think about moving toolID to experiment level instead of job level
		var toolId = experiment.Jobs[0].ToolID

		var model models.Model
		result := db.Where("cid = ?", toolId).First(&model)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Model not found", http.StatusNotFound)
			} else {
				http.Error(w, "Error fetching Model", http.StatusInternalServerError)
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

		ioList, err := ipwl.InitializeIo(model.CID, scatteringMethod, kwargs, db)
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
			if model.ToolType == "ray" {
				queue = models.QueueTypeRay
			}

			job := models.Job{
				ToolID:        ioItem.Model.S3,
				ExperimentID:  experiment.ID,
				WalletAddress: user.WalletAddress,
				Inputs:        datatypes.JSON(inputsJSON),
				Queue:         queue,
				CreatedAt:     time.Now().UTC(),
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

			requestTracker := models.RequestTracker{
				JobID:      job.ID,
				RetryCount: 0,
				State:      models.JobStateQueued,
				CreatedAt:  time.Now().UTC(),
			}

			result = db.Save(&requestTracker)
			if result.Error != nil {
				http.Error(w, fmt.Sprintf("Error creating RequestTracker entity: %v", result.Error), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(experiment); err != nil {
			http.Error(w, "Error encoding Experiment to JSON", http.StatusInternalServerError)
			return
		}
	}
}
