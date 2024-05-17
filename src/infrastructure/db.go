package infrastructure

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DatabaseErrorType int8

const (
	DuplicateError DatabaseErrorType = iota
	NotFound
	MappingError
)

type DatabaseError struct {
	Type    DatabaseErrorType
	Details string
}

func (err DatabaseError) Error() string {
	return err.Details
}

func CreateDatabase(config Config) *sqlx.DB {
	connectionString := fmt.Sprintf(
		"user=%s dbname=%s sslmode=disable password=%s host=%s",
		config.DatabaseUser, config.DatabaseName, config.DatabasePassword, config.DatabaseHost,
	)
	db, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		panic(fmt.Errorf("got error while connecting to database: %w", err))
	}

	return db
}
