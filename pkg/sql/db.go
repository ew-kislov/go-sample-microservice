package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func CreateDatabase(config *cfg.Config, logger *logrus.Logger) Database {
	return &database{
		db:     connectDb(config),
		logger: logger,
	}
}

func (d *database) Query(ctx context.Context, query string, params ...any) (QueryResult, error) {
	start := time.Now()

	d.logger.WithContext(ctx).WithFields(logrus.Fields{"args": params}).Infof("Query started: `%s`.", query)

	result, err := d.queryInternal(ctx, query, params...)

	fields := logrus.Fields{"args": params, "result": result, "durationMs": time.Since(start).Milliseconds()}

	if err != nil {
		d.logger.WithContext(ctx).WithFields(fields).Infof(
			"Query failed with error. Query: `%s`. Error: %s", query, err.Error(),
		)
	} else {
		d.logger.WithContext(ctx).WithFields(fields).Infof("Query finished: `%s`.", query)
	}

	return result, err
}

func (d *database) Exec(ctx context.Context, query string, params ...any) (*ExecResult, error) {
	start := time.Now()

	d.logger.WithContext(ctx).WithFields(logrus.Fields{"args": params}).Infof("Query started: `%s`.", query)

	result, err := d.execInternal(ctx, query, params)

	fields := logrus.Fields{"args": params, "result": result, "durationMs": time.Since(start).Milliseconds()}

	if err != nil {
		d.logger.WithContext(ctx).WithFields(fields).Infof(
			"Query failed with error. Query: `%s`. Error: %s", query, err.Error(),
		)
	} else {
		d.logger.WithContext(ctx).WithFields(fields).Infof("Query finished: `%s`.", query)
	}

	return result, err
}

func (d *database) Close() error {
	return d.db.Close()
}

func (d *database) queryInternal(_ context.Context, query string, params ...any) (QueryResult, error) {
	rows, err := d.db.Query(query, params...)

	if err == nil {
		return d.mapRowsToMap(rows)
	} else {
		return nil, d.handleQueryError(err)
	}
}

func (d *database) execInternal(_ context.Context, query string, params ...any) (*ExecResult, error) {
	result, err := d.db.Exec(query, params)

	if err != nil {
		return nil, d.handleQueryError(err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return nil, fmt.Errorf("error getting affected rows: %w", err)
	}

	return &ExecResult{RowsAffected: rowsAffected}, nil
}

func (*database) handleQueryError(err error) error {
	pgErr, ok := err.(*pq.Error)

	if !ok {
		return err
	}

	switch pgErr.Code {
	case "23505":
		return Error{Type: DuplicateError, Message: err.Error()}
	}

	return err
}

func (*database) mapRowsToMap(rows *sql.Rows) ([]map[string]any, error) {
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %w", err)
	}

	results := []map[string]any{}

	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		rowMap := make(map[string]any)
		for i, col := range columns {
			rowMap[col] = values[i]
		}

		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

func connectDb(config *cfg.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"user=%s dbname=%s sslmode=disable password=%s host=%s",
		config.DatabaseUser, config.DatabaseName, config.DatabasePassword, config.DatabaseHost,
	)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		panic(fmt.Errorf("error connecting database: %w", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("error while ping database: %w", err))
	}

	return db
}
