package authservice

import (
	"context"

	userservice "github.com/ew-kislov/go-sample-microservice/internal/service/user_service"

	"github.com/ew-kislov/go-sample-microservice/pkg"
)

type TokenPayload struct {
	Id int64 `mapstructure:"id"`
}

type SignUpResponse struct {
	UserId int64  `json:"userId"`
	Token  string `json:"token"`
}

type AuthService interface {
	SignUp(ctx context.Context, params userservice.CreateUserParams) (*SignUpResponse, error)
	Authenticate(ctx context.Context, token string) (*userservice.User, error)
}

type authService struct {
	config             pkg.Config
	encryptionProvider pkg.EncryptionProvider
	userService        userservice.UserService
}
