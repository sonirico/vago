package streams

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamChannel(t *testing.T) {
	type testCase struct {
		name     string
		input    []int
		expected []int
	}

	tests := []testCase{
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
			t.Parallel()

			ch := make(chan int, len(tc.input))
			for _, v := range tc.input {
				ch <- v
			}
			close(ch)

			stream := Channel[int](ch)

			var actual []int
			for stream.Next(context.Background()) {
				actual = append(actual, stream.Data())
				require.NoError(t, stream.Err())
			}

			if len(tc.expected) == 0 {
				assert.Empty(t, actual)
			} else {
				assert.Equal(t, tc.expected, actual)
			}
			assert.NoError(t, stream.Err())
		})
	}
}

func TestStreamChannel_NextUnblocksOnContextCancellation(t *testing.T) {
	t.Parallel()

	ch := make(chan int)
	stream := Channel[int](ch)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan bool, 1)
	go func() {
		done <- stream.Next(ctx)
	}()

	cancel()

	select {
	case ok := <-done:
		assert.False(t, ok, "Next should return false once ctx is cancelled")
	case <-time.After(2 * time.Second):
		t.Fatal("Next did not unblock after context cancellation")
	}

	assert.ErrorIs(t, stream.Err(), context.Canceled)
}

func ExampleChannel() {
	// Create a channel and send data to it
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)

	// Create a stream from the channel
	stream := Channel[int](ch)

	// Process all values from the channel
	for stream.Next(context.Background()) {
		fmt.Println(stream.Data())
	}
	// Output:
	// 10
	// 20
	// 30
}
