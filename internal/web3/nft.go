package web3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labdao/plex/internal/ipfs"
)

var recipientWallet = os.Getenv("RECIPIENT_WALLET")
var autotaskWebhook = os.Getenv("AUTOTASK_WEBHOOK")

type postData struct {
	RecipientAddress string `json:"recipientAddress"`
	Cid              string `json:"cid"`
}

func removeFilepathKeys(obj map[string]interface{}) {
	delete(obj, "filepath")

	for _, value := range obj {
		if asMap, ok := value.(map[string]interface{}); ok {
			removeFilepathKeys(asMap)
		} else if asSlice, ok := value.([]interface{}); ok {
			for _, itemInSlice := range asSlice {
				if asMap, ok := itemInSlice.(map[string]interface{}); ok {
					removeFilepathKeys(asMap)
				}
			}
		}
	}
}

func buildTokenMetadata(toolPath, ioPath string) (string, error) {
	toolBytes, err := ioutil.ReadFile(toolPath)
	if err != nil {
		return "", fmt.Errorf("error reading tool file: %v", err)
	}

	ioBytes, err := ioutil.ReadFile(ioPath)
	if err != nil {
		return "", fmt.Errorf("error reading io file: %v", err)
	}

	var toolMap map[string]interface{}
	var ioMap []map[string]interface{}

	err = json.Unmarshal(toolBytes, &toolMap)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling tool file: %v", err)
	}
	removeFilepathKeys(toolMap)

	err = json.Unmarshal(ioBytes, &ioMap)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling io file: %v", err)
	}

	for _, ioEntry := range ioMap {
		removeFilepathKeys(ioEntry)
	}

	tokenName := GenerateTokenName()

	graphs := []map[string]interface{}{}

	for _, ioEntry := range ioMap {
		graph := map[string]interface{}{
			"tool":    toolMap,
			"inputs":  ioEntry["inputs"],
			"outputs": ioEntry["outputs"],
			"state":   ioEntry["state"],
			"errMsg":  ioEntry["errMsg"],
		}
		graphs = append(graphs, graph)
	}

	outputMap := map[string]interface{}{
		"name":        tokenName,
		"description": "Research, Reimagined. All Scientists Welcome.",
		"image":       "ipfs://bafybeiba666bzbff5vu6rayvp5st2tk7tdltqnwjppzyvpljcycfhshdhq/",
		"graph":       graphs,
	}

	tokenMetadata, err := json.Marshal(outputMap)
	if err != nil {
		return "", err
	}

	return string(tokenMetadata), nil
}

func MintNFT(toolPath, ioJsonPath string) {
	if recipientWallet == "" {
		fmt.Println("RECIPIENT_WALLET must be set")
		os.Exit(1)
	}

	if autotaskWebhook == "" {
		fmt.Println("AUTOTASK_WEBHOOK must be set")
		os.Exit(1)
	}

	// Build NFT metadata
	fmt.Println("Preparing NFT metadata...")
	metadata, err := buildTokenMetadata(toolPath, ioJsonPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	tempFile, err := ioutil.TempFile("", "metadata-*.json")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(metadata)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	err = tempFile.Close()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Uploading NFT metadata to IPFS...")
	cid, err := ipfs.GetFileCid(tempFile.Name())
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("NFT metadata uploaded to IPFS: ipfs://%s\n", cid)

	fmt.Println("Triggering minting process via Defender Autotask...")
	err = triggerMinting(recipientWallet, cid)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func triggerMinting(recipientAddress, cid string) error {
	url := autotaskWebhook

	data := postData{
		RecipientAddress: recipientAddress,
		Cid:              cid,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO: Handle response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Response from Autotask: %s\n", string(body))

	return nil
}
