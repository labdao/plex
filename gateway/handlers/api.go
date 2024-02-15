package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labdao/plex/gateway/middleware"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"

	"gorm.io/gorm"
)

func AddAPIKeyHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.SendJSONError(w, "Only POST method is supported", http.StatusBadRequest)
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

		apiKeyString, err := utils.GenerateAPIKey(32, walletAddress)
		if err != nil {
			utils.SendJSONError(w, "Error generating API key", http.StatusInternalServerError)
			return
		}

		apiKey := models.APIKey{
			Key:       apiKeyString,
			Scope:     models.ScopeReadWrite, // default scope is read-write
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // default expiration time is 30 days
			UserID:    walletAddress,
		}

		result := db.Create(&apiKey)
		if result.Error != nil {
			utils.SendJSONError(w, "Error saving API key to database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(apiKey)
	}
}

func ListAPIKeysHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.SendJSONError(w, "Only GET method is supported", http.StatusBadRequest)
			return
		}

		did, err := middleware.GetUserDIDFromRequest(r, db)
		if err != nil {
			utils.SendJSONError(w, "Failed to get user DID from request: "+err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := middleware.GetUserByDID(did, db)
		if err != nil {
			utils.SendJSONError(w, "Failed to get user by DID: "+err.Error(), http.StatusNotFound)
			return
		}

		var apiKeys []models.APIKey
		if err := db.Where("user_id = ?", user.WalletAddress).Find(&apiKeys).Error; err != nil {
			utils.SendJSONError(w, "Failed to get API keys: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(apiKeys); err != nil {
			utils.SendJSONError(w, "Error encoding API keys to JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
