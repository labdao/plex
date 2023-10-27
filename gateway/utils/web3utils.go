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

	batchSize := 100

	var batchAddresses []common.Address
	var batchTokenIDs []*big.Int

	for i := 0; i < int(maxTokenID.Int64()); i += batchSize {
		batchAddresses = []common.Address{}
		batchTokenIDs = []*big.Int{}

		for j := i; j < i+batchSize && j < int(maxTokenID.Int64()); j++ {
			batchAddresses = append(batchAddresses, address)
			batchTokenIDs = append(batchTokenIDs, big.NewInt(int64(j)))
		}

		callData, err := parsedABI.Pack("balanceOfBatch", batchAddresses, batchTokenIDs)
		if err != nil {
			return false, fmt.Errorf("error packing balanceOfBatch call: %w", err)
		}
		msg.Data = callData

		result, err := client.CallContract(context.Background(), msg, nil)
		if err != nil {
			return false, fmt.Errorf("error calling contract for balanceOfBatch: %w", err)
		}

		balances := []*big.Int{}
		err = parsedABI.UnpackIntoInterface(&balances, "balanceOfBatch", result)
		if err != nil {
			return false, fmt.Errorf("error unpacking balanceOfBatch result: %w", err)
		}

		for _, balance := range balances {
			if balance.Cmp(big.NewInt(0)) > 0 {
				return true, nil
			}
		}
	}

	return false, nil
}
