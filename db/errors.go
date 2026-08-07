package db

import (
	"database/sql"

	"errors"
)

var (
	ErrDoesNotExist         = errors.New("does not exist")
	ErrConflict             = errors.New("conflict")
	NoUniqueConstraintError = errors.New(
		"It is mandatory including the PK ot having a unique key ina bulkable")
)

// ErrIsNoRows returns true if the error is a 'no rows' error from database/sql.
//
// Backend-specific "no rows" sentinels (e.g. pgx.ErrNoRows) are not checked
// here - core db has no driver dependency to check them against. A backend
// package that has its own sentinel provides its own ErrIsNoRows built on
// top of this one; see db/postgres.ErrIsNoRows.
func ErrIsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
