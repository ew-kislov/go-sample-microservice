package userservice

import (
	"context"
	"net/http"

	userrepository "github.com/ew-kislov/go-sample-microservice/internal/repository/user_repository"
	"github.com/ew-kislov/go-sample-microservice/pkg"
	"github.com/ew-kislov/go-sample-microservice/pkg/db"
)

func (us *userService) Create(ctx context.Context, params CreateUserParams) (int64, error) {
	salt, err := us.EncryptionProvider.GenerateSalt(16)

	if err != nil {
		return 0, pkg.ApiError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	hash := us.EncryptionProvider.GenerateHash(params.Password, salt)

	id, err := us.UserRepository.Create(
		ctx,
		userrepository.CreateUserParams{
			Email:       params.Email,
			DisplayName: params.DisplayName,
			Username:    params.Username,
			Salt:        salt,
			Hash:        hash,
		},
	)

	if databaseError, ok := err.(db.DatabaseError); ok && databaseError.Type == db.DuplicateError {
		return 0, pkg.ApiError{Code: http.StatusConflict, Message: "User with provided email or username already exists"}
	} else if err != nil {
		return 0, pkg.ApiError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return id, nil
}

func (us *userService) GetById(ctx context.Context, id int64) (*User, error) {
	userRaw, err := us.UserRepository.GetById(ctx, id)

	if databaseError, ok := err.(db.DatabaseError); ok && databaseError.Type == db.NotFound {
		return nil, pkg.ApiError{Code: http.StatusNotFound, Message: "User with provided id not found"}
	} else if err != nil {
		return nil, pkg.ApiError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	user := &User{
		Id:          userRaw.Id,
		Email:       userRaw.Email,
		Username:    userRaw.Username,
		DisplayName: userRaw.DisplayName,
	}

	return user, nil
}
