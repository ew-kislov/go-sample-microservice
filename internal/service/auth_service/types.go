package authservice

import (
	"context"

	userrepository "github.com/ew-kislov/go-sample-microservice/internal/repository/user_repository"
	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"
)

type TokenPayload struct {
	Id int64 `mapstructure:"id"`
}

type SignUpResponse struct {
	UserId int64
	Token  string
}

type SignUpParams struct {
	Email       string
	Username    string
	DisplayName string
	Password    string
}

type User struct {
	Id          int
	Email       string
	Username    string
	DisplayName string
}

type AuthService interface {
	SignUp(ctx context.Context, params SignUpParams) (*SignUpResponse, error)
	Authenticate(ctx context.Context, token string) (*User, error)
}

type authService struct {
	config         cfg.Config
	userRepository userrepository.UserRepository
}
