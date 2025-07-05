package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testBulkable struct {
	pk    []string
	uk    []string
	cols  []string
	row   []any
	pkInc bool
}

func (t *testBulkable) PK() []string            { return t.pk }
func (t *testBulkable) UniqueKeys() []string    { return t.uk }
func (t *testBulkable) IncludePKOnUpsert() bool { return t.pkInc }
func (t *testBulkable) Cols() []string          { return t.cols }
func (t *testBulkable) Row() []any              { return t.row }

func TestBulkRanger(t *testing.T) {
	items := []*testBulkable{
		{pk: []string{"id"}, cols: []string{"id", "v"}, row: []any{1, "a"}},
		{pk: []string{"id"}, cols: []string{"id", "v"}, row: []any{2, "b"}},
	}
	br := BulkRanger[*testBulkable](items)
	assert.Equal(t, 2, br.Len())
	assert.Equal(t, items[1], br.Get(1))
	var seen []any
	br.Range(func(b Bulkable) { seen = append(seen, b.Row()...) })
	assert.ElementsMatch(t, []any{1, "a", 2, "b"}, seen)
}

func ExampleBulkRanger_usage() {
	items := []*testBulkable{
		{pk: []string{"id"}, cols: []string{"id", "v"}, row: []any{1, "a"}},
		{pk: []string{"id"}, cols: []string{"id", "v"}, row: []any{2, "b"}},
	}
	br := BulkRanger[*testBulkable](items)
	_ = br.Len()
	_ = br.Get(0)
	br.Range(func(b Bulkable) {})
	// Output:
}
