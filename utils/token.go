package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken creates a random hex token of length 32 bytes (64 hex chars).
func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
