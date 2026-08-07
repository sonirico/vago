package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNoopRepoContext(t *testing.T) {
	t.Parallel()

	rc := NewNoopRepoContext(context.Background())

	assert.NotPanics(t, func() {
		_, err := rc.Querier().Exec("SELECT 1")
		require.ErrorIs(t, err, ErrNoopQuerier)
	})
}

func TestNewNoopRepoContextTx(t *testing.T) {
	t.Parallel()

	rc := NewNoopRepoContextTx(context.Background())

	assert.NotPanics(t, func() {
		_, err := rc.Querier().Exec("SELECT 1")
		require.ErrorIs(t, err, ErrNoopQuerier)
	})
}

func TestNoopQuerier(t *testing.T) {
	t.Parallel()

	t.Run("Exec returns ErrNoopQuerier", func(t *testing.T) {
		t.Parallel()
		q := newNoopQuerier()
		_, err := q.Exec("SELECT 1")
		require.ErrorIs(t, err, ErrNoopQuerier)
	})

	t.Run("ExecContext returns ErrNoopQuerier", func(t *testing.T) {
		t.Parallel()
		q := newNoopQuerier()
		_, err := q.ExecContext(context.Background(), "SELECT 1")
		require.ErrorIs(t, err, ErrNoopQuerier)
	})

	t.Run("Query returns ErrNoopQuerier", func(t *testing.T) {
		t.Parallel()
		q := newNoopQuerier()
		_, err := q.Query("SELECT 1")
		require.ErrorIs(t, err, ErrNoopQuerier)
	})

	t.Run("QueryContext returns ErrNoopQuerier", func(t *testing.T) {
		t.Parallel()
		q := newNoopQuerier()
		_, err := q.QueryContext(context.Background(), "SELECT 1")
		require.ErrorIs(t, err, ErrNoopQuerier)
	})

	t.Run("QueryRow.Scan returns ErrNoopQuerier", func(t *testing.T) {
		t.Parallel()
		q := newNoopQuerier()
		var dest int
		err := q.QueryRow("SELECT 1").Scan(&dest)
		require.ErrorIs(t, err, ErrNoopQuerier)
	})

	t.Run("QueryRowContext.Err returns ErrNoopQuerier", func(t *testing.T) {
		t.Parallel()
		q := newNoopQuerier()
		err := q.QueryRowContext(context.Background(), "SELECT 1").Err()
		require.ErrorIs(t, err, ErrNoopQuerier)
	})
}
