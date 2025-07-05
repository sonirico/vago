package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullJSON_Scan_Value(t *testing.T) {
	t.Parallel()

	type M map[string]int

	m := M{"a": 1, "b": 2}
	b, _ := json.Marshal(m)

	n := NullJSON[string, int]{}
	assert.NoError(t, n.Scan(b))
	assert.True(t, n.Valid)
	assert.Equal(t, map[string]int(m), n.JSON)

	v, err := n.Value()
	assert.NoError(t, err)
	assert.Equal(t, driver.Value(b), v)
}

func TestNullJSON_Scan_nil(t *testing.T) {
	n := NullJSON[string, int]{}
	assert.NoError(t, n.Scan(nil))
	assert.False(t, n.Valid)
	assert.Nil(t, n.JSON)
}

func ExampleNullJSON_usage() {
	n := NullJSON[string, int]{}
	_ = n.Scan([]byte(`{"a":1}`))
	_, _ = n.Value()
	fmt.Println(n.JSON)
	fmt.Println(n.Valid)
	// Output:

	// map[a:1]
	// true
}
