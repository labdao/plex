package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

verificationKey := os.Getenv("PRIVY_VERIFICATION_KEY")
appId := os.Getenv("PRIVY_APP_ID")

type PrivyClaims struct {
	AppId      string `json:"aud,omitempty"`
	Expiration uint64 `json:"exp,omitempty"`
	Issuer     string `json:"iss,omitempty"`
	UserId     string `json:"sub,omitempty"`
}

func (c *PrivyClaims) Valid() error {
	if c.AppId != appId {

	}
	return nil
}

func JWTMiddleware(db *gorm.DB, privyPublicKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			// splitToken := strings.Split(authHeader, "Bearer ")
			// if len(splitToken) != 2 {
			// 	http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			// 	return
			// }

			// tokenString := splitToken[1]

			// claims, err := models.ParseJWT(tokenString, privyPublicKey)
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusUnauthorized)
			// 	return
			// }

			// user, err := models.GetUserByWalletAddress(db, claims.WalletAddress)
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusUnauthorized)
			// 	return
			// }

			// r = r.WithContext(models.SetUserContext(r.Context(), user))

			next.ServeHTTP(w, r)
		})
	}
}
