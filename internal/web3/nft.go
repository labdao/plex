package web3

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/labdao/plex/internal/ipfs"
)

// check if wallet address is valid
// func isValidAddress(address string) bool {
// 	if !common.IsHexAddress(address) {
// 		return false
// 	}
// 	return common.IsHexAddress(common.HexToAddress(address).Hex())
// }

// type RelayerRequest struct {
// 	To   string `json:"to"`
// 	Data string `json:"data"`
// }

// func SubmitTransaction() {
// 	contractAddress := common.HexToAddress("contractAddress")
// 	contractAbi, err := abi.JSON(strings.NewReader("contractAbi"))
// 	if err != nil {
// 		log.Fatalf("Failed to parse contract ABI: %v", err)
// 	}

// 	// Generate call data
// 	data, err := contractAbi.Pack("functionName", "param1", "param2")
// 	if err != nil {
// 		log.Fatalf("Failed to generate call data: %v", err)
// 	}

// 	// Create Ethereum transaction
// 	tx := types.NewTransaction(0, contractAddress, big.NewInt(0), 100000, big.NewInt(0), data)

// 	// Convert the transaction to raw RLP-encoded bytes
// 	encodedTx, err := tx.MarshalBinary()
// 	if err != nil {
// 		log.Fatalf("Failed to encode transaction: %v", err)
// 	}

// 	// Create the Defender relayer request
// 	request := RelayerRequest{
// 		To:   contractAddress.Hex(),
// 		Data: common.Bytes2Hex(encodedTx),
// 	}

// 	// Convert the request to JSON
// 	jsonRequest, err := json.Marshal(request)
// 	if err != nil {
// 		log.Fatalf("Failed to encode request to JSON: %v", err)
// 	}

// 	// Make the HTTP POST request to the Defender relayer
// 	resp, err := http.Post("https://defender.openzeppelin.com/relayerApi/yourRelayer", "application/json", bytes.NewBuffer(jsonRequest))
// 	if err != nil {
// 		log.Fatalf("Failed to send request to relayer: %v", err)
// 	}
// }

// // Define contract details
// contractAddress := common.HexToAddress("contractAddress")
// contractAbi, err := abi.JSON(strings.NewReader("contractAbi"))
// if err != nil {
// 	log.Fatalf("Failed to parse contract ABI: %v", err)
// }

func removeFilepathKeys(obj map[string]interface{}) {
	delete(obj, "filepath") // Remove the filepath key at the current level

	for _, value := range obj {
		// Check if the value is another map
		if asMap, ok := value.(map[string]interface{}); ok {
			removeFilepathKeys(asMap) // Recursively check this map
		} else if asSlice, ok := value.([]interface{}); ok { // Check if the value is a slice
			for _, itemInSlice := range asSlice {
				// Check if the item in the slice is a map
				if asMap, ok := itemInSlice.(map[string]interface{}); ok {
					removeFilepathKeys(asMap) // Recursively check this map
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
	// Build NFT metadata
	fmt.Println("Preparing NFT metadata...")
	metadata, err := buildTokenMetadata(toolPath, ioJsonPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Create a temporary file
	tempFile, err := ioutil.TempFile("", "metadata-*.json")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer os.Remove(tempFile.Name()) // clean up

	// Write the metadata to the temporary file
	_, err = tempFile.WriteString(metadata)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Close the file
	err = tempFile.Close()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Upload the metadata to IPFS and return the CID
	fmt.Println("Uploading NFT metadata to IPFS...")
	cid, err := ipfs.GetFileCid(tempFile.Name())
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("NFT metadata uploaded to IPFS: ipfs://%s\n", cid)
}
