package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/sonirico/vago/lol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sonirico/vago/db"
)

func TestMigrate(t *testing.T) {
	t.Parallel()

	t.Run("unknown action is refused before opening the database", func(t *testing.T) {
		t.Parallel()

		// Arrange: a path inside a missing directory cannot be opened, so an
		// open attempt would surface as a different error.
		cfg := db.MigrationsConfig{Url: filepath.Join(t.TempDir(), "missing", "app.db")}

		// Act
		err := Migrate(cfg, "aplly", lol.ZeroTestLogger)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown migration action")
	})

	t.Run("up applies the migrations and down reverts the last one", func(t *testing.T) {
		t.Parallel()

		// Arrange
		cfg := newTestMigrations(t)

		// Act
		upErr := Migrate(cfg, string(db.ActionUp), lol.ZeroTestLogger)
		downErr := Migrate(cfg, string(db.ActionDown), lol.ZeroTestLogger)

		// Assert
		require.NoError(t, upErr)
		require.NoError(t, downErr)
		sdb, err := Open(lol.ZeroTestLogger, cfg.Url)
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, sdb.Close()) })
		assert.Equal(t, []string{"users"}, tableNames(t, sdb))
	})

	t.Run("up on an up-to-date database is not an error", func(t *testing.T) {
		t.Parallel()

		// Arrange
		cfg := newTestMigrations(t)
		require.NoError(t, Migrate(cfg, string(db.ActionUp), lol.ZeroTestLogger))

		// Act
		err := Migrate(cfg, string(db.ActionUp), lol.ZeroTestLogger)

		// Assert
		assert.NoError(t, err)
	})
}

// newTestMigrations writes two migrations (users, then posts) and points the
// config at them and at a fresh database file.
func newTestMigrations(t *testing.T) db.MigrationsConfig {
	t.Helper()
	dir := t.TempDir()
	migrations := map[string]string{
		"1_users.up.sql":   `CREATE TABLE users (id INTEGER PRIMARY KEY);`,
		"1_users.down.sql": `DROP TABLE users;`,
		"2_posts.up.sql":   `CREATE TABLE posts (id INTEGER PRIMARY KEY);`,
		"2_posts.down.sql": `DROP TABLE posts;`,
	}
	for name, body := range migrations {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644))
	}
	return db.MigrationsConfig{
		Url:            filepath.Join(dir, "app.db"),
		MigrationsPath: "file://" + dir,
	}
}

func tableNames(t *testing.T, sdb *sql.DB) []string {
	t.Helper()
	rows, err := sdb.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'schema_%' ORDER BY name`)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, rows.Close()) })
	var names []string
	for rows.Next() {
		var name string
		require.NoError(t, rows.Scan(&name))
		names = append(names, name)
	}
	require.NoError(t, rows.Err())
	return names
}
