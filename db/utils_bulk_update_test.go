package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockBulkUpdatable struct {
	pk   [2]string
	cols [][2]string
	vals []any
}

func (m *mockBulkUpdatable) PK() [2]string               { return m.pk }
func (m *mockBulkUpdatable) BulkUpdateCols() [][2]string { return m.cols }
func (m *mockBulkUpdatable) BulkUpdateValues() []any     { return m.vals }

type mockBulkUpdate []BulkUpdatable

func (m mockBulkUpdate) Len() int                { return len(m) }
func (m mockBulkUpdate) Get(i int) BulkUpdatable { return m[i] }
func (m mockBulkUpdate) Range(f func(x BulkUpdatable)) {
	for _, x := range m {
		f(x)
	}
}

func TestBulkUpdateSQL(t *testing.T) {
	tests := []struct {
		name         string
		rows         bulkUpdate
		tableName    string
		expectedStmt string
		expectedArgs []any
		expectedErr  error
	}{
		{
			name: "Bulk update two rows, two columns",
			rows: mockBulkUpdate{
				&mockBulkUpdatable{
					pk:   [2]string{"int", "id"},
					cols: [][2]string{{"int", "value"}, {"text", "name"}},
					vals: []any{1, "foo", 100, "bar"},
				},
				&mockBulkUpdatable{
					pk:   [2]string{"int", "id"},
					cols: [][2]string{{"int", "value"}, {"text", "name"}},
					vals: []any{2, "baz", 200, "qux"},
				},
			},
			tableName: "my_table",
			expectedStmt: `
				UPDATE my_table
				SET value = bulk_update_tmp.value::int,name = bulk_update_tmp.name::text
				FROM (VALUES ($1::int, $2::int, $3::text), ($4::int, $5::int, $6::text)) as bulk_update_tmp(id, value, name)
				WHERE my_table.id::int = bulk_update_tmp.id::int
			`,
			expectedArgs: []any{1, "foo", 100, "bar", 2, "baz", 200, "qux"},
			expectedErr:  nil,
		},
		{
			name:         "Empty rows returns error",
			rows:         mockBulkUpdate{},
			tableName:    "table",
			expectedStmt: "",
			expectedArgs: nil,
			expectedErr:  assert.AnError, // will check error is not nil
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stmt, args, err := BulkUpdateSQL(test.rows, test.tableName)
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

// ExampleBulkUpdateSQL demonstrates how to use BulkUpdateSQL to generate an SQL statement for bulk updates.
func ExampleBulkUpdateSQL() {
	rows := mockBulkUpdate{
		&mockBulkUpdatable{
			pk:   [2]string{"int", "id"},
			cols: [][2]string{{"int", "value"}, {"text", "name"}},
			vals: []any{1, "foo", 100, "bar"},
		},
		&mockBulkUpdatable{
			pk:   [2]string{"int", "id"},
			cols: [][2]string{{"int", "value"}, {"text", "name"}},
			vals: []any{2, "baz", 200, "qux"},
		},
	}
	query, args, _ := BulkUpdateSQL(rows, "my_table")
	fmt.Println("SQL:", normalizeSQL(query))
	fmt.Println("ARGS:", args)
	// Output:
	// SQL: UPDATE my_table SET value = bulk_update_tmp.value::int,name = bulk_update_tmp.name::text FROM (VALUES ($1::int, $2::int, $3::text), ($4::int, $5::int, $6::text)) as bulk_update_tmp(id, value, name) WHERE my_table.id::int = bulk_update_tmp.id::int
	// ARGS: [1 foo 100 bar 2 baz 200 qux]
}
