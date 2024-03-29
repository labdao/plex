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
	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"

	"gorm.io/gorm"
)

func AddToolHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request at /add-tool")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Println("Request body: ", string(body))

		var toolRequest struct {
			ToolJson json.RawMessage `json:"toolJson"`
		}
		err = json.Unmarshal(body, &toolRequest)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
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

		// Unmarshal the tool from the tool JSON
		var tool ipwl.Tool
		err = json.Unmarshal(toolRequest.ToolJson, &tool)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid toolJson format: %v", err), http.StatusBadRequest)
			return
		}

		toolJSON, err := json.Marshal(tool)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error re-marshalling tool data: %v", err), http.StatusInternalServerError)
			return
		}

		reader := bytes.NewReader(toolJSON)
		tempFile, err := utils.CreateAndWriteTempFile(reader, tool.Name+".json")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating temp file: %v", err), http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())

		cid, err := ipfs.WrapAndPinFile(tempFile.Name())
		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding to IPFS: %v", err), http.StatusInternalServerError)
			return
		}

		var toolGpu int
		if tool.GpuBool {
			toolGpu = 1
		} else {
			toolGpu = 0
		}

		var display bool = true
		var defaultTool bool = false

		var taskCategory string
		if tool.TaskCategory == "" {
			taskCategory = "community-models"
		} else {
			taskCategory = tool.TaskCategory
		}

		// Start transaction
		tx := db.Begin()

		// If the new tool is marked as default, reset the default tool in the same task category
		if defaultTool {
			if err := tx.Model(&models.Tool{}).
				Where("task_category = ? AND default_tool = TRUE", tool.TaskCategory).
				Update("default_tool", false).Error; err != nil {
				tx.Rollback()
				http.Error(w, fmt.Sprintf("Error resetting existing default tool: %v", err), http.StatusInternalServerError)
				return
			}
		}

		toolEntry := models.Tool{
			CID:           cid,
			WalletAddress: walletAddress,
			Name:          tool.Name,
			ToolJson:      toolJSON,
			Container:     tool.DockerPull,
			Memory:        *tool.MemoryGB,
			Cpu:           *tool.Cpu,
			Gpu:           toolGpu,
			Network:       tool.NetworkBool,
			Timestamp:     time.Now(),
			Display:       display,
			TaskCategory:  taskCategory,
			DefaultTool:   defaultTool,
		}

		result := tx.Create(&toolEntry)
		if result.Error != nil {
			tx.Rollback()
			if utils.IsDuplicateKeyError(result.Error) {
				http.Error(w, "A tool with the same CID already exists", http.StatusConflict)
			} else {
				http.Error(w, fmt.Sprintf("Error creating tool entity: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			http.Error(w, fmt.Sprintf("Transaction commit error: %v", err), http.StatusInternalServerError)
			return
		}

		utils.SendJSONResponseWithCID(w, toolEntry.CID)
	}
}

func UpdateToolHandler(db *gorm.DB) http.HandlerFunc {
	acceptedTaskCategories := map[string]bool{
		"protein-binder-design": true,
		"protein-folding":       true,
		"community-models":      true,
		// To do: add tool task category should also only accept these accepted categories
		// To do: remove hardcoding later to match one of the available slugs from tasks taskList.ts with available: true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.SendJSONError(w, "Only PUT method is supported", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		cid := vars["cid"]
		if cid == "" {
			utils.SendJSONError(w, "Missing CID parameter", http.StatusBadRequest)
			return
		}

		var requestData struct {
			TaskCategory *string `json:"taskCategory,omitempty"`
			Display      *bool   `json:"display,omitempty"`
			DefaultTool  *bool   `json:"defaultTool,omitempty"`
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
		if requestData.DefaultTool != nil {
			updateData["default_tool"] = *requestData.DefaultTool
		}

		if len(updateData) == 0 {
			utils.SendJSONError(w, "No valid fields provided for update", http.StatusBadRequest)
			return
		}

		result := tx.Model(&models.Tool{}).Where("cid = ?", cid).Updates(updateData)

		if result.Error != nil {
			tx.Rollback()
			http.Error(w, fmt.Sprintf("Error updating tool: %v", result.Error), http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			utils.SendJSONError(w, "Tool with the specified CID not found", http.StatusNotFound)
			return
		}

		if err := tx.Commit().Error; err != nil {
			http.Error(w, fmt.Sprintf("Transaction commit error: %v", err), http.StatusInternalServerError)
			return
		}

		utils.SendJSONResponse(w, "Tool updated successfully")
	}
}

func GetToolHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		cid := vars["cid"]
		if cid == "" {
			utils.SendJSONError(w, "Missing CID parameter", http.StatusBadRequest)
			return
		}

		var tool models.Tool
		result := db.Where("cid = ?", cid).First(&tool)
		if result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching tool: %v", result.Error), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tool); err != nil {
			http.Error(w, "Error encoding tool to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func ListToolsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		query := db.Model(&models.Tool{})

		// If display is provided, filter based on it, if not, display only tools where 'display' is true by default
		displayParam, displayProvided := r.URL.Query()["display"]
		if displayProvided && len(displayParam[0]) > 0 {
			query = query.Where("display = ?", displayParam[0] == "true")
		} else {
			query = query.Where("display = ?", true)
		}

		if cid := r.URL.Query().Get("cid"); cid != "" {
			query = query.Where("cid = ?", cid)
		}

		if name := r.URL.Query().Get("name"); name != "" {
			query = query.Where("name = ?", name)
		}

		if walletAddress := r.URL.Query().Get("walletAddress"); walletAddress != "" {
			query = query.Where("wallet_address = ?", walletAddress)
		}

		if taskCategory := r.URL.Query().Get("taskCategory"); taskCategory != "" {
			query = query.Where("task_category = ?", taskCategory)
		}

		var tools []models.Tool
		if result := query.Find(&tools); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching tools: %v", result.Error), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tools); err != nil {
			http.Error(w, "Error encoding tools to JSON", http.StatusInternalServerError)
			return
		}
	}
}
