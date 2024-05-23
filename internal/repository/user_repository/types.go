package userrepository

import (
	"context"
	"time"

	"github.com/ew-kislov/go-sample-microservice/pkg/db"
)

type CreateUserParams struct {
	Email       string
	Username    string
	DisplayName string
	Salt        string
	Hash        string
}

type User struct {
	Id          int       `mapstructure:"id"`
	Email       string    `mapstructure:"email"`
	Username    string    `mapstructure:"username"`
	DisplayName string    `mapstructure:"display_name"`
	Salt        string    `mapstructure:"salt"`
	Hash        string    `mapstructure:"hash"`
	CreatedAt   time.Time `mapstructure:"created_at"`
	UpdatedAt   time.Time `mapstructure:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, params CreateUserParams) (int64, error)
	GetById(ctx context.Context, id int64) (*User, error)
}

type userRepository struct {
	db db.Database
}

func NewUserRepository(db db.Database) UserRepository {
	return &userRepository{db}
}
