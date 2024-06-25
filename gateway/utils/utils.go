package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func ExtractAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("authorization header format must be 'Bearer {token}'")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func SendJSONResponse(w http.ResponseWriter, response interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		SendJSONError(w, fmt.Sprintf("Error encoding response to JSON: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func SendJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": message}); err != nil {
		log.Printf("Could not encode JSON: %v", err)
	}
}

func SendJSONResponseWithCID(w http.ResponseWriter, cid string) {
	response := map[string]string{"cid": cid}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		SendJSONError(w, fmt.Sprintf("Error encoding response to JSON: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func CheckRequestMethod(r *http.Request, method string) error {
	if r.Method != method {
		return fmt.Errorf("only %s method is supported", method)
	}
	return nil
}

func ReadRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("error parsing request body: %v", err)
	}
	return nil
}

func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
		strings.Contains(err.Error(), "Duplicate entry")
}

func CreateAndWriteTempFile(r io.Reader, filename string) (*os.File, error) {
	tempFile, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %v", err)
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, r)
	if err != nil {
		return nil, fmt.Errorf("error writing temp file: %v", err)
	}

	return tempFile, nil
}

func GenerateFileHash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()

	if _, err := hasher.Write([]byte(filename)); err != nil {
		return "", fmt.Errorf("failed to hash filename: %w", err)
	}

	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to hash file contents: %w", err)
	}

	hashBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashBytes), nil
}

func GetEnvAsInt(name string, defaultValue int) int {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid format for %s. Using default value. \n", name)
		return defaultValue
	}
	return value
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
