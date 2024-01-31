package handlers

import (
	"net/http"

	"gorm.io/gorm"
)

func GenerateAPIKeyHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
