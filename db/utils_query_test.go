package db

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type queryTestFixture struct {
	executor Executor
	mock     sqlmock.Sqlmock
}

func newTestQueryFixture(t *testing.T) queryTestFixture {
	t.Helper()

	conn, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	t.Cleanup(func() {
		mock.ExpectClose()
		require.NoError(t, conn.Close())
	})

	return queryTestFixture{
		executor: newDatabaseSqlExecutor(log, conn),
		mock:     mock,
	}
}

func TestQuery(t *testing.T) {
	errFn := errors.New("fn failed")

	type testCase struct {
		name     string
		setup    func(mock sqlmock.Sqlmock)
		fn       func(ctx Context) (int, error)
		wantData int
		wantErr  error
	}

	tests := []testCase{
		{
			name: "returns fn's data when fn succeeds",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 1;").
					WillReturnRows(sqlmock.NewRows([]string{"n"}).AddRow(1))
			},
			fn: func(ctx Context) (int, error) {
				var n int
				err := ctx.Querier().QueryRowContext(ctx, "SELECT 1;").Scan(&n)
				return n, err
			},
			wantData: 1,
			wantErr:  nil,
		},
		{
			name:  "propagates fn's error instead of discarding it",
			setup: func(mock sqlmock.Sqlmock) {},
			fn: func(ctx Context) (int, error) {
				return 0, errFn
			},
			wantData: 0,
			wantErr:  errFn,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fx := newTestQueryFixture(t)
			tc.setup(fx.mock)

			data, err := Query(context.Background(), fx.executor, tc.fn)

			assert.Equal(t, tc.wantData, data)
			if tc.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.wantErr)
			}
			assert.NoError(t, fx.mock.ExpectationsWereMet())
		})
	}
}

// TestQueryAfterCommitHook is the regression test for the corollary of the same defect in Do-based
// queries: because the callback handed to Do used to always return nil, Do's after-commit hooks ran
// even when fn failed. Propagating fn's error into the callback makes Do skip them, same as a
// failed transaction skips commit.
func TestQueryAfterCommitHook(t *testing.T) {
	errFn := errors.New("fn failed")

	type testCase struct {
		name        string
		fn          func(ctx Context) (int, error)
		wantErr     error
		wantHookRan bool
	}

	tests := []testCase{
		{
			name: "runs the after-commit hook when fn succeeds",
			fn: func(ctx Context) (int, error) {
				return 1, nil
			},
			wantErr:     nil,
			wantHookRan: true,
		},
		{
			name: "skips the after-commit hook when fn fails",
			fn: func(ctx Context) (int, error) {
				return 0, errFn
			},
			wantErr:     errFn,
			wantHookRan: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fx := newTestQueryFixture(t)
			var hookRan bool

			_, err := Query(context.Background(), fx.executor, func(ctx Context) (int, error) {
				ctx.AfterCommitDo(func(Context) { hookRan = true })
				return tc.fn(ctx)
			})

			if tc.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.wantHookRan, hookRan)
			assert.NoError(t, fx.mock.ExpectationsWereMet())
		})
	}
}

func TestQueryRO(t *testing.T) {
	errFn := errors.New("fn failed")

	type testCase struct {
		name     string
		setup    func(mock sqlmock.Sqlmock)
		fn       func(ctx Context) (int, error)
		wantData int
		wantErr  error
	}

	tests := []testCase{
		{
			name: "returns fn's data when fn succeeds",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 1;").
					WillReturnRows(sqlmock.NewRows([]string{"n"}).AddRow(1))
			},
			fn: func(ctx Context) (int, error) {
				var n int
				err := ctx.Querier().QueryRowContext(ctx, "SELECT 1;").Scan(&n)
				return n, err
			},
			wantData: 1,
			wantErr:  nil,
		},
		{
			name:  "propagates fn's error instead of discarding it",
			setup: func(mock sqlmock.Sqlmock) {},
			fn: func(ctx Context) (int, error) {
				return 0, errFn
			},
			wantData: 0,
			wantErr:  errFn,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fx := newTestQueryFixture(t)
			tc.setup(fx.mock)

			data, err := QueryRO(context.Background(), fx.executor, tc.fn)

			assert.Equal(t, tc.wantData, data)
			if tc.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.wantErr)
			}
			assert.NoError(t, fx.mock.ExpectationsWereMet())
		})
	}
}

func TestQueryRW(t *testing.T) {
	errFn := errors.New("fn failed")

	type testCase struct {
		name     string
		setup    func(mock sqlmock.Sqlmock)
		fn       func(ctx Context) (int, error)
		wantData int
		wantErr  error
	}

	tests := []testCase{
		{
			name: "returns fn's data when fn succeeds",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 1;").
					WillReturnRows(sqlmock.NewRows([]string{"n"}).AddRow(1))
			},
			fn: func(ctx Context) (int, error) {
				var n int
				err := ctx.Querier().QueryRowContext(ctx, "SELECT 1;").Scan(&n)
				return n, err
			},
			wantData: 1,
			wantErr:  nil,
		},
		{
			name:  "propagates fn's error instead of discarding it",
			setup: func(mock sqlmock.Sqlmock) {},
			fn: func(ctx Context) (int, error) {
				return 0, errFn
			},
			wantData: 0,
			wantErr:  errFn,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fx := newTestQueryFixture(t)
			tc.setup(fx.mock)

			data, err := QueryRW(context.Background(), fx.executor, tc.fn)

			assert.Equal(t, tc.wantData, data)
			if tc.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.wantErr)
			}
			assert.NoError(t, fx.mock.ExpectationsWereMet())
		})
	}
}

// TestQueryTx is the regression test for the defect where QueryTx always returned nil from the
// DoWithTx callback, so a failing fn still committed instead of rolling back, and Begin/Commit
// errors from the executor were silently discarded. sqlmock's ordered expectations are the
// artifact: if the code under test calls Commit when only Rollback is expected (or vice versa),
// ExpectationsWereMet reports the mismatch - a test that only inspected the returned error would
// have passed against the broken implementation, because the broken version still returned fn's
// error even while committing.
func TestQueryTx(t *testing.T) {
	errFn := errors.New("fn failed")
	errCommit := errors.New("commit failed")
	errRollback := errors.New("rollback failed")

	type testCase struct {
		name     string
		setup    func(mock sqlmock.Sqlmock)
		fn       func(ctx Context) (int, error)
		wantData int
		wantErrs []error
	}

	tests := []testCase{
		{
			name: "commits when fn succeeds",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT 1;").
					WillReturnRows(sqlmock.NewRows([]string{"n"}).AddRow(1))
				mock.ExpectCommit()
			},
			fn: func(ctx Context) (int, error) {
				var n int
				err := ctx.Querier().QueryRowContext(ctx, "SELECT 1;").Scan(&n)
				return n, err
			},
			wantData: 1,
			wantErrs: nil,
		},
		{
			name: "rolls back instead of committing when fn fails",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			fn: func(ctx Context) (int, error) {
				return 0, errFn
			},
			wantData: 0,
			wantErrs: []error{errFn},
		},
		{
			name: "surfaces the commit error when fn succeeds but commit fails",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT 1;").
					WillReturnRows(sqlmock.NewRows([]string{"n"}).AddRow(1))
				mock.ExpectCommit().WillReturnError(errCommit)
			},
			fn: func(ctx Context) (int, error) {
				var n int
				err := ctx.Querier().QueryRowContext(ctx, "SELECT 1;").Scan(&n)
				return n, err
			},
			wantData: 1,
			wantErrs: []error{errCommit},
		},
		{
			name: "joins fn's error with the rollback error instead of dropping one",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectRollback().WillReturnError(errRollback)
			},
			fn: func(ctx Context) (int, error) {
				return 0, errFn
			},
			wantData: 0,
			wantErrs: []error{errFn, errRollback},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fx := newTestQueryFixture(t)
			tc.setup(fx.mock)

			data, err := QueryTx(context.Background(), fx.executor, tc.fn)

			assert.Equal(t, tc.wantData, data)
			if len(tc.wantErrs) == 0 {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				for _, wantErr := range tc.wantErrs {
					assert.ErrorIs(t, err, wantErr)
				}
			}
			assert.NoError(t, fx.mock.ExpectationsWereMet())
		})
	}
}
