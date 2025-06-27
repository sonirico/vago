package streams

import (
	"testing"
)

func TestStreamChannel(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "empty channel",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{42},
			expected: []int{42},
		},
		{
			name:     "multiple elements",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare a channel and push input data
			ch := make(chan int, len(tc.input))
			for _, v := range tc.input {
				ch <- v
			}
			close(ch)

			// Create a new stream from the channel
			stream := Channel[int](ch)

			var actual []int
			for stream.Next() {
				actual = append(actual, stream.Data())
				if err := stream.Err(); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			// Verify that the actual read elements match the expected elements
			if len(actual) != len(tc.expected) {
				t.Fatalf("expected %d elements, got %d", len(tc.expected), len(actual))
			}
			for i, v := range tc.expected {
				if actual[i] != v {
					t.Errorf("expected element %d to be %d, got %d", i, v, actual[i])
				}
			}

			// Check Err is still nil at the end
			if stream.Err() != nil {
				t.Errorf("expected Err() to be nil at the end, got %v", stream.Err())
			}
		})
	}
}
