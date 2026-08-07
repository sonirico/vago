package db

import (
	"context"
)

type (
	Handler interface {
		Querier
		Begin(ctx context.Context) (Tx, error)
		Close() error
	}

	Result interface {
		LastInsertId() (int64, error)
		RowsAffected() (int64, error)
	}

	Rows interface {
		Err() error
		Next() bool
		Scan(dst ...any) error
		Close() error
		Columns() ([]string, error)
	}

	// Row is one result row. Err is NOT uniform across backends: the
	// database/sql-backed implementations return real pre-scan errors,
	// while the pgx one always returns nil because pgx.Row has no Err to
	// wrap. Code that must work on both backends has to read the error
	// from Scan; checking Err first and Scan second silently misses
	// failures on pgx alone, which is the hardest kind of difference to
	// notice because it only shows up on one deployment.
	Row interface {
		Scan(...any) error
		Err() error
	}

	Querier interface {
		Exec(query string, args ...any) (Result, error)
		ExecContext(ctx context.Context, query string, args ...any) (Result, error)
		Query(query string, args ...any) (Rows, error)
		QueryContext(ctx context.Context, query string, args ...any) (Rows, error)
		QueryRow(query string, args ...any) Row
		QueryRowContext(ctx context.Context, query string, args ...any) Row
	}

	Tx interface {
		Querier
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
)
