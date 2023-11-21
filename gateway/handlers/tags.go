package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"

	"gorm.io/gorm"
)

func AddTagHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := utils.CheckRequestMethod(r, http.MethodPost); err != nil {
			utils.SendJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		type TagRequest struct {
			Name string `json:"name"`
			Type string `json:"type"`
		}

		var req TagRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.SendJSONError(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.Type == "" {
			utils.SendJSONError(w, "Tag name and type are required", http.StatusBadRequest)
			return
		}

		tag := models.Tag{
			Name: req.Name,
			Type: req.Type,
		}

		result := db.Create(&tag)
		if result.Error != nil {
			if utils.IsDuplicateKeyError(result.Error) {
				utils.SendJSONError(w, "A tag with the same name already exists", http.StatusConflict)
			} else {
				utils.SendJSONError(w, fmt.Sprintf("Error creating tag: %v", result.Error), http.StatusInternalServerError)
			}
			return
		}

		utils.SendJSONResponse(w, map[string]string{"message": fmt.Sprintf("Tag %s created successfully", tag.Name)})
	}
}
