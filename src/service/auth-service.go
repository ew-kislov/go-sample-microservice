package service

import (
	"context"
	"go-sample-microservice/src/infrastructure"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

type CreateUserParams struct {
	Email       string
	Username    string
	DisplayName string
	Password    string
}

type TokenPayload struct {
	Id int64 `mapstructure:"id"`
}

type SignUpResponse struct {
	UserId int64  `json:"userId"`
	Token  string `json:"token"`
}

type AuthService struct {
	Config             infrastructure.Config
	EncryptionProvider infrastructure.EncryptionProvider
	UserService
}

func (as *AuthService) SignUp(ctx context.Context, params CreateUserParams) (*SignUpResponse, error) {
	id, err := as.UserService.Create(ctx, params)

	if err != nil {
		return nil, err
	}

	payload := TokenPayload{Id: id}

	var payloadMap map[string]interface{}

	mapstructure.Decode(payload, &payloadMap)

	token, err := as.EncryptionProvider.CreateJwt(payloadMap, as.Config.JwtSecret)

	if err != nil {
		return nil, err
	}

	return &SignUpResponse{UserId: id, Token: token}, nil
}

func (as *AuthService) Authenticate(ctx context.Context, token string) (*User, error) {
	payloadMap, err := as.EncryptionProvider.VerifyJwt(token, as.Config.JwtSecret)

	if err != nil {
		return nil, infrastructure.ApiError{Code: http.StatusUnauthorized, Details: "Malformed token"}
	}

	var payload TokenPayload

	mapstructure.Decode(payloadMap, &payload)

	user, err := as.UserService.GetById(ctx, payload.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
