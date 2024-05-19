package userservice

import (
	"context"

	userrepository "github.com/ew-kislov/go-sample-microservice/internal/repository/user_repository"
	"github.com/ew-kislov/go-sample-microservice/pkg"
)

type CreateUserParams struct {
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

type UserService interface {
	Create(ctx context.Context, params CreateUserParams) (int64, error)
	GetById(ctx context.Context, id int64) (*User, error)
}

type userService struct {
	UserRepository     userrepository.UserRepository
	EncryptionProvider pkg.EncryptionProvider
}

func NewUserService(
	userRepository userrepository.UserRepository,
	encryptionProvider pkg.EncryptionProvider,
) UserService {
	return &userService{UserRepository: userRepository, EncryptionProvider: encryptionProvider}
}
