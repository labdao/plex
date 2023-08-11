package web3

import (
	"regexp"
)

func IsValidEthereumAddress(address string) bool {
	pattern := `^0x[0-9a-fA-F]{40}$`
	match, _ := regexp.MatchString(pattern, address)
	return match
}
