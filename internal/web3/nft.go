package web3

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// check if wallet address is valid
// func isValidAddress(address string) bool {
// 	if !common.IsHexAddress(address) {
// 		return false
// 	}
// 	return common.IsHexAddress(common.HexToAddress(address).Hex())
// }

func BuildTokenMetadata(toolPath, ioPath string) (string, error) {
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

	err = json.Unmarshal(ioBytes, &ioMap)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling io file: %v", err)
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
