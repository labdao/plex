package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func SendJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func SendJSONResponseWithCID(w http.ResponseWriter, cid string) {
	response := map[string]string{"cid": cid}
	jsonResponse, _ := json.Marshal(response)
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
