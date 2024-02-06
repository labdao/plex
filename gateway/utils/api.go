package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"strings"
)

func GenerateAPIKey(length int, label string) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("error generating random bytes: %v", err)
	}

	key := base64.URLEncoding.EncodeToString(bytes)
	checksum := crc32.ChecksumIEEE([]byte(key))
	keyWithChecksum := fmt.Sprintf("lab_%s_%x", key, checksum)

	return keyWithChecksum, nil
}

// useful for validation before DB lookup
func ValidateAPIKey(key string) (bool, error) {
	parts := strings.Split(key, "_")
	if len(parts) != 3 || parts[0] != "lab" {
		return false, fmt.Errorf("invalid API key format")
	}

	encodedPart := parts[1]
	providedChecksum := parts[2]

	var checksum uint32
	_, err := fmt.Sscanf(providedChecksum, "%x", &checksum)
	if err != nil {
		return false, fmt.Errorf("error parsing checksum: %v", err)
	}

	decoded, err := base64.URLEncoding.DecodeString(encodedPart)
	if err != nil {
		return false, fmt.Errorf("base64 decoding failed: %v", err)
	}
	computedChecksum := crc32.ChecksumIEEE(decoded)

	return computedChecksum == checksum, nil
}
