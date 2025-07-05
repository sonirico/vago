package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sonirico/vago/lol"

	"errors"

	"github.com/DATA-DOG/go-sqlmock"
)

type test struct {
	prepare     func() *sql.DB
	description string
	input       func(ctx Context) error
	want        error
}

var (
	ErrTest = errors.New("test error")
	log     = lol.ZeroTestLogger
)

func TestDatabaseServiceDo(t *testing.T) {

	tests := []test{
		{
			description: "Should work correctly",
			prepare: func() *sql.DB {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				if err != nil {
					t.Fatalf(
						"an error '%s' was not expected when opening a stub database connection",
						err,
					)
				}

				users := sqlmock.NewRows([]string{"id", "title"}).
					AddRow(1, "one")

				mock.ExpectQuery("SELECT * FROM users WHERE id = $1;").
					WithArgs(1).
					WillReturnRows(users)

				return db
			},
			input: func(ctx Context) error {
				type row struct {
					Id    string
					Title string
				}

				var tmp row

				var id = 1

				statement := `SELECT * FROM users WHERE id = $1;`

				if err := ctx.Querier().QueryRowContext(
					ctx,
					statement,
					id,
				).Scan(
					&tmp.Id,
					&tmp.Title,
				); err != nil {
					return err
				}

				return nil
			},
			want: nil,
		},
		{
			description: "Should fail",
			prepare: func() *sql.DB {
				db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				if err != nil {
					t.Fatalf(
						"an error '%s' was not expected when opening a stub database connection",
						err,
					)
				}

				return db
			},
			input: func(ctx Context) error {
				return ErrTest
			},
			want: ErrTest,
		},
	}

	for _, testcase := range tests {
		db := testcase.prepare()

		defer db.Close()

		dbSvc := newDatabaseSqlExecutor(log, db)

		err := dbSvc.Do(context.Background(), testcase.input)

		if !errors.Is(err, testcase.want) {
			t.Errorf("Expected: %v, got: %v", testcase.want, err)
		}

	}
}

func TestDatabaseServiceDoWithTx(t *testing.T) {

	tests := []test{
		{
			description: "Should work correctly",
			prepare: func() *sql.DB {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				if err != nil {
					t.Fatalf(
						"an error '%s' was not expected when opening a stub database connection",
						err,
					)
				}

				users := sqlmock.NewRows([]string{"id", "title"}).
					AddRow(1, "one")

				mock.ExpectBegin()

				mock.ExpectQuery("SELECT * FROM users WHERE id = $1;").
					WithArgs(1).
					WillReturnRows(users)

				mock.ExpectCommit()

				return db
			},
			input: func(ctx Context) error {
				type row struct {
					Id    string
					Title string
				}

				var tmp row

				var id = 1

				statement := `SELECT * FROM users WHERE id = $1;`

				if err := ctx.Querier().QueryRowContext(
					ctx,
					statement,
					id,
				).Scan(
					&tmp.Id,
					&tmp.Title,
				); err != nil {
					return err
				}

				return nil
			},
			want: nil,
		},
		{
			description: "Should fail",
			prepare: func() *sql.DB {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				if err != nil {
					t.Fatalf(
						"an error '%s' was not expected when opening a stub database connection",
						err,
					)
				}

				mock.ExpectBegin()
				mock.ExpectRollback()

				return db
			},
			input: func(ctx Context) error {
				return ErrTest
			},
			want: ErrTest,
		},
	}

	for _, testcase := range tests {
		db := testcase.prepare()

		defer db.Close()

		dbSvc := newDatabaseSqlExecutor(log, db)

		err := dbSvc.DoWithTx(context.Background(), testcase.input)

		if !errors.Is(err, testcase.want) {
			t.Errorf("Expected: %v, got: %v", testcase.want, err)
		}

	}
}
