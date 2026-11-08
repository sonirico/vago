package sqlite

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/sonirico/vago/lol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sonirico/vago/db"
)

func TestOpenSqlite(t *testing.T) {
	t.Parallel()

	t.Run("creates the file and honours the DSN pragmas", func(t *testing.T) {
		t.Parallel()

		// Arrange
		uri := "file:" + filepath.Join(t.TempDir(), "app.db") + "?_pragma=foreign_keys(1)"

		// Act
		sdb, err := OpenSqlite(lol.ZeroTestLogger, uri)

		// Assert
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, sdb.Close()) })
		var foreignKeys int
		require.NoError(t, sdb.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys))
		assert.Equal(t, 1, foreignKeys)
	})

	t.Run("fails on an unwritable path", func(t *testing.T) {
		t.Parallel()

		// Arrange
		uri := filepath.Join(t.TempDir(), "missing", "app.db")

		// Act
		sdb, err := OpenSqlite(lol.ZeroTestLogger, uri)

		// Assert
		require.Error(t, err)
		assert.Nil(t, sdb)
	})

	t.Run("works through the database/sql executor in a transaction", func(t *testing.T) {
		t.Parallel()

		// Arrange
		sdb := newTestDB(t)
		executor := db.NewDatabaseSqlExecutor(lol.ZeroTestLogger, sdb)
		ctx := context.Background()
		_, err := sdb.Exec(`CREATE TABLE users (name TEXT NOT NULL)`)
		require.NoError(t, err)

		// Act
		err = executor.DoWithTx(ctx, func(ctx db.Context) error {
			_, err := ctx.Querier().ExecContext(ctx, `INSERT INTO users (name) VALUES (?)`, "alice")
			return err
		})

		// Assert
		require.NoError(t, err)
		var count int
		require.NoError(t, sdb.QueryRow(`SELECT count(*) FROM users`).Scan(&count))
		assert.Equal(t, 1, count)
	})
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	sdb, err := OpenSqlite(lol.ZeroTestLogger, filepath.Join(t.TempDir(), "app.db"))
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, sdb.Close()) })
	return sdb
}
