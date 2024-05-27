package sql

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

type DatabaseErrorType int8

const (
	DuplicateError DatabaseErrorType = iota
	NotFound
	MappingError
)

type DatabaseError struct {
	Type    DatabaseErrorType
	Message string
}

func (err DatabaseError) Error() string {
	return err.Message
}

type database struct {
	db     *sql.DB
	logger *logrus.Logger
}

type ExecResult struct {
	RowsAffected int64
}

type QueryResult []map[string]any

type Database interface {
	Query(ctx context.Context, query string, params ...any) (QueryResult, error)
	Exec(ctx context.Context, query string, params ...any) (*ExecResult, error)
	Close() error
}
