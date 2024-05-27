package userrepository

import (
	"context"
	"errors"

	"github.com/mitchellh/mapstructure"
)

func (ur *userRepository) Create(ctx context.Context, params *CreateUserParams) (int64, error) {
	result, err := ur.db.Query(
		ctx,
		"INSERT INTO users(email, username, display_name, salt, hash) values ($1, $2, $3, $4, $5) RETURNING id",
		params.Email, params.Username, params.DisplayName, params.Salt, params.Hash,
	)

	if err != nil {
		return 0, err
	}

	id, ok := result[0]["id"].(int64)

	if !ok {
		return 0, errors.New("could not get created id from result set")
	}

	return id, err
}

func (ur *userRepository) GetById(ctx context.Context, id int64) (*User, error) {
	result, err := ur.db.Query(ctx, "SELECT * from users WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	var user User

	err = mapstructure.Decode(result[0], &user)

	if err != nil {
		return nil, errors.New("could not decode result set to User type")
	}

	return &user, nil
}
