package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func CreateJwt(payload map[string]interface{}, secret string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(payload)).SignedString([]byte(secret))
}

func VerifyJwt(token string, secret string) (map[string]interface{}, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil })

	if err != nil {
		return nil, err
	}

	if !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return parsed.Claims.(jwt.MapClaims), nil
}
