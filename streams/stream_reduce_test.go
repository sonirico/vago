package streams

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	tests := []struct {
		name     string
		data     []int
		fn       func(int, int) int
		initial  int
		expected int
		wantErr  bool
	}{
		{
			name:     "sum numbers",
			data:     []int{1, 2, 3, 4, 5},
			fn:       func(acc, n int) int { return acc + n },
			initial:  0,
			expected: 15,
			wantErr:  false,
		},
		{
			name:     "multiply numbers",
			data:     []int{2, 3, 4},
			fn:       func(acc, n int) int { return acc * n },
			initial:  1,
			expected: 24,
			wantErr:  false,
		},
		{
			name: "find maximum",
			data: []int{5, 2, 8, 1, 9, 3},
			fn: func(acc, n int) int {
				if n > acc {
					return n
				}
				return acc
			},
			initial:  0,
			expected: 9,
			wantErr:  false,
		},
		{
			name:     "empty stream",
			data:     []int{},
			fn:       func(acc, n int) int { return acc + n },
			initial:  10,
			expected: 10,
			wantErr:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stream := MemReader(tc.data, nil)
			result, err := Reduce(stream, tc.fn, tc.initial)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestReduceWithError(t *testing.T) {
	stream := MemReader([]int{1, 2, 3}, errors.New("stream error"))
	_, err := Reduce(stream, func(acc, n int) int { return acc + n }, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "stream error")
}

func TestReduceSlice(t *testing.T) {
	tests := []struct {
		name     string
		data     []string
		fn       func([]string, string) []string
		expected []string
		wantErr  bool
	}{
		{
			name:     "collect all items",
			data:     []string{"a", "b", "c"},
			fn:       func(acc []string, item string) []string { return append(acc, item) },
			expected: []string{"a", "b", "c"},
			wantErr:  false,
		},
		{
			name: "filter and collect",
			data: []string{"apple", "a", "application", "ab"},
			fn: func(acc []string, item string) []string {
				if len(item) > 2 {
					return append(acc, item)
				}
				return acc
			},
			expected: []string{"apple", "application"},
			wantErr:  false,
		},
		{
			name: "reverse collect",
			data: []string{"first", "second", "third"},
			fn: func(acc []string, item string) []string {
				return append([]string{item}, acc...)
			},
			expected: []string{"third", "second", "first"},
			wantErr:  false,
		},
		{
			name:     "empty stream",
			data:     []string{},
			fn:       func(acc []string, item string) []string { return append(acc, item) },
			expected: []string{},
			wantErr:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stream := MemReader(tc.data, nil)
			result, err := ReduceSlice(stream, tc.fn)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestReduceSliceWithError(t *testing.T) {
	stream := MemReader([]string{"a", "b"}, errors.New("slice error"))
	_, err := ReduceSlice(stream, func(acc []string, item string) []string {
		return append(acc, item)
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "slice error")
}

func TestReduceMap(t *testing.T) {
	tests := []struct {
		name     string
		data     []string
		fn       func(map[string]int, string) map[string]int
		expected map[string]int
		wantErr  bool
	}{
		{
			name: "count occurrences",
			data: []string{"apple", "banana", "apple", "cherry", "banana", "apple"},
			fn: func(acc map[string]int, word string) map[string]int {
				acc[word]++
				return acc
			},
			expected: map[string]int{"apple": 3, "banana": 2, "cherry": 1},
			wantErr:  false,
		},
		{
			name: "track lengths",
			data: []string{"cat", "dog", "elephant", "cat"},
			fn: func(acc map[string]int, word string) map[string]int {
				acc[word] = len(word)
				return acc
			},
			expected: map[string]int{"cat": 3, "dog": 3, "elephant": 8},
			wantErr:  false,
		},
		{
			name: "conditional mapping",
			data: []string{"a", "bb", "ccc", "d", "ee"},
			fn: func(acc map[string]int, word string) map[string]int {
				if len(word) > 1 {
					acc[word] = len(word)
				}
				return acc
			},
			expected: map[string]int{"bb": 2, "ccc": 3, "ee": 2},
			wantErr:  false,
		},
		{
			name: "empty stream",
			data: []string{},
			fn: func(acc map[string]int, word string) map[string]int {
				acc[word]++
				return acc
			},
			expected: map[string]int{},
			wantErr:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stream := MemReader(tc.data, nil)
			result, err := ReduceMap(stream, tc.fn)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestReduceMapWithError(t *testing.T) {
	stream := MemReader([]string{"a", "b"}, errors.New("map error"))
	_, err := ReduceMap(stream, func(acc map[string]int, word string) map[string]int {
		acc[word]++
		return acc
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "map error")
}

// ExampleReduce demonstrates using Reduce to sum numbers from a stream.
func ExampleReduce() {
	// Create a stream of numbers from strings
	reader := strings.NewReader("10\n20\n30\n40\n50")
	lines := Lines(reader)

	// Convert strings to numbers and sum them
	sum, _ := Reduce(Map(lines, func(s string) int {
		n, _ := strconv.Atoi(s)
		return n
	}), func(acc, n int) int {
		return acc + n
	}, 0)

	fmt.Printf("Sum: %d\n", sum)

	// Output:
	// Sum: 150
}

// ExampleReduceSlice demonstrates collecting filtered items from a stream.
func ExampleReduceSlice() {
	// Create a stream of words
	reader := strings.NewReader("cat\ndog\nelephant\nant\nbutterfly\nbird")
	stream := Lines(reader)

	// Collect only words longer than 3 characters
	longWords, _ := ReduceSlice(stream, func(acc []string, word string) []string {
		if len(word) > 3 {
			return append(acc, word)
		}
		return acc
	})

	fmt.Printf("Long words: %v\n", longWords)

	// Output:
	// Long words: [elephant butterfly bird]
}

// ExampleReduceMap demonstrates reducing a stream to a map with aggregated values.
func ExampleReduceMap() {
	// Create a stream of words
	reader := strings.NewReader("apple\nbanana\napple\ncherry\nbanana\napple")
	stream := Lines(reader)

	// Count occurrences of each word
	counts, _ := ReduceMap(stream, func(acc map[string]int, word string) map[string]int {
		acc[word]++
		return acc
	})

	fmt.Printf("apple: %d\n", counts["apple"])
	fmt.Printf("banana: %d\n", counts["banana"])
	fmt.Printf("cherry: %d\n", counts["cherry"])

	// Output:
	// apple: 3
	// banana: 2
	// cherry: 1
}
