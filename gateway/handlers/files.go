package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/s3"

	"log"

	"gorm.io/gorm"
)

func AddFileHandler(db *gorm.DB, s3c *s3.S3Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request to add file")

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

		retrievedFile, _, err := r.FormFile("file")
		if err != nil {
			utils.SendJSONError(w, "Error retrieving file from multipart form", http.StatusBadRequest)
			return
		}
		defer retrievedFile.Close()

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
		// Files can be made public later by making experiment results public
		if !user.Admin {
			isPublic = false
		}

		log.Printf("Received file upload request for file: %s, walletAddress: %s \n", filename, walletAddress)

		tempFile, err := utils.CreateAndWriteTempFile(retrievedFile, filename)
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error creating temp file: %v", err), http.StatusInternalServerError)
			return
		}

		bucketName := os.Getenv("BUCKET_NAME")
		if bucketName == "" {
			utils.SendJSONError(w, "BUCKET_NAME environment variable not set", http.StatusInternalServerError)
			return
		}

		hash, err := utils.GenerateFileHash(tempFile.Name())
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error generating hash: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Println("Hash of", filename, "is", hash)
		defer os.Remove(filename)

		objectKey := hash + "/" + filename
		// S3 upload
		err = s3c.UploadFile(bucketName, objectKey, tempFile.Name())
		if err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error uploading file to bucket: %v", err), http.StatusInternalServerError)
			return
		}

		s3_uri := fmt.Sprintf("s3://%s/%s", bucketName, objectKey)
		file := models.File{
			FileHash:      hash,
			WalletAddress: user.WalletAddress,
			Filename:      filename,
			CreatedAt:     time.Now().UTC(),
			Public:        isPublic,
			S3URI:         s3_uri,
		}

		var existingFile models.File
		if err := db.Where("file_hash = ?", hash).First(&existingFile).Error; err == nil {
			var userHasFile bool
			var count int64
			db.Model(&file).Joins("JOIN user_files ON user_files.file_id = files.id").
				Where("user_files.wallet_address = ? AND files.id = ?", user.WalletAddress, hash).First(&file).Count(&count)
			userHasFile = count > 0
			if userHasFile {
				utils.SendJSONError(w, "A user file with the same ID already exists", http.StatusConflict)
				return
			} else {
				if err := db.Model(&user).Association("UserFiles").Append(&existingFile); err != nil {
					utils.SendJSONError(w, fmt.Sprintf("Error associating file with user: %v", err), http.StatusInternalServerError)
					return
				}
			}
			if isPublic && !existingFile.Public {
				existingFile.Public = true
				if err := db.Save(&existingFile).Error; err != nil {
					utils.SendJSONError(w, fmt.Sprintf("Error updating file public status: %v", err), http.StatusInternalServerError)
					return
				}
			}
		} else {
			result := db.Create(&file)
			if result.Error != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error saving file: %v", result.Error), http.StatusInternalServerError)
				return
			}

			if err := db.Model(&user).Association("UserFiles").Append(&file); err != nil {
				utils.SendJSONError(w, fmt.Sprintf("Error associating file with user: %v", err), http.StatusInternalServerError)
				return
			}
		}

		var uploadedTag models.Tag
		if err := db.Where("name = ?", "uploaded").First(&uploadedTag).Error; err != nil {
			utils.SendJSONError(w, "Tag 'uploaded' not found", http.StatusInternalServerError)
			return
		}

		if err := db.Model(&file).Association("Tags").Append([]models.Tag{uploadedTag}); err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error adding tag to file: %v", err), http.StatusInternalServerError)
			return
		}

		utils.SendJSONResponseWithID(w, file.ID)
	}
}

func GetFileHandler(db *gorm.DB) http.HandlerFunc {
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
		id := vars["id"]
		if id == "" {
			utils.SendJSONError(w, "Missing ID parameter", http.StatusBadRequest)
			return
		}

		var file models.File
		result := db.Preload("Tags").Where("id = ?", id).First(&file)
		if result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching file: %v", result.Error), http.StatusInternalServerError)
			return
		}

		var userAssociatedWithFile bool
		db.Model(&user).Association("UserFiles").Find(&file, models.File{ID: file.ID})
		userAssociatedWithFile = !errors.Is(db.Error, gorm.ErrRecordNotFound)

		if !userAssociatedWithFile && !file.Public {
			utils.SendJSONError(w, "Unauthorized access or file not found", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(file); err != nil {
			http.Error(w, "Error encoding file to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func ListFilesHandler(db *gorm.DB) http.HandlerFunc {
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

		query := db.Model(&models.File{}).
			Joins("LEFT JOIN user_files ON user_files.file_id = files.id AND user_files.wallet_address = ?", user.WalletAddress).
			Where("files.public = true OR (files.public = false AND user_files.wallet_address = ?)", user.WalletAddress)

		if id := r.URL.Query().Get("id"); id != "" {
			query = query.Where("files.id = ?", id)
		}
		if filename := r.URL.Query().Get("filename"); filename != "" {
			query = query.Where("files.filename LIKE ?", "%"+filename+"%")
		}
		if tsBefore := r.URL.Query().Get("tsBefore"); tsBefore != "" {
			parsedTime, err := time.Parse(time.RFC3339, tsBefore)
			if err != nil {
				utils.SendJSONError(w, "Invalid timestamp format, use RFC3339 format", http.StatusBadRequest)
				return
			}
			query = query.Where("user_files.created_at <= ?", parsedTime)
		}
		if tsAfter := r.URL.Query().Get("tsAfter"); tsAfter != "" {
			parsedTime, err := time.Parse(time.RFC3339, tsAfter)
			if err != nil {
				utils.SendJSONError(w, "Invalid timestamp format, use RFC3339 format", http.StatusBadRequest)
				return
			}
			query = query.Where("user_files.created_at >= ?", parsedTime)
		}

		var totalCount int64
		query.Count(&totalCount)

		defaultSort := "files.created_at desc"
		sortParam := r.URL.Query().Get("sort")
		if sortParam != "" {
			defaultSort = sortParam
		}
		query = query.Order(defaultSort).Offset(offset).Limit(pageSize)

		var files []models.File
		if result := query.Preload("Tags").Find(&files); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching files: %v", result.Error), http.StatusInternalServerError)
			return
		}

		totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

		response := map[string]interface{}{
			"data": files,
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

func DownloadFileHandler(db *gorm.DB, s3c *s3.S3Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		if id == "" {
			utils.SendJSONError(w, "Missing ID parameter", http.StatusBadRequest)
			return
		}

		user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
		if !ok {
			utils.SendJSONError(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		var file models.File
		if err := db.Where("id = ?", id).First(&file).Error; err != nil {
			utils.SendJSONError(w, "File not found", http.StatusNotFound)
			return
		}

		var userAssociatedWithFile bool
		db.Model(&user).Association("UserFiles").Find(&file, models.File{ID: file.ID})
		userAssociatedWithFile = !errors.Is(db.Error, gorm.ErrRecordNotFound)

		if !userAssociatedWithFile && !file.Public {
			utils.SendJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//if the file as S3Bucket and S3Location, download from S3, else error
		if file.S3URI == "" {
			utils.SendJSONError(w, "S3URI is null, fetch from IPFS is not supported", http.StatusNotFound)
		} else {
			filename := file.Filename

			if err := s3c.StreamFileToResponse(file.S3URI, w, filename); err != nil {
				utils.SendJSONError(w, "Error downloading file from S3", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+file.Filename)
		w.Header().Set("Content-Type", "application/octet-stream")

	}
}

func UpdateFileHandler(db *gorm.DB) http.HandlerFunc {
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
		id := vars["id"]
		if id == "" {
			utils.SendJSONError(w, "Missing ID parameter", http.StatusBadRequest)
			return
		}

		var file models.File
		result := db.Where("id = ? AND public = false", id).First(&file)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				utils.SendJSONError(w, "File not found or already public", http.StatusNotFound)
			} else {
				utils.SendJSONError(w, fmt.Sprintf("Error fetching file: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		if file.WalletAddress != user.WalletAddress {
			utils.SendJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		file.Public = true
		if err := db.Save(&file).Error; err != nil {
			utils.SendJSONError(w, fmt.Sprintf("Error updating file: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(file); err != nil {
			utils.SendJSONError(w, "Error encoding file to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func AddTagsToFile(db *gorm.DB, fileID int, tagNames []string) error {
	log.Println("Starting AddTagsToFile for File with ID:", fileID)

	var file models.File
	if err := db.Preload("Tags").Where("id = ?", fileID).First(&file).Error; err != nil {
		log.Printf("Error finding File with ID %d: %v\n", fileID, err)
		return fmt.Errorf("file not found: %v", err)
	}

	var tags []models.Tag
	if err := db.Where("name IN ?", tagNames).Find(&tags).Error; err != nil {
		log.Printf("Error finding tags: %v\n", err)
		return fmt.Errorf("error finding tags: %v", err)
	}

	existingTagMap := make(map[string]bool)
	for _, tag := range file.Tags {
		existingTagMap[tag.Name] = true
	}

	log.Println("Adding tags:", tagNames)
	for _, tag := range tags {
		if !existingTagMap[tag.Name] {
			file.Tags = append(file.Tags, tag)
		}
	}

	log.Println("Saving File with new tags to DB")
	if err := db.Save(&file).Error; err != nil {
		log.Printf("error saving File with ID %d: %v\n", fileID, err)
		return fmt.Errorf("error saving file: %v", err)
	}

	log.Println("File with ID", fileID, "successfully updated with new tags")

	return nil
}
