package streams

import (
	"fmt"
	"strings"
)

// ExampleReduce demonstrates reducing a stream to a map with aggregated values.
func ExampleReduce() {
	// Create a stream of words
	reader := strings.NewReader("apple\nbanana\napple\ncherry\nbanana\napple")
	stream := Lines(reader)

	// Count occurrences of each word
	counts, _ := Reduce(stream, func(acc map[string]int, word string) map[string]int {
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
