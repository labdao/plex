package utils

import (
	"context"
	"fmt"
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
	file, err := os.Open("gateway/utils/clean-abi.json")
	if err != nil {
		return nil, fmt.Errorf("error opening ABI file: %w", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading ABI file: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(byteValue)))
	if err != nil {
		return nil, fmt.Errorf("error parsing ABI: %w", err)
	}

	return &parsedABI, nil
}

func CheckNFTOwnership(walletAddress string) (bool, error) {
	parsedABI, err := getContractABI()
	if err != nil {
		return false, fmt.Errorf("error getting contract ABI: %w", err)
	}

	ethClient := os.Getenv("NEXT_PUBLIC_OPTIMISM_GOERLI_ALCHEMY_RPC")
	client, err := ethclient.Dial(ethClient)
	if err != nil {
		return false, fmt.Errorf("error dialing ethClient: %w", err)
	}

	contractAddress := common.HexToAddress("0xda70C0709d4213eE8441E4731A5F662C0406ed7e")
	address := common.HexToAddress(walletAddress)

	callData, err := parsedABI.Pack("tokenID")
	if err != nil {
		return false, fmt.Errorf("error packing tokenID call: %w", err)
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return false, fmt.Errorf("error calling contract for tokenID: %w", err)
	}

	maxTokenID := new(big.Int).SetBytes(result)

	for tokenID := big.NewInt(0); tokenID.Cmp(maxTokenID) <= 0; tokenID.Add(tokenID, big.NewInt(1)) {
		callData, err := parsedABI.Pack("balanceOf", address, tokenID)
		if err != nil {
			return false, fmt.Errorf("error packing balanceOf call for tokenID %s: %w", tokenID.String(), err)
		}

		msg.Data = callData

		result, err := client.CallContract(context.Background(), msg, nil)
		if err != nil {
			return false, fmt.Errorf("error calling contract for balanceOf with tokenID %s: %w", tokenID.String(), err)
		}

		balance := new(big.Int).SetBytes(result)

		if balance.Cmp(big.NewInt(0)) > 0 {
			return true, nil
		}
	}

	return false, nil
}
