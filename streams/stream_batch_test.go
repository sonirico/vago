package streams

import (
	"fmt"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
)

func TestBatchStream(t *testing.T) {
	type testCase struct {
		name            string
		inputItems      []int
		inputErr        error
		batchSize       int
		expectedBatches [][]int
		expectedErr     error
	}

	var (
		errStream  = errors.New("stream error")
		errNoItems = errors.New("no items and error")
	)

	testCases := []testCase{
		{
			name:            "Empty input",
			inputItems:      []int{},
			batchSize:       3,
			expectedBatches: [][]int{},
			expectedErr:     nil,
		},
		{
			name:            "Fewer than batchSize",
			inputItems:      []int{1, 2},
			batchSize:       3,
			expectedBatches: [][]int{{1, 2}},
			expectedErr:     nil,
		},
		{
			name:            "Exactly one batch",
			inputItems:      []int{1, 2, 3},
			batchSize:       3,
			expectedBatches: [][]int{{1, 2, 3}},
			expectedErr:     nil,
		},
		{
			name:            "Multiple full batches",
			inputItems:      []int{1, 2, 3, 4, 5, 6},
			batchSize:       3,
			expectedBatches: [][]int{{1, 2, 3}, {4, 5, 6}},
			expectedErr:     nil,
		},
		{
			name:            "Multiple batches with a remainder",
			inputItems:      []int{1, 2, 3, 4, 5},
			batchSize:       3,
			expectedBatches: [][]int{{1, 2, 3}, {4, 5}},
			expectedErr:     nil,
		},
		{
			name:            "Error scenario at end",
			inputItems:      []int{1, 2, 3},
			batchSize:       2,
			expectedBatches: [][]int{{1, 2}, {3}},
			inputErr:        errStream,
			// Since the error occurs after consuming first 2 items, the second Next will return false, and Err should reflect the error.
			expectedErr: errStream,
		},
		{
			name:            "Error scenario with no items returned",
			inputItems:      []int{},
			batchSize:       2,
			inputErr:        errNoItems,
			expectedBatches: [][]int{},
			expectedErr:     errNoItems,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mem := MemReader(tc.inputItems, tc.inputErr)
			batch := Batch(mem, tc.batchSize)

			gotBatches := make([][]int, 0)

			for batch.Next() {
				gotBatch := batch.Data()
				// make a copy to avoid aliasing issues in subsequent loops
				gotBatches = append(gotBatches, gotBatch)
			}
			gotErr := batch.Err()

			assert.ErrorIs(t, gotErr, tc.expectedErr)
			assert.Equal(t, tc.expectedBatches, gotBatches)
		})
	}
}

// ExampleBatch demonstrates grouping stream elements into batches.
func ExampleBatch() {
	// Create a stream from a slice of integers
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stream := MemReader(data, nil)

	// Group into batches of 3
	batchStream := Batch(stream, 3)

	// Collect the results
	result, _ := Consume(batchStream)
	for i, batch := range result {
		fmt.Printf("Batch %d: %v\n", i+1, batch)
	}
	// Output:
	// Batch 1: [1 2 3]
	// Batch 2: [4 5 6]
	// Batch 3: [7 8 9]
	// Batch 4: [10]
}
