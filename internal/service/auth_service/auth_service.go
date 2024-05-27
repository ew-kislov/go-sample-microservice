package authservice

import (
	"context"
	"errors"
	"net/http"

	userrepository "github.com/ew-kislov/go-sample-microservice/internal/repository/user_repository"
	"github.com/ew-kislov/go-sample-microservice/pkg/api"
	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"
	"github.com/ew-kislov/go-sample-microservice/pkg/encryption"
	"github.com/ew-kislov/go-sample-microservice/pkg/jwt"
	"github.com/ew-kislov/go-sample-microservice/pkg/sql"

	"github.com/mitchellh/mapstructure"
)

func NewAuthService(
	config *cfg.Config,
	userRepository userrepository.UserRepository,
) AuthService {
	return &authService{config, userRepository}
}

func (as *authService) SignUp(ctx context.Context, params SignUpParams) (*SignUpResponse, error) {
	salt, err := encryption.GenerateSalt(16)

	if err != nil {
		return nil, err
	}

	hash := encryption.GenerateHash(params.Password, salt)

	id, err := as.userRepository.Create(
		ctx,
		&userrepository.CreateUserParams{
			Email:       params.Email,
			DisplayName: params.DisplayName,
			Username:    params.Username,
			Salt:        salt,
			Hash:        hash,
		},
	)

	if databaseError, ok := err.(sql.DatabaseError); ok && databaseError.Type == sql.DuplicateError {
		return nil, api.Error{
			Code:    http.StatusConflict,
			Message: "User with provided email or username already exists",
		}
	} else if err != nil {
		return nil, err
	}

	payload := TokenPayload{Id: id}

	var payloadMap map[string]any

	err = mapstructure.Decode(payload, &payloadMap)

	if err != nil {
		return nil, errors.New("Could not decode map TokenPayload map")
	}

	token, err := jwt.CreateJwt(payloadMap, as.config.JwtSecret)

	if err != nil {
		return nil, err
	}

	return &SignUpResponse{UserId: id, Token: token}, nil
}

func (as *authService) Authenticate(ctx context.Context, token string) (*User, error) {
	payloadMap, err := jwt.VerifyJwt(token, as.config.JwtSecret)

	if err != nil {
		return nil, api.Error{Code: http.StatusUnauthorized, Message: "Malformed token"}
	}

	var payload TokenPayload

	err = mapstructure.Decode(payloadMap, &payload)

	if err != nil {
		return nil, errors.New("Could not decode map into TokenPayload")
	}

	userRaw, err := as.userRepository.GetById(ctx, payload.Id)

	if databaseError, ok := err.(sql.DatabaseError); ok && databaseError.Type == sql.NotFound {
		return nil, api.Error{Code: http.StatusNotFound, Message: "User with provided id not found"}
	} else if err != nil {
		return nil, err
	}

	user := &User{
		Id:          userRaw.Id,
		Email:       userRaw.Email,
		Username:    userRaw.Username,
		DisplayName: userRaw.DisplayName,
	}

	return user, nil
}
