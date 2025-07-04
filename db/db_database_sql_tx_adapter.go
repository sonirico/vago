package db

import (
	"context"
	"database/sql"
)

type txAdapter struct {
	tx *sql.Tx
}

func (s *txAdapter) Exec(query string, args ...any) (Result, error) {
	return s.tx.Exec(query, args...)
}

func (s *txAdapter) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	return s.tx.ExecContext(ctx, query, args...)
}

func (s *txAdapter) Query(query string, args ...any) (Rows, error) {
	return s.tx.Query(query, args...)
}

func (s *txAdapter) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	return s.tx.QueryContext(ctx, query, args...)
}

func (s *txAdapter) QueryRow(query string, args ...any) Row {
	return s.tx.QueryRow(query, args...)
}

func (s *txAdapter) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	return s.tx.QueryRowContext(ctx, query, args...)
}

func (s *txAdapter) Commit(_ context.Context) error {
	return s.tx.Commit()
}

func (s *txAdapter) Rollback(_ context.Context) error {
	return s.tx.Rollback()
}

func newSqlTxAdapter(tx *sql.Tx) *txAdapter {
	return &txAdapter{tx: tx}
}
