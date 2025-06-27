package streams

import (
	"errors"
	"io"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsume(t *testing.T) {
	tests := []struct {
		name     string
		stream   *MemoryStream[int]
		expected []int
		wantErr  error
	}{
		{
			name:     "Stream with data",
			stream:   MemReader([]int{1, 2, 3, 4, 5}, nil),
			expected: []int{1, 2, 3, 4, 5},
			wantErr:  nil,
		},
		{
			name:     "Empty stream",
			stream:   MemReader([]int{}, nil),
			expected: []int{},
			wantErr:  nil,
		},
		{
			name:     "Stream with EOF error (should return data)",
			stream:   MemReader([]int{10, 20, 30}, io.EOF),
			expected: []int{10, 20, 30},
			wantErr:  nil,
		},
		{
			name:     "Stream with generic error",
			stream:   MemReader([]int{100, 200}, errors.New("stream error")),
			expected: nil,
			wantErr:  errors.New("stream error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Consume(tc.stream)

			assert.Len(t, result, len(tc.expected), "Result does not match expected output")
			if len(tc.expected) > 0 {
				assert.Equal(t, tc.expected, result, "Result does not match expected output")
			}

			if tc.wantErr != nil {
				assert.Error(t, err, "Expected an error but got nil")
				assert.EqualError(t, err, tc.wantErr.Error(), "Unexpected error message")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
			}
		})
	}
}

// TestStreamPatterns demonstrates the three main stream patterns:
// 1. read -> write (simple pipe)
// 2. read -> filter -> map -> write (chained transformations)
// 3. read -> multiple writes (Multicast)
func TestStreamPatterns(t *testing.T) {
	t.Run("Simple Pipe (read -> write)", func(t *testing.T) {
		// Create source data
		sourceData := []string{"hello", "world", "golang", "streams"}
		src := MemReader(sourceData, nil)

		// Create destination
		dst := MemWriter[string]()

		// Simple pipe: read -> write
		bytesWritten, err := Pipe(src, dst)
		if err != nil {
			t.Fatalf("Pipe failed: %v", err)
		}

		t.Logf("Piped %d items (%d bytes)", len(dst.Items()), bytesWritten)

		// Verify all items were copied
		result := dst.Items()
		if len(result) != len(sourceData) {
			t.Errorf("Expected %d items, got %d", len(sourceData), len(result))
		}
	})

	t.Run("Chained Transformations (read -> filter -> map -> write)", func(t *testing.T) {
		// Create source data
		sourceData := []string{"hello", "world", "go", "programming", "test", "filter"}
		src := MemReader(sourceData, nil)

		// Chain transformations: filter strings with length > 3
		filtered := Filter(src, func(s string) bool { return len(s) > 3 })

		// Then map to uppercase
		mapped := Map(filtered, func(s string) string { return strings.ToUpper(s) })

		// Write to destination
		dst := MemWriter[string]()
		bytesWritten, err := Pipe(mapped, dst)
		if err != nil {
			t.Fatalf("Chained pipe failed: %v", err)
		}

		t.Logf(
			"Processed %d items (%d bytes) through filter+map chain",
			len(dst.Items()),
			bytesWritten,
		)

		// Verify transformation
		result := dst.Items()
		expected := []string{"HELLO", "WORLD", "PROGRAMMING", "TEST", "FILTER"}

		if len(result) != len(expected) {
			t.Errorf("Expected %d items, got %d", len(expected), len(result))
		}

		for i, expectedVal := range expected {
			if result[i] != expectedVal {
				t.Errorf("Expected item %d to be %s, got %s", i, expectedVal, result[i])
			}
		}
	})

	t.Run("Multicast (read -> multiple writes)", func(t *testing.T) {
		// Create source data
		sourceData := []int{1, 2, 3, 4, 5}
		src := MemReader(sourceData, nil)

		// Create multiple destinations
		dst1 := MemWriter[int]()
		dst2 := MemWriter[int]()
		dst3 := MemWriter[int]()

		// Multicast: read -> multiple writes
		bytesWritten, err := Multicast(src, dst1, dst2, dst3)
		if err != nil {
			t.Fatalf("Multicast failed: %v", err)
		}

		t.Logf("Multicasted to %d destinations: %v bytes", len(bytesWritten), bytesWritten)

		// Verify all destinations have the same data
		destinations := []*MemoryWriteStream[int]{dst1, dst2, dst3}
		for i, dst := range destinations {
			result := dst.Items()
			if len(result) != len(sourceData) {
				t.Errorf(
					"Destination %d: expected %d items, got %d",
					i,
					len(sourceData),
					len(result),
				)
			}

			for j, expectedVal := range sourceData {
				if result[j] != expectedVal {
					t.Errorf(
						"Destination %d, item %d: expected %d, got %d",
						i,
						j,
						expectedVal,
						result[j],
					)
				}
			}
		}
	})
}

// TestRealWorldExample shows a practical example combining all patterns
func TestRealWorldExample(t *testing.T) {
	// Simulate processing a log of user actions
	type UserAction struct {
		User   string
		Action string
		Score  int
	}

	// Source data
	actions := []UserAction{
		{"alice", "login", 1},
		{"bob", "purchase", 10},
		{"charlie", "view", 1},
		{"alice", "purchase", 15},
		{"bob", "logout", 1},
		{"david", "purchase", 25},
		{"alice", "view", 1},
	}

	src := MemReader(actions, nil)

	// Filter only high-value actions (score > 5)
	highValueFilter := func(action UserAction) bool { return action.Score > 5 }
	filtered := Filter(src, highValueFilter)

	// Map to user names only
	userMapper := func(action UserAction) string { return action.User }
	mapped := Map(filtered, userMapper)

	// Write to multiple destinations:
	// 1. All high-value users
	allUsers := MemWriter[string]()
	// 2. Alert system (same data)
	alertSystem := MemWriter[string]()

	// Process: read -> filter -> map -> Multicast to multiple writes
	bytesWritten, err := Multicast(mapped, allUsers, alertSystem)
	if err != nil {
		t.Fatalf("Real world example failed: %v", err)
	}

	t.Logf("Processed high-value users: %v bytes written", bytesWritten)

	// Verify results
	expectedUsers := []string{"bob", "alice", "david"}

	for i, dst := range []*MemoryWriteStream[string]{allUsers, alertSystem} {
		result := dst.Items()
		if len(result) != len(expectedUsers) {
			t.Errorf(
				"Destination %d: expected %d users, got %d",
				i,
				len(expectedUsers),
				len(result),
			)
		}

		for j, expectedUser := range expectedUsers {
			if result[j] != expectedUser {
				t.Errorf(
					"Destination %d, user %d: expected %s, got %s",
					i,
					j,
					expectedUser,
					result[j],
				)
			}
		}
	}
}

func TestReadAllBytes(t *testing.T) {
	t.Run("JSONTransform", func(t *testing.T) {
		data := []utilsTestMarshaler{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
		}
		stream := MemReader(data, nil)
		transform := JSONTransform(stream)

		result, err := ReadAllBytes[utilsTestMarshaler](transform)
		assert.NoError(t, err, "Should read all bytes without error")
		assert.NotEmpty(t, result, "Should return non-empty bytes")

		// Verify it's valid JSON
		expected := `[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]`
		assert.Equal(t, expected, string(result), "Should produce expected JSON")
	})

	t.Run("Transform with error", func(t *testing.T) {
		stream := MemReader([]utilsTestMarshaler{}, errors.New("stream error"))
		transform := JSONTransform(stream)

		result, err := ReadAllBytes[utilsTestMarshaler](transform)
		assert.Error(t, err, "Should return error from transform")
		assert.Nil(t, result, "Should return nil bytes on error")
	})
}

func TestConsumeErrSkip(t *testing.T) {
	t.Run("Stream with mixed data and errors", func(t *testing.T) {
		// Create a stream that will have some errors during iteration
		stream := MemReader([]int{1, 2, 3, 4, 5}, nil)

		// Consume all data (no errors in this case)
		result := ConsumeErrSkip(stream)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result, "Should consume all data")
	})

	t.Run("Empty stream", func(t *testing.T) {
		stream := MemReader([]int{}, nil)
		result := ConsumeErrSkip(stream)
		assert.Empty(t, result, "Should return empty slice for empty stream")
	})

	t.Run("Stream with error", func(t *testing.T) {
		stream := MemReader([]int{1, 2, 3}, errors.New("stream error"))
		result := ConsumeErrSkip(stream)
		// Should return empty slice when stream has error
		assert.Empty(t, result, "Should return empty result when stream has error")
	})
}

func TestWriteSeq(t *testing.T) {
	t.Run("Write sequence from slice", func(t *testing.T) {
		stream := MemWriter[int]()
		items := []int{1, 2, 3, 4, 5}

		bytesWritten, err := WriteSeq(stream, slices.Values(items))
		assert.NoError(t, err, "Should write sequence without error")
		assert.Positive(t, bytesWritten, "Should write positive bytes")

		result := stream.Items()
		assert.Equal(t, items, result, "Should write all items in sequence")
	})

	t.Run("Write empty sequence", func(t *testing.T) {
		stream := MemWriter[string]()
		items := []string{}

		bytesWritten, err := WriteSeq(stream, slices.Values(items))
		assert.NoError(t, err, "Should handle empty sequence")
		assert.Zero(t, bytesWritten, "Should write zero bytes for empty sequence")

		result := stream.Items()
		assert.Empty(t, result, "Should have no items")
	})

	t.Run("Stream with write error", func(t *testing.T) {
		stream := MemWriter[int]()
		stream.SetError(errors.New("write error"))
		items := []int{1, 2, 3}

		bytesWritten, err := WriteSeq(stream, slices.Values(items))
		assert.Error(t, err, "Should return error on write failure")
		assert.Zero(t, bytesWritten, "Should return zero bytes on error")
		assert.Contains(t, err.Error(), "write error", "Should contain write error message")
	})
}

func TestWriteSeqKeys(t *testing.T) {
	t.Run("Write keys from map", func(t *testing.T) {
		stream := MemWriter[string]()
		data := map[string]int{
			"apple":  1,
			"banana": 2,
			"cherry": 3,
		}

		bytesWritten, err := WriteSeqKeys(stream, maps.All(data))
		assert.NoError(t, err, "Should write keys without error")
		assert.Positive(t, bytesWritten, "Should write positive bytes")

		result := stream.Items()
		assert.Len(t, result, 3, "Should write all keys")

		// Check that all keys are present (order may vary)
		for key := range data {
			assert.Contains(t, result, key, "Should contain key %s", key)
		}
	})
}

func TestWriteSeqValues(t *testing.T) {
	t.Run("Write values from map", func(t *testing.T) {
		stream := MemWriter[int]()
		data := map[string]int{
			"apple":  1,
			"banana": 2,
			"cherry": 3,
		}

		bytesWritten, err := WriteSeqValues(stream, maps.All(data))
		assert.NoError(t, err, "Should write values without error")
		assert.Positive(t, bytesWritten, "Should write positive bytes")

		result := stream.Items()
		assert.Len(t, result, 3, "Should write all values")

		// Check that all values are present (order may vary)
		for _, value := range data {
			assert.Contains(t, result, value, "Should contain value %d", value)
		}
	})
}

func TestPipeErrorHandling(t *testing.T) {
	t.Run("Source stream error", func(t *testing.T) {
		src := MemReader([]int{1, 2}, errors.New("source error"))
		dst := MemWriter[int]()

		bytesWritten, err := Pipe(src, dst)
		assert.Error(t, err, "Should return error from source")
		assert.Contains(t, err.Error(), "read error", "Should contain read error message")
		assert.Zero(t, bytesWritten, "Should return zero bytes on error")
	})

	t.Run("Destination stream error", func(t *testing.T) {
		src := MemReader([]int{1, 2, 3}, nil)
		dst := MemWriter[int]()
		dst.SetError(errors.New("write error"))

		bytesWritten, err := Pipe(src, dst)
		assert.Error(t, err, "Should return error from destination")
		assert.Contains(t, err.Error(), "write error", "Should contain write error message")
		assert.Zero(t, bytesWritten, "Should return zero bytes on error")
	})

	t.Run("Successful pipe", func(t *testing.T) {
		src := MemReader([]int{1, 2, 3, 4, 5}, nil)
		dst := MemWriter[int]()

		bytesWritten, err := Pipe(src, dst)
		assert.NoError(t, err, "Should pipe successfully")
		assert.Positive(t, bytesWritten, "Should write positive bytes")

		result := dst.Items()
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result, "Should pipe all items")
	})
}

func TestMulticastErrorHandling(t *testing.T) {
	t.Run("No destinations", func(t *testing.T) {
		src := MemReader([]int{1, 2, 3}, nil)

		bytesWritten, err := Multicast(src)
		assert.NoError(t, err, "Should handle no destinations")
		assert.Empty(t, bytesWritten, "Should return empty bytes slice")
	})

	t.Run("Source stream error", func(t *testing.T) {
		src := MemReader([]int{1, 2}, errors.New("source error"))
		dst1 := MemWriter[int]()
		dst2 := MemWriter[int]()

		bytesWritten, err := Multicast(src, dst1, dst2)
		assert.Error(t, err, "Should return error from source")
		assert.Contains(t, err.Error(), "read error", "Should contain read error message")
		assert.Len(t, bytesWritten, 2, "Should return bytes array for all destinations")
	})

	t.Run("Destination stream error", func(t *testing.T) {
		src := MemReader([]int{1, 2, 3}, nil)
		dst1 := MemWriter[int]()
		dst2 := MemWriter[int]()
		dst2.SetError(errors.New("write error"))

		bytesWritten, err := Multicast(src, dst1, dst2)
		assert.Error(t, err, "Should return error from destination")
		assert.Contains(
			t,
			err.Error(),
			"write error to destination 1",
			"Should contain specific destination error",
		)
		assert.Len(t, bytesWritten, 2, "Should return bytes array for all destinations")
	})

	t.Run("Successful Multicast", func(t *testing.T) {
		src := MemReader([]int{1, 2, 3}, nil)
		dst1 := MemWriter[int]()
		dst2 := MemWriter[int]()
		dst3 := MemWriter[int]()

		bytesWritten, err := Multicast(src, dst1, dst2, dst3)
		assert.NoError(t, err, "Should Multicast successfully")
		assert.Len(t, bytesWritten, 3, "Should return bytes written for each destination")

		for i, bytes := range bytesWritten {
			assert.Positive(t, bytes, "Destination %d should have positive bytes written", i)
		}

		// Verify all destinations have the same data
		expected := []int{1, 2, 3}
		assert.Equal(t, expected, dst1.Items(), "Destination 1 should have all items")
		assert.Equal(t, expected, dst2.Items(), "Destination 2 should have all items")
		assert.Equal(t, expected, dst3.Items(), "Destination 3 should have all items")
	})
}

// Mock types for testing
type utilsTestMarshaler struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
