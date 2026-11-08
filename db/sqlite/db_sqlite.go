package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/sonirico/vago/lol"
	_ "modernc.org/sqlite"
)

// Open opens the SQLite database at uri (a file path or a "file:" URI with
// modernc.org/sqlite query parameters) and bounds the pool to one connection:
// SQLite allows a single writer, and a pool of one turns "database is locked"
// errors into waits inside database/sql.
//
// Example URI:
//
//	file:app.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)
func Open(log lol.Logger, uri string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", uri)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Debug("sqlite connected")

	return db, nil
}
