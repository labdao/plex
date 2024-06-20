package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labdao/plex/gateway/models"
	"github.com/labdao/plex/gateway/utils"
	"gorm.io/gorm"
)

var verificationKey string
var appId string

type contextKey string

const (
	UserContextKey contextKey = "user"
)

type PrivyClaims struct {
	SessionID  string `json:"sid,omitempty"`
	UserId     string `json:"sub,omitempty"`
	Issuer     string `json:"iss,omitempty"`
	AppId      string `json:"aud,omitempty"`
	IssuedAt   int64  `json:"iat,omitempty"`
	Expiration int64  `json:"exp,omitempty"`
}

func (c *PrivyClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{c.AppId}, nil
}

func (c *PrivyClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.IssuedAt, 0)), nil
}

func (c *PrivyClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.Expiration, 0)), nil
}

func (c *PrivyClaims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

func (c *PrivyClaims) GetJWTID() (string, error) {
	return c.SessionID, nil
}

func (c *PrivyClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c *PrivyClaims) GetSubject() (string, error) {
	return c.UserId, nil
}

func (c *PrivyClaims) Valid() error {
	if c.AppId != appId {
		return errors.New("Aud claim must be the same as Privy App ID.")
	}
	if c.Issuer != "privy.io" {
		return errors.New("Iss claim must be privy.io.")
	}
	if c.Expiration < int64(time.Now().UTC().Unix()) {
		return errors.New("Token has expired.")
	}
	return nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if token.Method.Alg() != "ES256" {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return jwt.ParseECPublicKeyFromPEM([]byte(verificationKey))
}

func SetupConfig(appID string, publicKey string) (string, string) {
	appId = appID
	verificationKey = publicKey
	fmt.Printf("Config setup. Verification Key: %v\n App ID: %v\n", verificationKey, appId)

	return appId, verificationKey
}

func IsJWT(token string) bool {
	parts := strings.Split(token, ".")
	return len(parts) == 3
}

func ValidateJWT(tokenString string, db *gorm.DB) (*PrivyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &PrivyClaims{}, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("JWT signature is invalid: %v", err)
	}

	if claims, ok := token.Claims.(*PrivyClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("JWT claims are invalid")
	}
}

func ValidateAPIKey(apiKey string, db *gorm.DB) bool {
	var key models.APIKey
	if err := db.Where("key = ?", apiKey).First(&key).Error; err != nil {
		return false
	}

	return true
}

func GetUserDIDFromRequest(r *http.Request, db *gorm.DB) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Missing Authorization header")
	} else {
		fmt.Printf("Auth header: %v\n", authHeader)
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("Invalid Authorization header")
	}

	tokenString := splitToken[1]
	fmt.Println("Token string:", tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &PrivyClaims{}, keyFunc)
	if err != nil {
		return "", fmt.Errorf("JWT signature is invalid: %v", err)
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok {
		fmt.Printf("Raw JWT claims: %+v\n", *claims)
	}

	privyClaim, ok := token.Claims.(*PrivyClaims)
	if !ok || privyClaim.Valid() != nil {
		return "", errors.New("JWT claims are invalid")
	} else {
		fmt.Printf("JWT claims are valid: %+v\n", privyClaim)
		return privyClaim.UserId, nil
	}
}

func GetUserByDID(did string, db *gorm.DB) (*models.User, error) {
	var user models.User
	if err := db.Where("did = ?", did).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("User with DID %s not found", did)
		}
		return nil, fmt.Errorf("Error fetching user: %v", err)
	}

	return &user, nil
}

func GetWalletAddressFromJWTClaims(claims *PrivyClaims, db *gorm.DB) (string, error) {
	user, err := GetUserByDID(claims.UserId, db)
	if err != nil {
		return "", err
	}
	return user.WalletAddress, nil
}

func GetUserByAPIKey(apiKey string, db *gorm.DB) (*models.User, error) {
	var key models.APIKey
	if err := db.Where("key = ?", apiKey).First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("API key not found")
		}
		return nil, fmt.Errorf("Error fetching API key: %v", err)
	}

	var user models.User
	if err := db.Where("wallet_address = ?", key.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("User not found for the provided API key")
		}
		return nil, fmt.Errorf("Error fetching user: %v", err)
	}

	return &user, nil
}

func GetWalletAddressFromAPIKey(apiKey string, db *gorm.DB) (string, error) {
	user, err := GetUserByAPIKey(apiKey, db)
	if err != nil {
		return "", err
	}

	return user.WalletAddress, nil
}

func AuthMiddleware(db *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token, err := utils.ExtractAuthHeader(r)
			if err != nil {
				fmt.Println("Error extracting JWT from header:", err)
				http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
				return
			}

			var user *models.User
			if IsJWT(token) {
				claims, err := ValidateJWT(token, db)
				if err != nil {
					fmt.Println("JWT validation error:", err)
					http.Error(w, "Invalid JWT", http.StatusUnauthorized)
					return
				}
				user, err = GetUserByDID(claims.UserId, db)
				if err != nil {
					fmt.Println("error fetching user from JWT DID:", err)
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			} else {
				user, err = GetUserByAPIKey(token, db)
				if err != nil {
					fmt.Println("error fetching user from api key:", err)
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}

			// Store the user model in the context so handlers can access user
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			r = r.WithContext(ctx)

			next(w, r)
		}
	}
}

func AdminCheckMiddleware(db *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserContextKey).(*models.User)
			if !ok {
				http.Error(w, "User context not found", http.StatusUnauthorized)
				return
			}

			if !user.Admin {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next(w, r)
		}
	}
}
