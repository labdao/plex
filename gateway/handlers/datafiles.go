package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/ipfs"

	"gorm.io/gorm"
)

func AddDataFileHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request at /create-datafile")

		if err := utils.CheckRequestMethod(r, http.MethodPost); err != nil {
			utils.SendJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			utils.SendJSONError(w, "Error parsing multipart form", http.StatusBadRequest)
			return
		}
		fmt.Println("Parsed multipart form")

		file, _, err := r.FormFile("file")
		if err != nil {
			utils.SendJSONError(w, "Error retrieving file from multipart form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		walletAddress := r.FormValue("walletAddress")
		filename := r.FormValue("filename")
		publicBool := r.FormValue("public")
		visibleBool := r.FormValue("visible")

		fmt.Printf("Received file upload request for file: %s, walletAddress: %s, public: %s, visible: %s\n", filename, walletAddress, publicBool, visibleBool)

		tempFile, err := utils.CreateAndWriteTempFile(file, filename)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error creating temp file: %v", err), http.StatusInternalServerError)
			return
		}
		defer os.Remove(filename)

		cid, err := ipfs.WrapAndPinFile(tempFile.Name())
		if err != nil {
			utils.SendJSONError(w, "Error pinning file to IPFS", http.StatusInternalServerError)
			return
		}

		dataFile := models.DataFile{
			CID:           cid,
			WalletAddress: walletAddress,
			Filename:      filename,
			Timestamp:     time.Now(),
		}

		if result := db.Create(&dataFile); result.Error != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error saving datafile: %v", result.Error), http.StatusInternalServerError)
			return
		}

		utils.SendJSONResponseWithCID(w, dataFile.CID)
	}
}

// get a single datafile
func GetDataFileHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		// Get the ID from the URL
		params := mux.Vars(r)
		id := params["id"]

		var dataFile models.DataFile
		if result := db.First(&dataFile, id); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching datafile: %v", result.Error), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dataFile); err != nil {
			http.Error(w, "Error encoding datafile to JSON", http.StatusInternalServerError)
			return
		}
	}
}

// gets all datafiles
func GetDataFilesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		var dataFiles []models.DataFile
		if result := db.Find(&dataFiles); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching datafiles: %v", result.Error), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dataFiles); err != nil {
			http.Error(w, "Error encoding datafiles to JSON", http.StatusInternalServerError)
			return
		}
	}
}
