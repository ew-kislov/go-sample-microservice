package pkg

import (
	"crypto/rand"
	"encoding/hex"
)

func RandHexString(n int) (string, error) {
	bytes := make([]byte, n/2)

	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
