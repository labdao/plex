package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/ipfs"

	"log"

	"gorm.io/gorm"
)

func AddDataFileHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request to add datafile")

		if err := utils.CheckRequestMethod(r, http.MethodPost); err != nil {
			utils.SendJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			utils.SendJSONError(w, "Error parsing multipart form", http.StatusBadRequest)
			return
		}
		log.Println("Parsed multipart form")

		file, _, err := r.FormFile("file")
		if err != nil {
			utils.SendJSONError(w, "Error retrieving file from multipart form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
		if !ok {
			utils.SendJSONError(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		walletAddress := user.WalletAddress

		filename := r.FormValue("filename")
		publicValue := r.FormValue("public")

		isPublic, err := strconv.ParseBool(publicValue)
		if err != nil {
			isPublic = false
		}

		// Silently ignore public flag if user is not an admin
		// We still allow them to upload, it will just be private
		// Data files can be made public later by making experiment results public
		if !user.Admin {
			isPublic = false
		}

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
			Public:        isPublic,
		}

		var existingDataFile models.DataFile
		if err := db.Where("cid = ?", cid).First(&existingDataFile).Error; err == nil {
			var userHasDataFile bool
			var count int64
			db.Model(&dataFile).Joins("JOIN user_datafiles ON user_datafiles.c_id = data_files.cid").
				Where("user_datafiles.wallet_address = ? AND data_files.cid = ?", user.WalletAddress, cid).First(&dataFile).Count(&count)
			userHasDataFile = count > 0
			if userHasDataFile {
				utils.SendJSONError(w, "A user data file with the same CID already exists", http.StatusConflict)
				return
			} else {
				if err := db.Model(&user).Association("UserDatafiles").Append(&existingDataFile); err != nil {
					utils.SendJSONError(w, fmt.Sprintf("Error associating datafile with user: %v", err), http.StatusInternalServerError)
					return
				}
			}
			if isPublic && !existingDataFile.Public {
				existingDataFile.Public = true
				if err := db.Save(&existingDataFile).Error; err != nil {
					utils.SendJSONError(w, fmt.Sprintf("Error updating datafile public status: %v", err), http.StatusInternalServerError)
					return
				}
			}
		} else {
			result := db.Create(&dataFile)
			if result.Error != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error saving datafile: %v", result.Error), http.StatusInternalServerError)
				return
			}

			if err := db.Model(&user).Association("UserDatafiles").Append(&dataFile); err != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error associating datafile with user: %v", err), http.StatusInternalServerError)
				return
			}
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

func GetDataFileHandler(db *gorm.DB) http.HandlerFunc {
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

		var userAssociatedWithFile bool
		db.Model(&user).Association("UserDatafiles").Find(&dataFile, models.DataFile{CID: cid})
		userAssociatedWithFile = !errors.Is(db.Error, gorm.ErrRecordNotFound)

		if !userAssociatedWithFile && !dataFile.Public {
			utils.SendJSONError(w, "Unauthorized access or data file not found", http.StatusUnauthorized)
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

		user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
		if !ok {
			utils.SendJSONError(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		var page, pageSize int = 1, 50

		if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
			page = p
		}
		if ps, err := strconv.Atoi(r.URL.Query().Get("pageSize")); err == nil && ps > 0 {
			pageSize = ps
		}

		offset := (page - 1) * pageSize

		query := db.Model(&models.DataFile{}).
			Joins("LEFT JOIN user_datafiles ON user_datafiles.data_file_c_id = data_files.cid AND user_datafiles.wallet_address = ?", user.WalletAddress).
			Where("data_files.public = true OR (data_files.public = false AND user_datafiles.wallet_address = ?)", user.WalletAddress)

		if cid := r.URL.Query().Get("cid"); cid != "" {
			query = query.Where("data_files.cid = ?", cid)
		}
		if filename := r.URL.Query().Get("filename"); filename != "" {
			query = query.Where("data_files.filename LIKE ?", "%"+filename+"%")
		}
		if tsBefore := r.URL.Query().Get("tsBefore"); tsBefore != "" {
			parsedTime, err := time.Parse(time.RFC3339, tsBefore)
			if err != nil {
				utils.SendJSONError(w, "Invalid timestamp format, use RFC3339 format", http.StatusBadRequest)
				return
			}
			query = query.Where("user_datafiles.timestamp <= ?", parsedTime)
		}
		if tsAfter := r.URL.Query().Get("tsAfter"); tsAfter != "" {
			parsedTime, err := time.Parse(time.RFC3339, tsAfter)
			if err != nil {
				utils.SendJSONError(w, "Invalid timestamp format, use RFC3339 format", http.StatusBadRequest)
				return
			}
			query = query.Where("user_datafiles.timestamp >= ?", parsedTime)
		}

		var totalCount int64
		query.Count(&totalCount)

		defaultSort := "data_files.timestamp desc"
		sortParam := r.URL.Query().Get("sort")
		if sortParam != "" {
			defaultSort = sortParam
		}
		query = query.Order(defaultSort).Offset(offset).Limit(pageSize)

		var dataFiles []models.DataFile
		if result := query.Preload("Tags").Find(&dataFiles); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching datafiles: %v", result.Error), http.StatusInternalServerError)
			return
		}

		totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

		response := map[string]interface{}{
			"data": dataFiles,
			"pagination": map[string]int{
				"currentPage": page,
				"totalPages":  totalPages,
				"pageSize":    pageSize,
				"totalCount":  int(totalCount),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response to JSON", http.StatusInternalServerError)
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

		user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
		if !ok {
			utils.SendJSONError(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		var dataFile models.DataFile
		if err := db.Preload("Tags").Where("cid = ?", cid).First(&dataFile).Error; err != nil {
			utils.SendJSONError(w, "Data file not found", http.StatusNotFound)
			return
		}

		var userAssociatedWithFile bool
		db.Model(&user).Association("UserDatafiles").Find(&dataFile, models.DataFile{CID: cid})
		userAssociatedWithFile = !errors.Is(db.Error, gorm.ErrRecordNotFound)

		if !userAssociatedWithFile && !dataFile.Public {
			utils.SendJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// First attempt with the initial ipfsPath
		ipfsPath := determineIPFSPath(cid, dataFile)
		tempFilePath, err := ipfs.DownloadFileToTemp(ipfsPath, dataFile.Filename)
		if err != nil {
			// If the first attempt fails, try with an alternative ipfsPath
			altIPFSPath := determineAltIPFSPath(cid, dataFile)
			tempFilePath, err = ipfs.DownloadFileToTemp(altIPFSPath, dataFile.Filename)
			if err != nil {
				utils.SendJSONError(w, "Error downloading file from IPFS", http.StatusInternalServerError)
				return
			}
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

func UpdateDataFileHandler(db *gorm.DB) http.HandlerFunc {
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

		vars := mux.Vars(r)
		cid := vars["cid"]
		if cid == "" {
			utils.SendJSONError(w, "Missing CID parameter", http.StatusBadRequest)
			return
		}

		var dataFile models.DataFile
		result := db.Where("cid = ? AND public = false", cid).First(&dataFile)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				utils.SendJSONError(w, "Data file not found or already public", http.StatusNotFound)
			} else {
				utils.SendJSONError(w, fmt.Sprintf("Error fetching datafile: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		if dataFile.WalletAddress != user.WalletAddress {
			utils.SendJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		dataFile.Public = true
		if err := db.Save(&dataFile).Error; err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error updating datafile: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dataFile); err != nil {
			utils.SendJSONError(w, "Error encoding datafile to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func determineIPFSPath(cid string, dataFile models.DataFile) string {
	isGenerated := checkIfGenerated(dataFile)
	if dataFile.WalletAddress != "" || isGenerated {
		return cid + "/" + dataFile.Filename
	}
	return cid
}

func determineAltIPFSPath(cid string, dataFile models.DataFile) string {
	isGenerated := checkIfGenerated(dataFile)
	if dataFile.WalletAddress == "" && !isGenerated {
		return cid + "/" + dataFile.Filename
	}
	return cid
}

func checkIfGenerated(dataFile models.DataFile) bool {
	for _, tag := range dataFile.Tags {
		if tag.Name == "generated" {
			return true
		}
	}
	return false
}

func AddTagsToDataFile(db *gorm.DB, dataFileCID string, tagNames []string) error {
	log.Println("Starting AddTagsToDataFile for DataFile with CID:", dataFileCID)

	var dataFile models.DataFile
	if err := db.Preload("Tags").Where("cid = ?", dataFileCID).First(&dataFile).Error; err != nil {
		log.Printf("Error finding DataFile with CID %s: %v\n", dataFileCID, err)
		return fmt.Errorf("data file not found: %v", err)
	}

	var tags []models.Tag
	if err := db.Where("name IN ?", tagNames).Find(&tags).Error; err != nil {
		log.Printf("Error finding tags: %v\n", err)
		return fmt.Errorf("error finding tags: %v", err)
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
		log.Printf("error saving DataFile with CID %s: %v\n", dataFileCID, err)
		return fmt.Errorf("error saving datafile: %v", err)
	}

	log.Println("DataFile with CID", dataFileCID, "successfully updated with new tags")

	return nil
}
