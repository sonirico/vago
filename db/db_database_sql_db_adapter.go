package db

import (
	"context"
	"database/sql"
)

type sqlAdapter struct {
	db *sql.DB
}

func (s *sqlAdapter) Close() error {
	return s.db.Close()
}

func (s *sqlAdapter) Exec(query string, args ...any) (Result, error) {
	return s.db.Exec(query, args...)
}

func (s *sqlAdapter) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *sqlAdapter) Query(query string, args ...any) (Rows, error) {
	return s.db.Query(query, args...)
}

func (s *sqlAdapter) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *sqlAdapter) QueryRow(query string, args ...any) Row {
	return s.db.QueryRow(query, args...)
}

func (s *sqlAdapter) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *sqlAdapter) Begin(_ context.Context) (Tx, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	return newSqlTxAdapter(tx), nil
}

func newSqlAdapter(db *sql.DB) *sqlAdapter {
	return &sqlAdapter{db: db}
}
