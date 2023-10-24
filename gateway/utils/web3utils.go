package utils

import (
	"context"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getContractABI() (*abi.ABI, error) {
	file, err := os.Open("../../contracts/artifacts/build-info/a7fa49d6da0732841f00ec7c540205a6.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(byteValue)))
	if err != nil {
		return nil, err
	}

	return &parsedABI, nil
}

func CheckNFTOwnership(walletAddress string) (bool, error) {
	parsedABI, err := getContractABI()
	if err != nil {
		return false, err
	}

	ethClient := os.Getenv("OPTIMISM_GOERLI_ALCHEMY_RPC")
	client, err := ethclient.Dial(ethClient)
	if err != nil {
		return false, err
	}

	contractAddress := common.HexToAddress("0xda70C0709d4213eE8441E4731A5F662C0406ed7e")
	address := common.HexToAddress(walletAddress)

	// get maximum token ID
	callData, err := parsedABI.Pack("tokenID")
	if err != nil {
		return false, err
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return false, err
	}

	maxTokenID := new(big.Int)
	err = parsedABI.UnpackIntoInterface(maxTokenID, "tokenID", result)
	if err != nil {
		return false, err
	}

	for tokenID := big.NewInt(0); tokenID.Cmp(maxTokenID) <= 0; tokenID.Add(tokenID, big.NewInt(1)) {
		callData, err := parsedABI.Pack("balanceOf", address, tokenID)
		if err != nil {
			return false, err
		}

		msg.Data = callData

		result, err := client.CallContract(context.Background(), msg, nil)
		if err != nil {
			return false, err
		}

		balance := new(big.Int)
		err = parsedABI.UnpackIntoInterface(balance, "balanceOf", result)
		if err != nil {
			return false, err
		}

		if balance.Cmp(big.NewInt(0)) > 0 {
			return true, nil
		}
	}

	return false, nil
}
