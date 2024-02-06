package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"

	"gorm.io/gorm"
)

func GenerateAPIKeyHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKeyString, err := utils.GenerateAPIKey(32, "test")
		if err != nil {
			http.Error(w, "Error generating API key", http.StatusInternalServerError)
			return
		}

		createdAt := time.Now()
		expiresAt := createdAt.Add(30 * 24 * time.Hour)

		scope := models.ScopeReadWrite

		fmt.Printf("Generated API key: %s, Scope: %s, CreatedAt: %s, ExpiresAt: %s\n", apiKeyString, scope, createdAt.Format(time.RFC3339), expiresAt.Format(time.RFC3339))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"apiKey": "%s", "scope": "%s", "createdAt": "%s", "expiresAt": "%s"}`, apiKeyString, scope, createdAt.Format(time.RFC3339), expiresAt.Format(time.RFC3339))))
	}
}
