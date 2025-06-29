package streams

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupStream(t *testing.T) {
	t.Run("groups consecutive items with same key", func(t *testing.T) {
		// Arrange
		data := []string{"a", "a", "b", "b", "b", "a", "c"}
		memStream := MemReader(data, nil)

		// Use identity function as key extractor
		GroupStream := Group(memStream, func(s string) string { return s })

		// Act & Assert

		// First group: ["a", "a"]
		require.True(t, GroupStream.Next())
		assert.Equal(t, []string{"a", "a"}, GroupStream.Data())
		assert.NoError(t, GroupStream.Err())

		// Second group: ["b", "b", "b"]
		require.True(t, GroupStream.Next())
		assert.Equal(t, []string{"b", "b", "b"}, GroupStream.Data())
		assert.NoError(t, GroupStream.Err())

		// Third group: ["a"] (different from first group because not consecutive)
		require.True(t, GroupStream.Next())
		assert.Equal(t, []string{"a"}, GroupStream.Data())
		assert.NoError(t, GroupStream.Err())

		// Fourth group: ["c"]
		require.True(t, GroupStream.Next())
		assert.Equal(t, []string{"c"}, GroupStream.Data())
		assert.NoError(t, GroupStream.Err())

		// No more groups
		assert.False(t, GroupStream.Next())
		assert.NoError(t, GroupStream.Err())

		// Cleanup
		assert.NoError(t, GroupStream.Close())
	})

	t.Run("works with custom key extractor", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name string
			Age  int
		}

		data := []Person{
			{"Alice", 25},
			{"Bob", 25},
			{"Charlie", 25},
			{"David", 30},
			{"Eve", 30},
			{"Frank", 25}, // Different age group, not consecutive with first age 25 group
		}

		memStream := MemReader(data, nil)

		// Group by age
		GroupStream := Group(memStream, func(p Person) int { return p.Age })

		// Act & Assert

		// First group: people aged 25
		require.True(t, GroupStream.Next())
		group1 := GroupStream.Data()
		assert.Len(t, group1, 3)
		assert.Equal(t, 25, group1[0].Age)
		assert.Equal(t, 25, group1[1].Age)
		assert.Equal(t, 25, group1[2].Age)
		assert.NoError(t, GroupStream.Err())

		// Second group: people aged 30
		require.True(t, GroupStream.Next())
		group2 := GroupStream.Data()
		assert.Len(t, group2, 2)
		assert.Equal(t, 30, group2[0].Age)
		assert.Equal(t, 30, group2[1].Age)
		assert.NoError(t, GroupStream.Err())

		// Third group: Frank aged 25 (separate group)
		require.True(t, GroupStream.Next())
		group3 := GroupStream.Data()
		assert.Len(t, group3, 1)
		assert.Equal(t, 25, group3[0].Age)
		assert.Equal(t, "Frank", group3[0].Name)
		assert.NoError(t, GroupStream.Err())

		// No more groups
		assert.False(t, GroupStream.Next())
		assert.NoError(t, GroupStream.Err())

		// Cleanup
		assert.NoError(t, GroupStream.Close())
	})

	t.Run("handles empty stream", func(t *testing.T) {
		// Arrange
		var data []string
		memStream := MemReader(data, nil)
		GroupStream := Group(memStream, func(s string) string { return s })

		// Act & Assert
		assert.False(t, GroupStream.Next())
		assert.NoError(t, GroupStream.Err())
		assert.NoError(t, GroupStream.Close())
	})

	t.Run("handles single item", func(t *testing.T) {
		// Arrange
		data := []string{"single"}
		memStream := MemReader(data, nil)
		GroupStream := Group(memStream, func(s string) string { return s })

		// Act & Assert
		require.True(t, GroupStream.Next())
		assert.Equal(t, []string{"single"}, GroupStream.Data())
		assert.NoError(t, GroupStream.Err())

		assert.False(t, GroupStream.Next())
		assert.NoError(t, GroupStream.Err())
		assert.NoError(t, GroupStream.Close())
	})

	t.Run("works with iterator", func(t *testing.T) {
		// Arrange
		data := []int{1, 1, 2, 2, 3, 1}
		memStream := MemReader(data, nil)
		GroupStream := Group(memStream, func(i int) int { return i })

		// Act
		var groups [][]int
		for group := range Iter(GroupStream) {
			groups = append(groups, group)
		}

		// Assert
		expected := [][]int{
			{1, 1},
			{2, 2},
			{3},
			{1},
		}
		assert.Equal(t, expected, groups)
		assert.NoError(t, GroupStream.Close())
	})

	t.Run("works with Consume", func(t *testing.T) {
		// Arrange
		data := []string{"apple", "apricot", "banana", "blueberry", "cherry", "coconut"}
		memStream := MemReader(data, nil)
		GroupStream := Group(memStream, func(s string) rune {
			return rune(s[0])
		})

		// Act
		groups, err := Consume(GroupStream)

		// Assert
		require.NoError(t, err)
		expected := [][]string{
			{"apple", "apricot"},
			{"banana", "blueberry"},
			{"cherry", "coconut"},
		}
		assert.Equal(t, expected, groups)
		assert.NoError(t, GroupStream.Close())
	})
}

// ExampleGroup demonstrates grouping consecutive items with the same key.
func ExampleGroup() {
	// Create a stream from a slice of strings
	data := []string{"apple", "apricot", "banana", "blueberry", "cherry", "coconut"}
	stream := MemReader(data, nil)

	// Group by first letter
	Grouped := Group(stream, func(s string) rune {
		return rune(s[0])
	})

	// Collect the results
	result, _ := Consume(Grouped)
	for _, group := range result {
		fmt.Printf("Group: %v\n", group)
	}
	// Output:
	// Group: [apple apricot]
	// Group: [banana blueberry]
	// Group: [cherry coconut]
}
