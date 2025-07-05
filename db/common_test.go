package db

import (
	"context"
	"strings"
	"testing"
)

type mockTx struct{}

func (m mockTx) Exec(query string, args ...any) (Result, error) { return mockResult{}, nil }
func (m mockTx) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	return mockResult{}, nil
}
func (m mockTx) Query(query string, args ...any) (Rows, error) { return mockRows{}, nil }
func (m mockTx) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	return mockRows{}, nil
}
func (m mockTx) QueryRow(query string, args ...any) Row { return mockRow{} }
func (m mockTx) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	return mockRow{}
}
func (mockTx) Commit(ctx context.Context) error   { return nil }
func (mockTx) Rollback(ctx context.Context) error { return nil }

func TestHandlerInterface(t *testing.T) {
	var _ Handler = (*mockHandler)(nil)
}

type mockHandler struct{}

// Exec implements Handler.
func (m *mockHandler) Exec(query string, args ...any) (Result, error) {
	return mockResult{}, nil
}

// ExecContext implements Handler.
func (m *mockHandler) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	return mockResult{}, nil
}

// Query implements Handler.
func (m *mockHandler) Query(query string, args ...any) (Rows, error) {
	return mockRows{}, nil
}

// QueryContext implements Handler.
func (m *mockHandler) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	return mockRows{}, nil
}

// QueryRow implements Handler.
func (m *mockHandler) QueryRow(query string, args ...any) Row {
	return mockRow{}
}

// QueryRowContext implements Handler.
func (m *mockHandler) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	return mockRow{}
}

func (mockHandler) Querier() Querier                      { return nil }
func (mockHandler) Begin(ctx context.Context) (Tx, error) { return mockTx{}, nil }
func (mockHandler) Close() error                          { return nil }

func TestResultInterface(t *testing.T) {
	var _ Result = (*mockResult)(nil)
}

type mockResult struct{}

func (mockResult) LastInsertId() (int64, error) { return 1, nil }
func (mockResult) RowsAffected() (int64, error) { return 2, nil }

func TestRowsInterface(t *testing.T) {
	var _ Rows = (*mockRows)(nil)
}

type mockRows struct{}

func (mockRows) Err() error                 { return nil }
func (mockRows) Next() bool                 { return false }
func (mockRows) Scan(dst ...any) error      { return nil }
func (mockRows) Close() error               { return nil }
func (mockRows) Columns() ([]string, error) { return []string{"a"}, nil }

func TestRowInterface(t *testing.T) {
	var _ Row = (*mockRow)(nil)
}

type mockRow struct{}

func (mockRow) Scan(...any) error { return nil }
func (mockRow) Err() error        { return nil }

func TestQuerierInterface(t *testing.T) {
	var _ Querier = (*mockQuerier)(nil)
}

type mockQuerier struct{}

func (mockQuerier) Exec(query string, args ...any) (Result, error) { return mockResult{}, nil }
func (mockQuerier) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	return mockResult{}, nil
}
func (mockQuerier) Query(query string, args ...any) (Rows, error) { return mockRows{}, nil }
func (mockQuerier) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	return mockRows{}, nil
}
func (mockQuerier) QueryRow(query string, args ...any) Row { return mockRow{} }
func (mockQuerier) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	return mockRow{}
}

type mockBulkable struct {
	PKVal         []string
	UniqueKeysVal []string
	IncludePKVal  bool
	ColsVal       []string
	RowVal        []any
}

func (m *mockBulkable) PK() []string {
	return m.PKVal
}

func (m *mockBulkable) UniqueKeys() []string {
	return m.UniqueKeysVal
}

func (m *mockBulkable) IncludePKOnUpsert() bool {
	return m.IncludePKVal
}

func (m *mockBulkable) Cols() []string {
	return m.ColsVal
}

func (m *mockBulkable) Row() []any {
	return m.RowVal
}

// normalizeSQL standardizes SQL strings by removing newlines, tabs, and redundant spaces.
func normalizeSQL(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.TrimSpace(s)
	out := make([]rune, 0, len(s))
	lastSpace := false
	for _, r := range s {
		if r == ' ' {
			if !lastSpace {
				out = append(out, r)
			}
			lastSpace = true
		} else {
			out = append(out, r)
			lastSpace = false
		}
	}
	return string(out)
}
