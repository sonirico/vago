package streams

import (
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
			mem := NewMemory(tc.inputItems, tc.inputErr)
			batch := NewBatchStream(mem, tc.batchSize)

			gotBatches := make([][]int, 0)

			for batch.Next() {
				gotBatch := batch.Data()
				// make a copy to avoid aliasing issues in subsequent loops
				copied := make([]int, len(gotBatch))
				copy(copied, gotBatch)
				gotBatches = append(gotBatches, copied)
			}
			gotErr := batch.Err()

			assert.ErrorIs(t, gotErr, tc.expectedErr)
			assert.Equal(t, tc.expectedBatches, gotBatches)
		})
	}
}
