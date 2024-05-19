package userrepository

import (
	"context"

	"github.com/ew-kislov/go-sample-microservice/pkg"
)

func (ur *userRepository) Create(ctx context.Context, params CreateUserParams) (int64, error) {
	return ur.BaseRepository.Insert(
		ctx,
		pkg.InsertParams{
			Table:   "users",
			Columns: []string{"email", "username", "display_name", "salt", "hash"},
			Data:    params,
		},
	)
}

func (ur *userRepository) GetById(ctx context.Context, id int64) (*User, error) {
	user := User{}

	err := ur.BaseRepository.SelectById(ctx, pkg.SelectByIdParams{Table: "users", Id: id}, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
