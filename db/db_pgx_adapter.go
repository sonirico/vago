package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxAdapter struct {
	db *pgxpool.Pool
}

func newPgxAdapter(db *pgxpool.Pool) *pgxAdapter {
	return &pgxAdapter{db: db}
}

func (p *pgxAdapter) Exec(query string, args ...any) (Result, error) {
	return p.ExecContext(context.Background(), query, args...)
}

func (p *pgxAdapter) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	tag, err := p.db.Exec(ctx, query, args...)
	return &pgxResult{tag}, err
}

func (p *pgxAdapter) Query(query string, args ...any) (Rows, error) {
	return p.QueryContext(context.Background(), query, args...)
}

func (p *pgxAdapter) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	rows, err := p.db.Query(ctx, query, args...)
	return &pgxRows{rows}, err
}

func (p *pgxAdapter) QueryRow(query string, args ...any) Row {
	return p.QueryRowContext(context.Background(), query, args...)
}

func (p *pgxAdapter) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	return &pgxRow{p.db.QueryRow(ctx, query, args...)}
}

func (p *pgxAdapter) Begin(ctx context.Context) (Tx, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &pgxTxAdapter{tx}, nil
}

func (p *pgxAdapter) Close() error {
	p.db.Close()
	return nil
}

// pgxResult wraps pgx.CommandTag to implement Result interface
type pgxResult struct {
	tag pgconn.CommandTag
}

func (r *pgxResult) LastInsertId() (int64, error) {
	return 0, nil // pgx doesn't support LastInsertId
}

func (r *pgxResult) RowsAffected() (int64, error) {
	return r.tag.RowsAffected(), nil
}

// pgxRows wraps pgx.Rows to implement Rows interface
type pgxRows struct {
	rows pgx.Rows
}

func (r *pgxRows) Err() error {
	return r.rows.Err()
}

func (r *pgxRows) Next() bool {
	return r.rows.Next()
}

func (r *pgxRows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *pgxRows) Close() error {
	r.rows.Close()
	return nil
}

func (r *pgxRows) Columns() ([]string, error) {
	fields := r.rows.FieldDescriptions()
	columns := make([]string, len(fields))
	for i, field := range fields {
		columns[i] = string(field.Name)
	}
	return columns, nil
}

// pgxRow wraps pgx.Row to implement Row interface
type pgxRow struct {
	row pgx.Row
}

func (r *pgxRow) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

func (r *pgxRow) Err() error {
	return nil
}

// pgxTxAdapter implements Tx interface using pgx.Tx
type pgxTxAdapter struct {
	tx pgx.Tx
}

func (t *pgxTxAdapter) Exec(query string, args ...any) (Result, error) {
	return t.ExecContext(context.Background(), query, args...)
}

func (t *pgxTxAdapter) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	tag, err := t.tx.Exec(ctx, query, args...)
	return &pgxResult{tag}, err
}

func (t *pgxTxAdapter) Query(query string, args ...any) (Rows, error) {
	return t.QueryContext(context.Background(), query, args...)
}

func (t *pgxTxAdapter) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	rows, err := t.tx.Query(ctx, query, args...)
	return &pgxRows{rows}, err
}

func (t *pgxTxAdapter) QueryRow(query string, args ...any) Row {
	return t.QueryRowContext(context.Background(), query, args...)
}

func (t *pgxTxAdapter) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	return &pgxRow{t.tx.QueryRow(ctx, query, args...)}
}

func (t *pgxTxAdapter) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *pgxTxAdapter) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}
