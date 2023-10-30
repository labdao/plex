package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labdao/plex/gateway/models"

	"gorm.io/gorm"
)

type RequestBody struct {
	WalletAddress string `json:"walletAddress"`
}

func MemberOnlyMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			// Replace the request body so it can be read again in the main handler
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			var body RequestBody
			err = json.Unmarshal(bodyBytes, &body)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			walletAddress := body.WalletAddress
			var user models.User

			if err := db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					http.Error(w, "User not found", http.StatusUnauthorized)
				} else {
					http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
				}
				return
			}

			if !user.IsMember {
				http.Error(w, "Access denied: User is not a member", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r) // Call the next handler if user is a member
		})
	}
}
