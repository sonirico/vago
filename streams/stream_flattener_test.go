package streams

import (
	"fmt"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
)

func TestFlattenerStream(t *testing.T) {
	var (
		empty     []int
		errStream = errors.New("stream error")
	)

	tests := []struct {
		name        string
		inner       *MemoryStream[[]int]
		expected    []int
		expectedErr error
	}{
		{
			name: "basic flattening",
			inner: MemReader([][]int{
				{1, 2}, {3}, {4, 5},
			}, nil),
			expected:    []int{1, 2, 3, 4, 5},
			expectedErr: nil,
		},
		{
			name: "basic flattening with empty",
			inner: MemReader([][]int{
				{1, 2}, {3}, {}, {4, 5},
			}, nil),
			expected:    []int{1, 2, 3, 4, 5},
			expectedErr: nil,
		},
		{
			name:        "empty inner stream",
			inner:       MemReader([][]int{}, nil),
			expected:    empty,
			expectedErr: nil,
		},
		{
			name:        "single empty slice",
			inner:       MemReader([][]int{{}}, nil),
			expected:    empty,
			expectedErr: nil,
		},
		{
			name:        "multiple empty slices",
			inner:       MemReader([][]int{{}, {}}, nil),
			expected:    empty,
			expectedErr: nil,
		},
		{
			name: "error",
			inner: MemReader([][]int{
				{1, 2}, {3},
			}, errStream),
			expected:    empty,
			expectedErr: errStream,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := Flatten[int](tt.inner)
			result, err := Consume(stream)

			assert.ErrorIsf(
				t,
				err,
				tt.expectedErr,
				"expected error %v, got %v",
				tt.expectedErr,
				err,
			)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ExampleFlatten demonstrates flattening a stream of slices.
func ExampleFlatten() {
	// Create a stream from a slice of slices
	data := [][]int{{1, 2}, {3, 4, 5}, {6}, {7, 8, 9}}
	stream := MemReader(data, nil)

	// Flatten the slices
	flattened := Flatten(stream)

	// Collect the results
	result, _ := Consume(flattened)
	fmt.Println(result)
	// Output: [1 2 3 4 5 6 7 8 9]
}
