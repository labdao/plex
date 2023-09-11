package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/ipfs"

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

		log.Println("Request body: ", string(body))

		var requestData map[string]interface{}
		err = json.Unmarshal(body, &requestData)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		tool, ok := requestData["toolData"].(map[string]interface{})
		if !ok {
			http.Error(w, "Invalid or missing tool", http.StatusBadRequest)
			return
		}

		walletAddress, ok := requestData["walletAddress"].(string)
		if !ok {
			http.Error(w, "Invalid or missing walletAddress", http.StatusBadRequest)
			return
		}

		// TODO: Validate Tool

		toolJSON, err := json.Marshal(tool)
		if err != nil {
			http.Error(w, "Error serializing tool to JSON", http.StatusInternalServerError)
			return
		}

		toolName, ok := tool["name"].(string)
		if !ok {
			http.Error(w, "Invalid or missing tool name", http.StatusBadRequest)
			return
		}

		reader := bytes.NewReader(toolJSON)
		tempFile, err := utils.CreateAndWriteTempFile(reader, toolName+".json")
		if err != nil {
			http.Error(w, "Error creating temp file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())

		cid, err := ipfs.WrapAndPinFile(tempFile.Name())
		if err != nil {
			http.Error(w, "Error adding to IPFS", http.StatusInternalServerError)
			return
		}

		// Store serialized Tool in DB
		toolEntity := models.ToolEntity{
			CID:           cid,
			ToolJSON:      string(toolJSON),
			WalletAddress: walletAddress,
		}

		result := db.Create(&toolEntity)
		if result.Error != nil {
			if utils.IsDuplicateKeyError(result.Error) {
				http.Error(w, "A tool with the same CID already exists", http.StatusConflict)
			} else {
				http.Error(w, fmt.Sprintf("Error creating tool entity: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		utils.SendJSONResponseWithCID(w, toolEntity.CID)
	}
}

func GetToolHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		// Get the ID from the URL
		params := mux.Vars(r)
		cid := params["cid"]

		var tool models.ToolEntity
		if result := db.First(&tool, "cid = ?", cid); result.Error != nil {
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

func GetToolsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		var tools []models.ToolEntity
		if result := db.Find(&tools); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching tools: %v", result.Error), http.StatusInternalServerError)
			return
		}

		log.Println("Fetching tools from DB: ", tools)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tools); err != nil {
			http.Error(w, "Error encoding tools to JSON", http.StatusInternalServerError)
			return
		}
	}
}
