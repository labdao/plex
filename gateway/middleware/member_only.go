package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labdao/plex/gateway/models"

	"gorm.io/gorm"
)

func MemberOnlyMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			walletAddress := vars["walletAddress"]
			var user models.User

			if err := db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
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
