package db

import (
	"context"
	"database/sql"
	"fmt"
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

func TestExecutorDo(t *testing.T) {

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

		ex := newDatabaseSqlExecutor(log, db)

		err := ex.Do(context.Background(), testcase.input)

		if !errors.Is(err, testcase.want) {
			t.Errorf("Expected: %v, got: %v", testcase.want, err)
		}

	}
}

func TestExecutorDoWithTx(t *testing.T) {

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

		ex := newDatabaseSqlExecutor(log, db)

		err := ex.DoWithTx(context.Background(), testcase.input)

		if !errors.Is(err, testcase.want) {
			t.Errorf("Expected: %v, got: %v", testcase.want, err)
		}

	}
}

// Example for Do and DoWithTx usage with a database service and context.
func ExampleExecutor() {
	db, mock, _ := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer db.Close()
	log := lol.ZeroTestLogger

	// Setup mock expectations
	mock.ExpectQuery("SELECT 1;").
		WillReturnRows(sqlmock.NewRows([]string{"n"}).AddRow(1))

	ex := newDatabaseSqlExecutor(log, db)

	err := ex.Do(context.Background(), func(ctx Context) error {
		var n int
		return ctx.Querier().
			QueryRowContext(ctx, "SELECT 1;").Scan(&n)
	})
	fmt.Println("Do error:", err)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT 2;").
		WillReturnRows(sqlmock.NewRows([]string{"n"}).AddRow(2))
	mock.ExpectCommit()

	err = ex.DoWithTx(context.Background(), func(ctx Context) error {
		var n int
		return ctx.Querier().QueryRowContext(ctx, "SELECT 2;").Scan(&n)
	})
	fmt.Println("DoWithTx error:", err)

	// Output:
	// Do error: <nil>
	// DoWithTx error: <nil>
}
