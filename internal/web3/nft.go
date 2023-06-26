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

type Response struct {
	Status string `json:"status"`
}

func buildTokenMetadata(ioPath, imageCid, tokenName string) (string, error) {
	ioBytes, err := ioutil.ReadFile(ioPath)
	if err != nil {
		return "", fmt.Errorf("error reading io file: %v", err)
	}

	var ioMap []map[string]interface{}

	err = json.Unmarshal(ioBytes, &ioMap)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling io file: %v", err)
	}

	graphs := []map[string]interface{}{}

	for _, ioEntry := range ioMap {
		graph := map[string]interface{}{
			"tool":    ioEntry["tool"],
			"inputs":  ioEntry["inputs"],
			"outputs": ioEntry["outputs"],
			"state":   ioEntry["state"],
			"errMsg":  ioEntry["errMsg"],
		}
		graphs = append(graphs, graph)
	}

	if !ipfs.IsValidCID(imageCid) {
		return "", fmt.Errorf("invalid image CID: %s", imageCid)
	}

	outputMap := map[string]interface{}{
		"name":        tokenName,
		"description": "Research, Reimagined. All Scientists Welcome.",
		"image":       "ipfs://" + imageCid,
		"graph":       graphs,
	}

	tokenMetadata, err := json.Marshal(outputMap)
	if err != nil {
		return "", err
	}

	return string(tokenMetadata), nil
}

func MintNFT(ioJsonPath string, imageCid, tokenName string) {
	if autotaskWebhook == "" {
		fmt.Println("AUTOTASK_WEBHOOK must be set")
		fmt.Println("Please visit https://try.labdao.xyz for instructions")
		os.Exit(1)
	}

	if recipientWallet == "" {
		fmt.Println("RECIPIENT_WALLET must be set")
		fmt.Println("Please visit https://try.labdao.xyz for instructions")
		os.Exit(1)
	}

	// Build NFT metadata
	fmt.Println("Preparing NFT metadata...")
	metadata, err := buildTokenMetadata(ioJsonPath, imageCid, tokenName)
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
	cid, err := ipfs.PinFile(tempFile.Name())
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if result.Status == "success" {
		fmt.Println("\U0001F331\U0001F331\U0001F331\U0001F331\U0001F331") // 5 saplings for open science
		fmt.Println("Minting process successful.")
		fmt.Println("Thank you for making science more reproducible, open, and collaborative!")
		fmt.Println("You can view your ProofOfScience NFT at https://testnets.opensea.io/.")
		fmt.Println("\U0001F331\U0001F331\U0001F331\U0001F331\U0001F331") // 5 more saplings for open science
	} else {
		fmt.Println("Minting process failed.")
		fmt.Printf("Response from Autotask: %s\n", string(body))
	}

	return nil
}

// TODO: move MintNFT functionality to delegatedMintNFT; add logic for which minting function to call within MintNFT

// func delegatedMintNFT() {
// 	fmt.Println("Delegated minting not yet implemented")
// }

// func userMintNFT() {
// 	fmt.Println("User minting not yet implemented")
// }
