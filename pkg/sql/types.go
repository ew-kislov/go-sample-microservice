package sql

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

type ErrorType int8

const (
	DuplicateError ErrorType = iota
	NotFound
	MappingError
)

type Error struct {
	Type    ErrorType
	Message string
}

func (err Error) Error() string {
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
