package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockBulkable struct {
	PKVal         []string
	UniqueKeysVal []string
	IncludePKVal  bool
	ColsVal       []string
	RowVal        []any
}

func (m *MockBulkable) PK() []string {
	return m.PKVal
}

func (m *MockBulkable) UniqueKeys() []string {
	return m.UniqueKeysVal
}

func (m *MockBulkable) IncludePKOnUpsert() bool {
	return m.IncludePKVal
}

func (m *MockBulkable) Cols() []string {
	return m.ColsVal
}

func (m *MockBulkable) Row() []any {
	return m.RowVal
}

func TestBuildSqlStmt(t *testing.T) {
	type args struct {
		rows             BulkableRanger
		tableName        string
		updateOnConflict bool
	}

	tests := []struct {
		name         string
		args         args
		expectedStmt string
		expectedArgs []any
		expectedErr  error
	}{
		{
			name: "Upsert (update on conflict) unique key and include pk",
			args: args{
				rows: BulkRanger[Bulkable](
					[]Bulkable{
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{"unique_key"},
							IncludePKVal:  true,
							ColsVal:       []string{"id", "unique_key", "value", "extra_value"},
							RowVal:        []any{1, "abc", 100, "extra1"},
						},
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{"unique_key"},
							IncludePKVal:  true,
							ColsVal:       []string{"id", "unique_key", "value", "extra_value"},
							RowVal:        []any{2, "def", 200, "extra2"},
						},
					},
				),
				tableName:        "my_table",
				updateOnConflict: true,
			},
			expectedErr: nil,
			expectedStmt: `
		INSERT INTO my_table
		(id,unique_key,value,extra_value)
		VALUES
		($1,$2,$3,$4),($5,$6,$7,$8)
		ON CONFLICT (unique_key) 
			DO UPDATE SET
			id = EXCLUDED.id,unique_key = EXCLUDED.unique_key,value = EXCLUDED.value,extra_value = EXCLUDED.extra_value
		RETURNING id,unique_key,value,extra_value
		`,
			expectedArgs: []any{1, "abc", 100, "extra1", 2, "def", 200, "extra2"},
		},
		{
			name: "Upsert (update on conflict) unique key and exclude pk",
			args: args{
				rows: BulkRanger[Bulkable](
					[]Bulkable{
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{"unique_key"},
							IncludePKVal:  false,
							ColsVal:       []string{"unique_key", "value", "extra_value"},
							RowVal:        []any{"abc", 100, "extra1"},
						},
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{"unique_key"},
							IncludePKVal:  false,
							ColsVal:       []string{"unique_key", "value", "extra_value"},
							RowVal:        []any{"def", 200, "extra2"},
						},
					},
				),
				tableName:        "my_table",
				updateOnConflict: true,
			},
			expectedErr: nil,
			expectedStmt: `
		INSERT INTO my_table
		(unique_key,value,extra_value)
		VALUES
		($1,$2),($3,$4)
		ON CONFLICT (unique_key) 
			DO UPDATE SET
			unique_key = EXCLUDED.unique_key,value = EXCLUDED.value,extra_value = EXCLUDED.extra_value
		RETURNING unique_key,value,extra_value
		`,
			expectedArgs: []any{"abc", 100, "extra1", "def", 200, "extra2"},
		},
		{
			name: "Upsert (update on conflict)  no unique key and include pk",
			args: args{
				rows: BulkRanger[Bulkable](
					[]Bulkable{
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{},
							IncludePKVal:  true,
							ColsVal:       []string{"id", "value", "extra_value"},
							RowVal:        []any{1, 100, "extra1"},
						},
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{},
							IncludePKVal:  true,
							ColsVal:       []string{"id", "value", "extra_value"},
							RowVal:        []any{2, 200, "extra2"},
						},
					},
				),
				tableName:        "my_table",
				updateOnConflict: true,
			},
			expectedErr: nil,
			expectedStmt: `
		INSERT INTO my_table
		(id,value,extra_value)
		VALUES
		($1,$2,$3),($4,$5,$6)
		ON CONFLICT (id) 
			DO UPDATE SET
			id = EXCLUDED.id,value = EXCLUDED.value,extra_value = EXCLUDED.extra_value
		RETURNING id,value,extra_value
		`,
			expectedArgs: []any{1, 100, "extra1", 2, 200, "extra2"},
		},
		{
			name: "Upsert (update on conflict)  no unique key and exclude pk",
			args: args{
				rows: BulkRanger[Bulkable](
					[]Bulkable{
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{},
							IncludePKVal:  false,
							ColsVal:       []string{"value", "extra_value"},
							RowVal:        []any{100, "extra1"},
						},
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{},
							IncludePKVal:  false,
							ColsVal:       []string{"value", "extra_value"},
							RowVal:        []any{200, "extra2"},
						},
					},
				),
				tableName:        "my_table",
				updateOnConflict: true,
			},
			expectedErr:  NoUniqueConstraintError,
			expectedStmt: "",
			expectedArgs: nil,
		},
		{
			name: "Upsert (ignore on conflict) unique key and include pk", args: args{
				rows: BulkRanger[Bulkable](
					[]Bulkable{
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{"unique_key"},
							IncludePKVal:  true,
							ColsVal:       []string{"id", "unique_key", "value", "extra_value"},
							RowVal:        []any{3, "xyz", 300, "extra3"},
						},
					},
				),
				tableName:        "my_table",
				updateOnConflict: false,
			},
			expectedErr: nil,
			expectedStmt: `
		INSERT INTO my_table
		(id,unique_key,value,extra_value)
		VALUES
		($1,$2,$3,$4)
		ON CONFLICT (unique_key) DO NOTHING
		RETURNING id,unique_key,value,extra_value
		`,
			expectedArgs: []any{3, "xyz", 300, "extra3"},
		},
		{
			name: "Upsert (ignore on conflict) unique key and exclude pk",
			args: args{
				rows: BulkRanger[Bulkable](
					[]Bulkable{
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{"unique_key"},
							IncludePKVal:  false,
							ColsVal:       []string{"unique_key", "value", "extra_value"},
							RowVal:        []any{"abc", 100, "extra1"},
						},
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{"unique_key"},
							IncludePKVal:  false,
							ColsVal:       []string{"unique_key", "value", "extra_value"},
							RowVal:        []any{"def", 200, "extra2"},
						},
					},
				),
				tableName:        "my_table",
				updateOnConflict: false,
			},
			expectedErr: nil,
			expectedStmt: `
		INSERT INTO my_table
		(unique_key,value,extra_value)
		VALUES
		($1,$2),($3,$4)
		ON CONFLICT (unique_key) DO NOTHING
		RETURNING unique_key,value,extra_value
		`,
			expectedArgs: []any{"abc", 100, "extra1", "def", 200, "extra2"},
		},
		{
			name: "Upsert (ignore on conflict)  no unique key and include pk",
			args: args{
				rows: BulkRanger[Bulkable](
					[]Bulkable{
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{},
							IncludePKVal:  true,
							ColsVal:       []string{"id", "value", "extra_value"},
							RowVal:        []any{1, 100, "extra1"},
						},
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{},
							IncludePKVal:  true,
							ColsVal:       []string{"id", "value", "extra_value"},
							RowVal:        []any{2, 200, "extra2"},
						},
					},
				),
				tableName:        "my_table",
				updateOnConflict: false,
			},
			expectedErr: nil,
			expectedStmt: `
		INSERT INTO my_table
		(id,value,extra_value)
		VALUES
		($1,$2,$3),($4,$5,$6)
		ON CONFLICT (id) DO NOTHING
		RETURNING id,value,extra_value
		`,
			expectedArgs: []any{1, 100, "extra1", 2, 200, "extra2"},
		},
		{
			name: "Upsert (ignore on conflict)  no unique key and exclude pk",
			args: args{
				rows: BulkRanger[Bulkable](
					[]Bulkable{
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{},
							IncludePKVal:  false,
							ColsVal:       []string{"value", "extra_value"},
							RowVal:        []any{100, "extra1"},
						},
						&MockBulkable{
							PKVal:         []string{"id"},
							UniqueKeysVal: []string{},
							IncludePKVal:  false,
							ColsVal:       []string{"value", "extra_value"},
							RowVal:        []any{200, "extra2"},
						},
					},
				),
				tableName:        "my_table",
				updateOnConflict: false,
			},
			expectedErr:  NoUniqueConstraintError,
			expectedStmt: "",
			expectedArgs: nil,
		},
	}

	for _, test := range tests {
		t.Run("TestBuildSqlStmt: "+test.name, func(t *testing.T) {
			stmt, args, err := buildSqlStmt(
				test.args.rows,
				test.args.tableName,
				test.args.updateOnConflict,
			)

			if test.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.ErrorIs(t, err, test.expectedErr)
			}

			assert.Equal(t, test.expectedStmt, stmt)
			assert.Equal(t, test.expectedArgs, args)
		})
	}
}
