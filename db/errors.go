package db

import (
	"database/sql"

	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrDoesNotExist         = errors.New("does not exist")
	ErrConflict             = errors.New("conflict")
	NoUniqueConstraintError = errors.New(
		"It is mandatory including the PK ot having a unique key ina bulkable")
)

func ErrIsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows)
}
