package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateCSRFState(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
