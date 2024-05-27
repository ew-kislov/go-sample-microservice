package encryption

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateSalt(length int64) (string, error) {
	salt := make([]byte, length/2)
	_, err := rand.Read(salt)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(salt), nil
}

func GenerateHash(payload, salt string) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(payload))
	_, _ = hash.Write([]byte(salt))
	hashedPassword := hash.Sum(nil)

	return hex.EncodeToString(hashedPassword)
}
