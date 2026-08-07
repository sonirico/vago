package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type executorRWOTestFixture struct {
	executor ExecutorRWO
	roMock   sqlmock.Sqlmock
	rwMock   sqlmock.Sqlmock
}

func newTestExecutorRWOFixture(t *testing.T) executorRWOTestFixture {
	t.Helper()

	roConn, roMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	t.Cleanup(func() {
		roMock.ExpectClose()
		require.NoError(t, roConn.Close())
	})

	rwConn, rwMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	t.Cleanup(func() {
		rwMock.ExpectClose()
		require.NoError(t, rwConn.Close())
	})

	return executorRWOTestFixture{
		executor: NewRWOExecutor(log, newSqlAdapter(roConn), newSqlAdapter(rwConn)),
		roMock:   roMock,
		rwMock:   rwMock,
	}
}

// doer is the shared shape of ExecutorRO and ExecutorRW needed to run a plain query, letting a
// single table drive both RO() and RW() through the same assertions.
type doer interface {
	Do(ctx context.Context, fn func(ctx Context) error) error
}

// TestExecutorRWOAccessors is the regression test for the defect where RW() was wired to the
// ro Handler instead of rw. Two distinct sqlmock handles are the artifact: the expectation is set
// on the handle the accessor is supposed to use, and ExpectationsWereMet on the OTHER handle
// (which never sees any expectation) proves the query did not silently land there instead.
func TestExecutorRWOAccessors(t *testing.T) {
	type testCase struct {
		name       string
		invoke     func(fx executorRWOTestFixture) doer
		wantMock   func(fx executorRWOTestFixture) sqlmock.Sqlmock
		unusedMock func(fx executorRWOTestFixture) sqlmock.Sqlmock
	}

	tests := []testCase{
		{
			name:       "RO uses the ro handler",
			invoke:     func(fx executorRWOTestFixture) doer { return fx.executor.RO() },
			wantMock:   func(fx executorRWOTestFixture) sqlmock.Sqlmock { return fx.roMock },
			unusedMock: func(fx executorRWOTestFixture) sqlmock.Sqlmock { return fx.rwMock },
		},
		{
			name:       "RW uses the rw handler",
			invoke:     func(fx executorRWOTestFixture) doer { return fx.executor.RW() },
			wantMock:   func(fx executorRWOTestFixture) sqlmock.Sqlmock { return fx.rwMock },
			unusedMock: func(fx executorRWOTestFixture) sqlmock.Sqlmock { return fx.roMock },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fx := newTestExecutorRWOFixture(t)
			tc.wantMock(fx).ExpectQuery("SELECT 1;").
				WillReturnRows(sqlmock.NewRows([]string{"n"}).AddRow(1))

			executor := tc.invoke(fx)
			err := executor.Do(context.Background(), func(ctx Context) error {
				var n int
				return ctx.Querier().QueryRowContext(ctx, "SELECT 1;").Scan(&n)
			})

			require.NoError(t, err)
			assert.NoError(t, tc.wantMock(fx).ExpectationsWereMet())
			assert.NoError(t, tc.unusedMock(fx).ExpectationsWereMet())
		})
	}
}
