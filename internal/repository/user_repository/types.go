package userrepository

import (
	"context"
	"time"

	"github.com/ew-kislov/go-sample-microservice/pkg"

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

type UserRepository interface {
	Create(ctx context.Context, params CreateUserParams) (int64, error)
	GetById(ctx context.Context, id int64) (*User, error)
}

type userRepository struct {
	Db             *sqlx.DB
	BaseRepository pkg.BaseRepository
}

func NewUserRepository(db *sqlx.DB, baseRepository pkg.BaseRepository) UserRepository {
	return &userRepository{Db: db, BaseRepository: baseRepository}
}
