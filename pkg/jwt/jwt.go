package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

func CreateJwt(payload map[string]any, secret string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(payload)).SignedString([]byte(secret))
}

func VerifyJwt(token, secret string) (map[string]any, error) {
	parsed, err := jwt.Parse(token, func(_ *jwt.Token) (any, error) { return []byte(secret), nil })

	if err != nil {
		return nil, err
	}

	if !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	payload, ok := parsed.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("Could not cast JWT token claim to map")
	}

	return payload, nil
}
