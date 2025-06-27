package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompactStream(t *testing.T) {
	t.Run("groups consecutive items with same key", func(t *testing.T) {
		// Arrange
		data := []string{"a", "a", "b", "b", "b", "a", "c"}
		memStream := MemReader(data, nil)

		// Use identity function as key extractor
		compactStream := Compact(memStream, func(s string) string { return s })

		// Act & Assert

		// First group: ["a", "a"]
		require.True(t, compactStream.Next())
		assert.Equal(t, []string{"a", "a"}, compactStream.Data())
		assert.NoError(t, compactStream.Err())

		// Second group: ["b", "b", "b"]
		require.True(t, compactStream.Next())
		assert.Equal(t, []string{"b", "b", "b"}, compactStream.Data())
		assert.NoError(t, compactStream.Err())

		// Third group: ["a"] (different from first group because not consecutive)
		require.True(t, compactStream.Next())
		assert.Equal(t, []string{"a"}, compactStream.Data())
		assert.NoError(t, compactStream.Err())

		// Fourth group: ["c"]
		require.True(t, compactStream.Next())
		assert.Equal(t, []string{"c"}, compactStream.Data())
		assert.NoError(t, compactStream.Err())

		// No more groups
		assert.False(t, compactStream.Next())
		assert.NoError(t, compactStream.Err())

		// Cleanup
		assert.NoError(t, compactStream.Close())
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
		compactStream := Compact(memStream, func(p Person) int { return p.Age })

		// Act & Assert

		// First group: people aged 25
		require.True(t, compactStream.Next())
		group1 := compactStream.Data()
		assert.Len(t, group1, 3)
		assert.Equal(t, 25, group1[0].Age)
		assert.Equal(t, 25, group1[1].Age)
		assert.Equal(t, 25, group1[2].Age)
		assert.NoError(t, compactStream.Err())

		// Second group: people aged 30
		require.True(t, compactStream.Next())
		group2 := compactStream.Data()
		assert.Len(t, group2, 2)
		assert.Equal(t, 30, group2[0].Age)
		assert.Equal(t, 30, group2[1].Age)
		assert.NoError(t, compactStream.Err())

		// Third group: Frank aged 25 (separate group)
		require.True(t, compactStream.Next())
		group3 := compactStream.Data()
		assert.Len(t, group3, 1)
		assert.Equal(t, 25, group3[0].Age)
		assert.Equal(t, "Frank", group3[0].Name)
		assert.NoError(t, compactStream.Err())

		// No more groups
		assert.False(t, compactStream.Next())
		assert.NoError(t, compactStream.Err())

		// Cleanup
		assert.NoError(t, compactStream.Close())
	})

	t.Run("handles empty stream", func(t *testing.T) {
		// Arrange
		var data []string
		memStream := MemReader(data, nil)
		compactStream := Compact(memStream, func(s string) string { return s })

		// Act & Assert
		assert.False(t, compactStream.Next())
		assert.NoError(t, compactStream.Err())
		assert.NoError(t, compactStream.Close())
	})

	t.Run("handles single item", func(t *testing.T) {
		// Arrange
		data := []string{"single"}
		memStream := MemReader(data, nil)
		compactStream := Compact(memStream, func(s string) string { return s })

		// Act & Assert
		require.True(t, compactStream.Next())
		assert.Equal(t, []string{"single"}, compactStream.Data())
		assert.NoError(t, compactStream.Err())

		assert.False(t, compactStream.Next())
		assert.NoError(t, compactStream.Err())
		assert.NoError(t, compactStream.Close())
	})

	t.Run("works with iterator", func(t *testing.T) {
		// Arrange
		data := []int{1, 1, 2, 2, 3, 1}
		memStream := MemReader(data, nil)
		compactStream := Compact(memStream, func(i int) int { return i })

		// Act
		var groups [][]int
		for group := range Iter(compactStream) {
			// Make a copy to avoid slice aliasing issues
			groupCopy := make([]int, len(group))
			copy(groupCopy, group)
			groups = append(groups, groupCopy)
		}

		// Assert
		expected := [][]int{
			{1, 1},
			{2, 2},
			{3},
			{1},
		}
		assert.Equal(t, expected, groups)
		assert.NoError(t, compactStream.Close())
	})
}
