package pkg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type BaseRepository struct {
	Logger *logrus.Logger
	Db     *sqlx.DB
}

type InsertParams struct {
	Table   string
	Columns []string
	Data    interface{}
}

type SelectByIdParams struct {
	Table string
	Id    int64
}

func (repository *BaseRepository) Insert(ctx context.Context, params InsertParams) (int64, error) {
	query := fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES (%s) RETURNING id",
		params.Table,
		strings.Join(params.Columns, ", "),
		toNamedParams(params.Columns),
	)

	start := time.Now()

	repository.Logger.WithContext(ctx).WithFields(logrus.Fields{"args": params.Data}).Infof("Query started: `%s`.", query)

	statement, err := repository.Db.PrepareNamedContext(ctx, query)

	if err != nil {
		return 0, err
	}

	var id int64

	err = statement.Get(&id, params.Data)

	fields := logrus.Fields{"args": params.Data, "durationMs": time.Since(start).Milliseconds()}

	if err != nil {
		repository.Logger.WithContext(ctx).WithFields(fields).Infof("Query failed with error. Query: `%s`. Error: %s", query, err.Error())
	} else {
		repository.Logger.WithContext(ctx).WithFields(fields).Infof("Query finished: `%s`.", query)
	}

	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
		return 0, DatabaseError{Type: DuplicateError, Message: err.Error()}
	}

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repository *BaseRepository) SelectById(ctx context.Context, params SelectByIdParams, result interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", params.Table)

	start := time.Now()

	startFields := logrus.Fields{"args": logrus.Fields{"id": params.Id}}
	repository.Logger.WithContext(ctx).WithFields(startFields).Infof("Query started: `%s`.", query)

	err := repository.Db.GetContext(ctx, result, fmt.Sprintf("SELECT * FROM %s WHERE id = $1", params.Table), params.Id)

	finishFields := logrus.Fields{"args": logrus.Fields{"id": params.Id}, "durationMs": time.Since(start).Milliseconds()}

	if err != nil {
		repository.Logger.WithContext(ctx).WithFields(finishFields).Infof("Query failed with error. Query: `%s`. Error: %s", query, err.Error())
		return err
	} else {
		repository.Logger.WithContext(ctx).WithFields(finishFields).Infof("Query finished: `%s`.", query)
	}

	if result == nil {
		return DatabaseError{Type: NotFound, Message: "Record not found"}
	}

	return nil
}

func toNamedParams(columns []string) string {
	params := make([]string, len(columns))
	for i, v := range columns {
		params[i] = ":" + v
	}
	return strings.Join(params, ", ")
}
