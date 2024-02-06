package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labdao/plex/gateway/models"
	"gorm.io/gorm"
)

var verificationKey string
var appId string

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
	if c.Expiration < int64(time.Now().Unix()) {
		return errors.New("Token has expired.")
	}
	return nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if token.Method.Alg() != "ES256" {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	verificationKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEfhkUDY7OF5Dfx5yxehJsf7svxjOj5Ix6C+PihnsYSlsD4r8UQMu+RKYJw+Cyu2tSsvXJT7czfy0RM29YcmInrw==
-----END PUBLIC KEY-----`

	return jwt.ParseECPublicKeyFromPEM([]byte(verificationKey))
}

func SetupConfig(key string, appID string) {
	verificationKey = key
	appId = appID
	fmt.Printf("Config setup. Verification key: %v, app ID: %v\n", verificationKey, appId)
}

func validateJWT(authHeader string, w http.ResponseWriter, r *http.Request) (*PrivyClaims, bool) {
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		fmt.Println("Invalid JWT format")
		return nil, false
	}

	tokenString := splitToken[1]
	token, err := jwt.ParseWithClaims(tokenString, &PrivyClaims{}, keyFunc)
	if err != nil {
		fmt.Println("JWT signature is invalid.")
		return nil, false
	}

	privyClaim, ok := token.Claims.(*PrivyClaims)
	if !ok || privyClaim.Valid() != nil {
		fmt.Println("JWT claims are invalid.")
		return nil, false
	}

	return privyClaim, true
}

func validateAPIKey(apiKey string, db *gorm.DB) bool {
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

func GetUserByAPIKey(apiKey string, db *gorm.DB) (*models.User, error) {
	var key models.APIKey
	if err := db.Where("key = ?", apiKey).First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("API key not found")
		}
		return nil, fmt.Errorf("Error fetching API key: %v", err)
	}

	var user models.User
	if err := db.Where("id = ?", key.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("User not found for the provided API key")
		}
		return nil, fmt.Errorf("Error fetching user: %v", err)
	}

	return &user, nil
}

func AuthMiddleware(db *gorm.DB, privyPublicKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				fmt.Println("Missing Authorization header")
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			if strings.Contains(authHeader, "Bearer") {
				if _, valid := validateJWT(authHeader, w, r); !valid {
					http.Error(w, "Invalid JWT", http.StatusUnauthorized)
					return
				}
			} else {
				if !validateAPIKey(authHeader, db) {
					http.Error(w, "Invalid API Key", http.StatusUnauthorized)
					return
				}
			}

			next(w, r)
		}
	}
}
