package repository

import (
	"context"
	"go-sample-microservice/src/infrastructure"
	"time"

	"github.com/jmoiron/sqlx"
)

type CreateUserParams struct {
	Email       string `db:"email"`
	Username    string `db:"username"`
	DisplayName string `db:"display_name"`
	Salt        string `db:"salt"`
	Hash        string `db:"hash"`
}

type User struct {
	Id          int       `db:"id"`
	Email       string    `db:"email"`
	Username    string    `db:"username"`
	DisplayName string    `db:"display_name"`
	Salt        string    `db:"salt"`
	Hash        string    `db:"hash"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type UserRepository struct {
	Db             *sqlx.DB
	BaseRepository infrastructure.BaseRepository
}

func (ur *UserRepository) Create(ctx context.Context, params CreateUserParams) (int64, error) {
	return ur.BaseRepository.Insert(
		ctx,
		infrastructure.InsertParams{
			Table:   "users",
			Columns: []string{"email", "username", "display_name", "salt", "hash"},
			Data:    params,
		},
	)
}

func (ur *UserRepository) GetById(ctx context.Context, id int64) (*User, error) {
	user := User{}

	err := ur.BaseRepository.SelectById(ctx, infrastructure.SelectByIdParams{Table: "users", Id: id}, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
