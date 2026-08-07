package db

import (
	"context"
	"errors"
)

// ErrNoopQuerier is returned by every method of the noop Querier backing
// NewNoopRepoContext and NewNoopRepoContextTx. There is no underlying
// connection to report a more specific failure from.
var ErrNoopQuerier = errors.New("db: noop querier has no underlying connection")

// noopQuerier is a Querier that performs no I/O and reports ErrNoopQuerier
// on every operation, instead of touching a zero-value driver object.
type noopQuerier struct{}

func newNoopQuerier() *noopQuerier {
	return &noopQuerier{}
}

func (n *noopQuerier) Exec(query string, args ...any) (Result, error) {
	return nil, ErrNoopQuerier
}

func (n *noopQuerier) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	return nil, ErrNoopQuerier
}

func (n *noopQuerier) Query(query string, args ...any) (Rows, error) {
	return nil, ErrNoopQuerier
}

func (n *noopQuerier) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	return nil, ErrNoopQuerier
}

func (n *noopQuerier) QueryRow(query string, args ...any) Row {
	return newNoopRow()
}

func (n *noopQuerier) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	return newNoopRow()
}

// noopRow is a Row that reports ErrNoopQuerier from both Err and Scan,
// mirroring how database/sql defers a QueryRow failure until Scan is called.
type noopRow struct{}

func newNoopRow() *noopRow {
	return &noopRow{}
}

func (n *noopRow) Scan(dest ...any) error {
	return ErrNoopQuerier
}

func (n *noopRow) Err() error {
	return ErrNoopQuerier
}
