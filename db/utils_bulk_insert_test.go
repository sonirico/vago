package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBulkInsertSQL(t *testing.T) {
	tests := []struct {
		name         string
		rows         BulkableRanger
		tableName    string
		expectedStmt string
		expectedArgs []any
		expectedErr  error
	}{
		{
			name: "Bulk insert two rows, three columns",
			rows: BulkRanger[Bulkable]([]Bulkable{
				&mockBulkable{
					ColsVal: []string{"id", "name", "value"},
					RowVal:  []any{1, "foo", 100},
				},
				&mockBulkable{
					ColsVal: []string{"id", "name", "value"},
					RowVal:  []any{2, "bar", 200},
				},
			}),
			tableName: "my_table",
			expectedStmt: `
				INSERT INTO my_table
				(id,name,value)
				VALUES
				($1,$2,$3),($4,$5,$6)
			`,
			expectedArgs: []any{1, "foo", 100, 2, "bar", 200},
			expectedErr:  nil,
		},
		{
			name: "Bulk insert one row, two columns",
			rows: BulkRanger[Bulkable]([]Bulkable{
				&mockBulkable{
					ColsVal: []string{"id", "value"},
					RowVal:  []any{10, 999},
				},
			}),
			tableName: "test_table",
			expectedStmt: `
				INSERT INTO test_table
				(id,value)
				VALUES
				($1,$2)
			`,
			expectedArgs: []any{10, 999},
			expectedErr:  nil,
		},
		{
			name:         "Empty rows returns error",
			rows:         BulkRanger[Bulkable]([]Bulkable{}),
			tableName:    "table",
			expectedStmt: "",
			expectedArgs: nil,
			expectedErr:  assert.AnError, // will check error is not nil
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stmt, args, err := BulkInsertSQL(test.rows, test.tableName)
			if test.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
			assert.Equal(t, normalizeSQL(test.expectedStmt), normalizeSQL(stmt))
			assert.Equal(t, test.expectedArgs, args)
		})
	}
}

// ExampleBulkInsertSQL demonstrates how to use BulkInsertSQL to generate an SQL statement for bulk insertion.
func ExampleBulkInsertSQL() {
	rows := BulkRanger[Bulkable]([]Bulkable{
		&mockBulkable{
			ColsVal: []string{"id", "name", "value"},
			RowVal:  []any{1, "foo", 100},
		},
		&mockBulkable{
			ColsVal: []string{"id", "name", "value"},
			RowVal:  []any{2, "bar", 200},
		},
	})
	query, args, _ := BulkInsertSQL(rows, "my_table")
	fmt.Println("SQL:", normalizeSQL(query))
	fmt.Println("ARGS:", args)
	// Output:
	// SQL: INSERT INTO my_table (id,name,value) VALUES ($1,$2,$3),($4,$5,$6)
	// ARGS: [1 foo 100 2 bar 200]
}
