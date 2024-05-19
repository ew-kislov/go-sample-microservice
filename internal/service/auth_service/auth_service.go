package authservice

import (
	"context"
	"net/http"

	userservice "github.com/ew-kislov/go-sample-microservice/internal/service/user_service"

	"github.com/ew-kislov/go-sample-microservice/pkg"
	"github.com/mitchellh/mapstructure"
)

func NewAuthService(
	config pkg.Config,
	encryptionProvider pkg.EncryptionProvider,
	userService userservice.UserService,
) AuthService {
	return &authService{config, encryptionProvider, userService}
}

func (as *authService) SignUp(ctx context.Context, params userservice.CreateUserParams) (*SignUpResponse, error) {
	id, err := as.userService.Create(ctx, userservice.CreateUserParams(params))

	if err != nil {
		return nil, err
	}

	payload := TokenPayload{Id: id}

	var payloadMap map[string]interface{}

	mapstructure.Decode(payload, &payloadMap)

	token, err := as.encryptionProvider.CreateJwt(payloadMap, as.config.JwtSecret)

	if err != nil {
		return nil, err
	}

	return &SignUpResponse{UserId: id, Token: token}, nil
}

func (as *authService) Authenticate(ctx context.Context, token string) (*userservice.User, error) {
	payloadMap, err := as.encryptionProvider.VerifyJwt(token, as.config.JwtSecret)

	if err != nil {
		return nil, pkg.ApiError{Code: http.StatusUnauthorized, Message: "Malformed token"}
	}

	var payload TokenPayload

	mapstructure.Decode(payloadMap, &payload)

	user, err := as.userService.GetById(ctx, payload.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
