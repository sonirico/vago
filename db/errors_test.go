package db

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestErrIsNoRows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  error
		expect bool
	}{
		{"nil", nil, false},
		{"other", errors.New("foo"), false},
		{"sql.ErrNoRows", sql.ErrNoRows, true},
		{"pgx.ErrNoRows", pgx.ErrNoRows, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, ErrIsNoRows(tt.input))
		})
	}
}

func ExampleErrIsNoRows() {
	results := []bool{
		ErrIsNoRows(nil),
		ErrIsNoRows(errors.New("foo")),
		ErrIsNoRows(sql.ErrNoRows),
		ErrIsNoRows(pgx.ErrNoRows),
	}
	fmt.Println(results)
	// Output:
	// [false false true true]
}
