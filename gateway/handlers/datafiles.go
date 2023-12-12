package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/ipfs"

	"log"

	"gorm.io/gorm"
)

func AddDataFilesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request to add datafiles")

		if err := utils.CheckRequestMethod(r, http.MethodPost); err != nil {
			utils.SendJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			utils.SendJSONError(w, "Error parsing multipart form", http.StatusBadRequest)
			return
		}

		files := r.MultipartForm.File["files"]
		walletAddress := r.FormValue("walletAddress")
		var successfulCIDs []string

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {

			}
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			utils.SendJSONError(w, "Error retrieving file from multipart form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		walletAddress := r.FormValue("walletAddress")
		filename := r.FormValue("filename")

		log.Printf("Received file upload request for file: %s, walletAddress: %s \n", filename, walletAddress)

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

		result := db.Create(&dataFile)
		if result.Error != nil {
			log.Println("error saving to DB")
			if utils.IsDuplicateKeyError(result.Error) {
				utils.SendJSONError(w, "A data file with the same CID already exists", http.StatusConflict)
			} else {
				utils.SendJSONError(w, fmt.Sprintf("Error saving datafile: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		var uploadedTag models.Tag
		if err := db.Where("name = ?", "uploaded").First(&uploadedTag).Error; err != nil {
			utils.SendJSONError(w, "Tag 'uploaded' not found", http.StatusInternalServerError)
			return
		}

		if err := db.Model(&dataFile).Association("Tags").Append([]models.Tag{uploadedTag}); err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error adding tag to datafile: %v", err), http.StatusInternalServerError)
			return
		}

		utils.SendJSONResponseWithCID(w, dataFile.CID)
	}
}

// Get a single datafile by CID
func GetDataFileHandler(db *gorm.DB) http.HandlerFunc {
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

		var dataFile models.DataFile
		result := db.Preload("Tags").Where("cid = ?", cid).First(&dataFile)
		if result.Error != nil {
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

func ListDataFilesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		query := db.Model(&models.DataFile{})

		if cid := r.URL.Query().Get("cid"); cid != "" {
			query = query.Where("cid = ?", cid)
		}

		if walletAddress := r.URL.Query().Get("walletAddress"); walletAddress != "" {
			query = query.Where("wallet_address = ?", walletAddress)
		}

		if filename := r.URL.Query().Get("filename"); filename != "" {
			query = query.Where("filename LIKE ?", "%"+filename+"%")
		}

		if tsBefore := r.URL.Query().Get("tsBefore"); tsBefore != "" {
			parsedTime, err := time.Parse(time.RFC3339, tsBefore)
			if err != nil {
				utils.SendJSONError(w, "Invalid timestamp format, use RFC3339 format", http.StatusBadRequest)
				return
			}
			query = query.Where("timestamp <= ?", parsedTime)
		}

		if tsAfter := r.URL.Query().Get("tsAfter"); tsAfter != "" {
			parsedTime, err := time.Parse(time.RFC3339, tsAfter)
			if err != nil {
				utils.SendJSONError(w, "Invalid timestamp format, use RFC3339 format", http.StatusBadRequest)
				return
			}
			query = query.Where("timestamp >= ?", parsedTime)
		}

		var dataFiles []models.DataFile
		if result := query.Preload("Tags").Find(&dataFiles); result.Error != nil {
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

func DownloadDataFileHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cid := vars["cid"]
		if cid == "" {
			utils.SendJSONError(w, "Missing CID parameter", http.StatusBadRequest)
			return
		}

		var dataFile models.DataFile
		if err := db.Where("cid = ?", cid).First(&dataFile).Error; err != nil {
			utils.SendJSONError(w, "Data file not found", http.StatusNotFound)
			return
		}

		ipfsPath := cid
		if dataFile.WalletAddress != "" {
			ipfsPath = cid + "/" + dataFile.Filename
		}

		tempFilePath, err := ipfs.DownloadFileToTemp(ipfsPath, dataFile.Filename)
		if err != nil {
			utils.SendJSONError(w, "Error downloading file from IPFS", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFilePath)

		file, err := os.Open(tempFilePath)
		if err != nil {
			utils.SendJSONError(w, "Error opening downloaded file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Disposition", "attachment; filename="+dataFile.Filename)
		w.Header().Set("Content-Type", "application/octet-stream")

		if _, err := io.Copy(w, file); err != nil {
			utils.SendJSONError(w, "Error sending file", http.StatusInternalServerError)
			return
		}
	}
}

func AddTagsToDataFile(db *gorm.DB, dataFileCID string, tagNames []string) error {
	log.Println("Starting AddTagsToDataFile for DataFile with CID:", dataFileCID)

	var dataFile models.DataFile
	if err := db.Preload("Tags").Where("cid = ?", dataFileCID).First(&dataFile).Error; err != nil {
		log.Printf("Error finding DataFile with CID %s: %v\n", dataFileCID, err)
		return fmt.Errorf("Data file not found: %v", err)
	}

	var tags []models.Tag
	if err := db.Where("name IN ?", tagNames).Find(&tags).Error; err != nil {
		log.Printf("Error finding tags: %v\n", err)
		return fmt.Errorf("Error finding tags: %v", err)
	}

	existingTagMap := make(map[string]bool)
	for _, tag := range dataFile.Tags {
		existingTagMap[tag.Name] = true
	}

	log.Println("Adding tags:", tagNames)
	for _, tag := range tags {
		if !existingTagMap[tag.Name] {
			dataFile.Tags = append(dataFile.Tags, tag)
		}
	}

	log.Println("Saving DataFile with new tags to DB")
	if err := db.Save(&dataFile).Error; err != nil {
		log.Printf("Error saving DataFile with CID %s: %v\n", dataFileCID, err)
		return fmt.Errorf("Error saving datafile: %v", err)
	}

	log.Println("DataFile with CID", dataFileCID, "successfully updated with new tags")
	return nil
}
