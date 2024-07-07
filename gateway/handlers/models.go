package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/labdao/plex/internal/s3"
	"gorm.io/gorm"
)

func AddModelHandler(db *gorm.DB, s3c *s3.S3Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request at /add-model")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// http.Error(w, "Bad request", http.StatusBadRequest)
			utils.SendJSONError(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Println("Request body: ", string(body))

		var modelRequest struct {
			ModelJson json.RawMessage `json:"modelJson"`
		}
		err = json.Unmarshal(body, &modelRequest)
		if err != nil {
			// http.Error(w, "Invalid JSON", http.StatusBadRequest)
			utils.SendJSONError(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		token, err := utils.ExtractAuthHeader(r)
		if err != nil {
			utils.SendJSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var walletAddress string
		if middleware.IsJWT(token) {
			claims, err := middleware.ValidateJWT(token, db)
			if err != nil {
				utils.SendJSONError(w, "Invalid JWT", http.StatusUnauthorized)
				return
			}
			walletAddress, err = middleware.GetWalletAddressFromJWTClaims(claims, db)
			if err != nil {
				utils.SendJSONError(w, err.Error(), http.StatusUnauthorized)
				return
			}
		} else {
			walletAddress, err = middleware.GetWalletAddressFromAPIKey(token, db)
			if err != nil {
				utils.SendJSONError(w, err.Error(), http.StatusUnauthorized)
				return
			}
		}

		// Unmarshal the model from the model JSON
		var model ipwl.Model
		err = json.Unmarshal(modelRequest.ModelJson, &model)
		if err != nil {
			// http.Error(w, fmt.Sprintf("Invalid modelJson format: %v", err), http.StatusBadRequest)
			utils.SendJSONError(w, fmt.Sprintf("Invalid modelJson format: %v", err), http.StatusBadRequest)
			return
		}

		modelJSON, err := json.Marshal(model)
		if err != nil {
			// http.Error(w, fmt.Sprintf("Error re-marshalling model data: %v", err), http.StatusInternalServerError)
			utils.SendJSONError(w, fmt.Sprintf("Error re-marshalling model data: %v", err), http.StatusInternalServerError)
			return
		}

		reader := bytes.NewReader(modelJSON)
		tempFile, err := utils.CreateAndWriteTempFile(reader, model.Name+".json")
		if err != nil {
			// http.Error(w, fmt.Sprintf("Error creating temp file: %v", err), http.StatusInternalServerError)
			utils.SendJSONError(w, fmt.Sprintf("Error creating temp file: %v", err), http.StatusInternalServerError)
			return
		}

		bucketName := os.Getenv("BUCKET_NAME")
		if bucketName == "" {
			utils.SendJSONError(w, "Missing BUCKET_NAME environment variable", http.StatusInternalServerError)
			return
		}

		hash, err := utils.GenerateFileHash(tempFile.Name())
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error hashing file: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Println("Hash of the model manifest", tempFile.Name(), " file: ", hash)
		defer os.Remove(tempFile.Name())

		objectKey := hash + "/" + tempFile.Name()
		// S3 upload
		err = s3c.UploadFile(bucketName, objectKey, tempFile.Name())
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error uploading to bucket: %v", err), http.StatusInternalServerError)
			return
		}
		s3_uri := fmt.Sprintf("s3://%s/%s", bucketName, objectKey)
		var display bool = true
		var defaultModel bool = false

		var taskCategory string
		if model.TaskCategory == "" {
			taskCategory = "community-models"
		} else {
			taskCategory = model.TaskCategory
		}

		var maxRunningTime int
		//if maxruntime is not provided, set it to 45 minutes
		if model.MaxRunningTime == 0 {
			maxRunningTime = 2700
		} else {
			maxRunningTime = model.MaxRunningTime
		}
		// Start transaction
		tx := db.Begin()

		// If the new model is marked as default, reset the default model in the same task category
		if defaultModel {
			if err := tx.Model(&models.Model{}).
				Where("task_category = ? AND default_model = TRUE", model.TaskCategory).
				Update("default_model", false).Error; err != nil {
				tx.Rollback()
				// http.Error(w, fmt.Sprintf("Error resetting existing default model: %v", err), http.StatusInternalServerError)
				utils.SendJSONError(w, fmt.Sprintf("Error resetting existing default model: %v", err), http.StatusInternalServerError)
				return
			}
		}
		var user models.User
		if err := db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
			return
		}

		modelEntry := models.Model{
			WalletAddress:      user.WalletAddress,
			Name:               model.Name,
			ModelJson:          modelJSON,
			CreatedAt:          time.Now().UTC(),
			Display:            display,
			TaskCategory:       taskCategory,
			DefaultModel:       defaultModel,
			MaxRunningTime:     maxRunningTime,
			RayServiceEndpoint: model.RayServiceEndpoint,
			ComputeCost:        model.ComputeCost,
			S3URI:              s3_uri,
		}

		result := tx.Create(&modelEntry)
		if result.Error != nil {
			tx.Rollback()
			if utils.IsDuplicateKeyError(result.Error) {
				// http.Error(w, "A model with the same ID already exists", http.StatusConflict)
				utils.SendJSONError(w, "This model already exists", http.StatusConflict)
			} else {
				// http.Error(w, fmt.Sprintf("Error creating model entity: %v", result.Error), http.StatusInternalServerError)
				utils.SendJSONError(w, fmt.Sprintf("Error creating model entity: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			// http.Error(w, fmt.Sprintf("Transaction commit error: %v", err), http.StatusInternalServerError)
			utils.SendJSONError(w, fmt.Sprintf("Transaction commit error: %v", err), http.StatusInternalServerError)
			return
		}

		utils.SendJSONResponseWithID(w, modelEntry.ID)
	}
}

func UpdateModelHandler(db *gorm.DB) http.HandlerFunc {
	acceptedTaskCategories := map[string]bool{
		"protein-binder-design": true,
		"protein-folding":       true,
		"community-models":      true,
		// To do: add model task category should also only accept these accepted categories
		// To do: remove hardcoding later to match one of the available slugs from tasks taskList.ts with available: true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.SendJSONError(w, "Only PUT method is supported", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]
		if id == "" {
			utils.SendJSONError(w, "Missing ID parameter", http.StatusBadRequest)
			return
		}

		var requestData struct {
			TaskCategory   *string `json:"taskCategory,omitempty"`
			Display        *bool   `json:"display,omitempty"`
			DefaultModel   *bool   `json:"defaultModel,omitempty"`
			MaxRunningTime *int    `json:"maxRunningTime,omitempty"`
			ComputeCost    *int    `json:"computeCost,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		if requestData.TaskCategory != nil {
			if _, ok := acceptedTaskCategories[*requestData.TaskCategory]; !ok {
				utils.SendJSONError(w, "Task category not accepted", http.StatusBadRequest)
				return
			}
		}

		tx := db.Begin()

		updateData := make(map[string]interface{})
		if requestData.TaskCategory != nil {
			updateData["task_category"] = *requestData.TaskCategory
		}
		if requestData.Display != nil {
			updateData["display"] = *requestData.Display
		}
		if requestData.DefaultModel != nil {
			updateData["default_model"] = *requestData.DefaultModel
		}
		if requestData.MaxRunningTime != nil {
			updateData["max_running_time"] = *requestData.MaxRunningTime
		}
		if requestData.ComputeCost != nil {
			updateData["compute_cost"] = *requestData.ComputeCost
		}

		if len(updateData) == 0 {
			utils.SendJSONError(w, "No valid fields provided for update", http.StatusBadRequest)
			return
		}

		result := tx.Model(&models.Model{}).Where("id = ?", id).Updates(updateData)

		if result.Error != nil {
			tx.Rollback()
			http.Error(w, fmt.Sprintf("Error updating model: %v", result.Error), http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			utils.SendJSONError(w, "Model with the specified ID not found", http.StatusNotFound)
			return
		}

		if err := tx.Commit().Error; err != nil {
			http.Error(w, fmt.Sprintf("Transaction commit error: %v", err), http.StatusInternalServerError)
			return
		}

		utils.SendJSONResponse(w, "Model updated successfully")
	}
}

func GetModelHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]
		if id == "" {
			utils.SendJSONError(w, "Missing ID parameter", http.StatusBadRequest)
			return
		}

		var model models.Model
		result := db.Where("id = ?", id).First(&model)
		if result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching model: %v", result.Error), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(model); err != nil {
			http.Error(w, "Error encoding model to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func ListModelsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		query := db.Model(&models.Model{})

		// If display is provided, filter based on it, if not, display only models where 'display' is true by default
		displayParam, displayProvided := r.URL.Query()["display"]
		if displayProvided && len(displayParam[0]) > 0 {
			query = query.Where("display = ?", displayParam[0] == "true")
		} else {
			query = query.Where("display = ?", true)
		}

		if id := r.URL.Query().Get("id"); id != "" {
			query = query.Where("id = ?", id)
		}

		if name := r.URL.Query().Get("name"); name != "" {
			query = query.Where("name = ?", name)
		}

		if walletAddress := r.URL.Query().Get("wallet_address"); walletAddress != "" {
			query = query.Where("wallet_address = ?", walletAddress)
		}

		if taskCategory := r.URL.Query().Get("taskCategory"); taskCategory != "" {
			query = query.Where("task_category = ?", taskCategory)
		}

		var models []models.Model
		if result := query.Find(&models); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching models: %v", result.Error), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(models); err != nil {
			http.Error(w, "Error encoding models to JSON", http.StatusInternalServerError)
			return
		}
	}
}
