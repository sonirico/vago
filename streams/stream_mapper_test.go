package streams

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ExampleMap demonstrates transforming elements in a stream.
func ExampleMap() {
	// Create a stream from a slice of integers
	data := []int{1, 2, 3, 4, 5}
	stream := MemReader(data, nil)

	// Transform integers to their string representation
	stringStream := Map(stream, func(n int) string {
		return fmt.Sprintf("number_%d", n)
	})

	// Collect the results
	result, _ := Consume(stringStream)
	fmt.Println(result)
	// Output: [number_1 number_2 number_3 number_4 number_5]
}

func TestMapperStream(t *testing.T) {
	// Test mapping integers to strings
	sourceData := []int{1, 2, 3, 4, 5}
	source := MemReader(sourceData, nil)

	// Map to strings
	mapped := Map(source, func(n int) string {
		return fmt.Sprintf("item_%d", n)
	})

	result, err := Consume(mapped)
	assert.NoError(t, err)

	expected := []string{"item_1", "item_2", "item_3", "item_4", "item_5"}
	assert.Equal(t, expected, result)
}
