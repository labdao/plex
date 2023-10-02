package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labdao/plex/gateway/utils"
	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"
	"gorm.io/gorm"
)

func pinIoList(ios []ipwl.IO) (string, error) {
	// Convert IO slice to JSON
	data, err := json.Marshal(ios)
	if err != nil {
		return "", fmt.Errorf("failed to marshal IO slice: %v", err)
	}

	// Create a temporary file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	if err != nil {
		return "", fmt.Errorf("cannot create temporary file: %v", err)
	}

	// Write JSON data to the temporary file
	if _, err = tmpFile.Write(data); err != nil {
		return "", fmt.Errorf("failed to write to temporary file: %v", err)
	}

	cid, err := ipfs.PinFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to pin file: %v", err)
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close the file: %v", err)
	}

	return cid, nil
}

func AddGraphHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received Post request at /graph")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		log.Println("Request body: ", string(body))

		requestData := make(map[string]json.RawMessage)
		err = json.Unmarshal(body, &requestData)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var walletAddress string
		err = json.Unmarshal(requestData["walletAddress"], &walletAddress)
		if err != nil || walletAddress == "" {
			http.Error(w, "Invalid or missing walletAddress", http.StatusBadRequest)
			return
		}

		var toolCid string
		err = json.Unmarshal(requestData["toolCid"], &toolCid)
		if err != nil || walletAddress == "" {
			http.Error(w, "Invalid or missing Tool CID", http.StatusBadRequest)
			return
		}

		var scatteringMethod string
		err = json.Unmarshal(requestData["scatteringMethod"], &scatteringMethod)
		if err != nil || walletAddress == "" {
			http.Error(w, "Invalid or missing Scattering Method", http.StatusBadRequest)
			return
		}

		var kwargs map[string][]string
		err = json.Unmarshal(requestData["kwargs"], &kwargs)
		if err != nil {
			http.Error(w, "Invalid or missing kwargs", http.StatusBadRequest)
			return
		}

		// add wallet
		ioList, err := ipwl.InitializeIo(toolCid, scatteringMethod, kwargs)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Initialized IO List")

		log.Println("Submitting IO List")
		submittedIoList := ipwl.SubmitIoList(ioList, "", 60, []string{})
		log.Println("pinning submitted IO List")
		submittedIoListCid, err := pinIoList(submittedIoList)
		if err != nil {
			log.Fatal(err)
		}
		utils.SendJSONResponseWithCID(w, submittedIoListCid)
	}
}
