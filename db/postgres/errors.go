package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/sonirico/vago/db"
)

// ErrIsNoRows returns true if the error is a 'no rows' error from either
// database/sql or pgx. It extends db.ErrIsNoRows (which only checks
// database/sql's sentinel) with pgx's own, since core db has no pgx
// dependency to check that against.
func ErrIsNoRows(err error) bool {
	return db.ErrIsNoRows(err) || errors.Is(err, pgx.ErrNoRows)
}
