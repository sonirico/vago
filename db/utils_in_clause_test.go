package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []any
		inArgs  []int
		expectS string
		expectA []any
	}{
		{"empty", nil, nil, "()", nil},
		{"one", nil, []int{1}, "($1)", []any{1}},
		{"three", []any{"foo"}, []int{1, 2, 3}, "($2,$3,$4)", []any{"foo", 1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, a := In(tt.args, tt.inArgs)
			assert.Equal(t, tt.expectS, s)
			assert.Equal(t, tt.expectA, a)
		})
	}
}

// ExampleIn demonstrates how to use the In function to generate an SQL IN clause and its arguments.
func ExampleIn() {
	args := []any{"foo"}
	inArgs := []int{1, 2, 3}
	s, a := In(args, inArgs)
	fmt.Println(s)
	fmt.Println(a)
	// Output:
	// ($2,$3,$4)
	// [foo 1 2 3]
}
