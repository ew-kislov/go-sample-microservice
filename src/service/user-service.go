package service

import (
	"context"
	"go-sample-microservice/src/infrastructure"
	"go-sample-microservice/src/repository"
	"net/http"
)

type User struct {
	Id          int
	Email       string
	Username    string
	DisplayName string
}

type UserService struct {
	UserRepository     repository.UserRepository
	EncryptionProvider infrastructure.EncryptionProvider
}

func (us *UserService) Create(ctx context.Context, params CreateUserParams) (int64, error) {
	salt, err := us.EncryptionProvider.GenerateSalt(16)

	if err != nil {
		return 0, infrastructure.ApiError{Code: http.StatusInternalServerError, Details: err.Error()}
	}

	hash := us.EncryptionProvider.GenerateHash(params.Password, salt)

	id, err := us.UserRepository.Create(
		ctx,
		repository.CreateUserParams{
			Email:       params.Email,
			DisplayName: params.DisplayName,
			Username:    params.Username,
			Salt:        salt,
			Hash:        hash,
		},
	)

	if databaseError, ok := err.(infrastructure.DatabaseError); ok && databaseError.Type == infrastructure.DuplicateError {
		return 0, infrastructure.ApiError{Code: http.StatusConflict, Message: "User with provided email or username already exists"}
	} else if err != nil {
		return 0, infrastructure.ApiError{Code: http.StatusInternalServerError, Details: err.Error()}
	}

	return id, nil
}

func (us *UserService) GetById(ctx context.Context, id int64) (*User, error) {
	userRaw, err := us.UserRepository.GetById(ctx, id)

	if databaseError, ok := err.(infrastructure.DatabaseError); ok && databaseError.Type == infrastructure.NotFound {
		return nil, infrastructure.ApiError{Code: http.StatusNotFound, Message: "User with provided id not found"}
	} else if err != nil {
		return nil, infrastructure.ApiError{Code: http.StatusInternalServerError, Details: err.Error()}
	}

	user := &User{
		Id:          userRaw.Id,
		Email:       userRaw.Email,
		Username:    userRaw.Username,
		DisplayName: userRaw.DisplayName,
	}

	return user, nil
}
