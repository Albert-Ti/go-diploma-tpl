package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHash(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
