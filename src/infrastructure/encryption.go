package infrastructure

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type EncryptionProvider struct {
}

func (*EncryptionProvider) GenerateSalt(length int64) (string, error) {
	salt := make([]byte, length/2)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(salt), nil
}

func (*EncryptionProvider) GenerateHash(payload string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(payload))
	hash.Write([]byte(salt))
	hashedPassword := hash.Sum(nil)

	return hex.EncodeToString(hashedPassword)
}

func (*EncryptionProvider) CreateJwt(payload map[string]interface{}, secret string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(payload)).SignedString([]byte(secret))
}

func (*EncryptionProvider) VerifyJwt(token string, secret string) (map[string]interface{}, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil })

	if err != nil {
		return nil, err
	}

	if !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return parsed.Claims.(jwt.MapClaims), nil
}
