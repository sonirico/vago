package maps

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/sonirico/vago/fp"
	"github.com/sonirico/vago/tuples"
)

func TestMapTo(t *testing.T) {
	type (
		testCase struct {
			name     string
			payload  map[int]int
			expected map[string]string
		}
	)

	tests := []testCase{
		{
			name:     "nil map is noop",
			payload:  nil,
			expected: nil,
		},
		{
			name:     "empty map returns empty map",
			payload:  map[int]int{},
			expected: map[string]string{},
		},
		{
			name:     "filled map",
			payload:  map[int]int{100: 3, 29: 2},
			expected: map[string]string{"100": "9", "29": "4"},
		},
	}

	predicate := func(k, v int) (string, string) {
		return strconv.FormatInt(int64(k), 10), strconv.FormatInt(int64(v*v), 10)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Map(test.payload, predicate)

			if !Equals(test.expected, actual, assertMapValueEq) {
				t.Errorf("unexpected map\nwant %v\nhave %v",
					test.expected, actual)
			}
		})
	}
}

func TestFilterMap(t *testing.T) {
	type (
		testCase struct {
			name     string
			payload  map[int]int
			expected map[string]string
		}
	)

	tests := []testCase{
		{
			name:     "nil map is noop",
			payload:  nil,
			expected: nil,
		},
		{
			name:     "empty map returns empty map",
			payload:  map[int]int{},
			expected: map[string]string{},
		},
		{
			name:     "filled map",
			payload:  map[int]int{101: 3, 22: 2},
			expected: map[string]string{"101": "9"},
		},
	}

	predicate := func(k, v int) fp.Option[tuples.Tuple2[string, string]] {
		if k%2 == 0 {
			return fp.None[tuples.Tuple2[string, string]]()
		}

		return fp.Some(tuples.Tuple2[string, string]{
			V1: strconv.FormatInt(int64(k), 10),
			V2: strconv.FormatInt(int64(v*v), 10),
		})
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := FilterMap(test.payload, predicate)

			if !Equals(test.expected, actual, assertMapValueEq) {
				t.Errorf("unexpected map\nwant %v\nhave %v",
					test.expected, actual)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	type (
		testCase struct {
			name     string
			payload  map[int]int
			expected map[int]int
		}
	)

	tests := []testCase{
		{
			name:     "nil map is noop",
			payload:  nil,
			expected: nil,
		},
		{
			name:     "empty map returns empty map",
			payload:  map[int]int{},
			expected: map[int]int{},
		},
		{
			name:     "filled map",
			payload:  map[int]int{101: 3, 22: 2},
			expected: map[int]int{22: 2},
		},
	}

	predicate := func(k, v int) bool {
		return k%2 == 0
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Filter(test.payload, predicate)

			if !Equals(test.expected, actual, func(x, y int) bool { return x == y }) {
				t.Errorf("unexpected map\nwant %v\nhave %v",
					test.expected, actual)
			}
		})
	}
}

func TestFilterInPlace(t *testing.T) {
	type (
		testCase struct {
			name     string
			payload  map[int]int
			expected map[int]int
		}
	)

	tests := []testCase{
		{
			name:     "nil map is noop",
			payload:  nil,
			expected: nil,
		},
		{
			name:     "empty map returns empty map",
			payload:  map[int]int{},
			expected: map[int]int{},
		},
		{
			name:     "filled map",
			payload:  map[int]int{101: 3, 22: 2},
			expected: map[int]int{22: 2},
		},
	}

	predicate := func(k, v int) bool {
		return k%2 == 0
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := FilterInPlace(test.payload, predicate)

			if !Equals(test.expected, actual, func(x, y int) bool { return x == y }) {
				t.Errorf("unexpected map\nwant %v\nhave %v",
					test.expected, actual)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	type (
		testCase struct {
			name     string
			payload  map[int]int
			expected int
		}
	)

	tests := []testCase{
		{
			name:     "nil map yields zero value",
			payload:  nil,
			expected: 0,
		},
		{
			name:     "empty map returns zero value",
			payload:  map[int]int{},
			expected: 0,
		},
		{
			name:     "filled map",
			payload:  map[int]int{101: 3, 22: 2},
			expected: 128,
		},
	}

	predicate := func(acc, k, v int) int {
		return acc + k + v
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Reduce(test.payload, predicate)

			if test.expected != actual {
				t.Errorf("unexpected map reduce result. \nwant %v\nhave %v",
					test.expected, actual)
			}
		})
	}
}

func TestFold(t *testing.T) {
	type (
		testCase struct {
			name     string
			payload  map[int]int
			expected int
		}
	)

	tests := []testCase{
		{
			name:     "nil map yields initial",
			payload:  nil,
			expected: 1,
		},
		{
			name:     "empty map returns initial value",
			payload:  map[int]int{},
			expected: 1,
		},
		{
			name:     "filled map",
			payload:  map[int]int{101: 3, 22: 2},
			expected: 129,
		},
	}

	predicate := func(acc, k, v int) int {
		return acc + k + v
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Fold(test.payload, predicate, 1)

			if test.expected != actual {
				t.Errorf("unexpected map reduce result. \nwant %v\nhave %v",
					test.expected, actual)
			}
		})
	}
}

func assertMapValueEq(x, y string) bool {
	return x == y
}

// ExampleMap demonstrates transforming keys and values in a map.
func ExampleMap() {
	// Create a map of numbers to their names
	numbers := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	// Transform to string keys and uppercase values
	transformed := Map(numbers, func(key int, value string) (string, string) {
		return fmt.Sprintf("num_%d", key), strings.ToUpper(value)
	})

	fmt.Println(transformed["num_1"])
	// Output: ONE
}

// ExampleFilter demonstrates filtering a map by key-value pairs.
func ExampleFilter() {
	// Create a map of products to prices
	prices := map[string]int{
		"apple":  100,
		"banana": 50,
		"cherry": 200,
		"date":   75,
	}

	// Keep only items that cost more than 75
	expensive := Filter(prices, func(product string, price int) bool {
		return price > 75
	})

	fmt.Printf("Expensive items count: %d\n", len(expensive))
	// Output: Expensive items count: 2
}

// ExampleFilterMap demonstrates filtering and transforming in a single operation.
func ExampleFilterMap() {
	// Create a map of names to ages
	ages := map[string]int{
		"Alice": 25,
		"Bob":   17,
		"Carol": 30,
		"Dave":  16,
	}

	// Keep only adults and transform to ID format
	adults := FilterMap(ages, func(name string, age int) fp.Option[tuples.Tuple2[string, string]] {
		if age >= 18 {
			id := fmt.Sprintf("ID_%s_%d", name, age)
			return fp.Some(tuples.Tuple2[string, string]{V1: name, V2: id})
		}
		return fp.None[tuples.Tuple2[string, string]]()
	})

	fmt.Printf("Adult count: %d\n", len(adults))
	// Output: Adult count: 2
}

// ExampleReduce demonstrates reducing a map to a single value.
func ExampleReduce() {
	// Create a map of item quantities
	inventory := map[string]int{
		"apples":  10,
		"bananas": 5,
		"oranges": 8,
	}

	// Calculate total items (Reduce starts with zero value)
	total := Reduce(inventory, func(acc int, key string, value int) int {
		return acc + value
	})

	fmt.Printf("Total items: %d\n", total)
	// Output: Total items: 23
}

// ExampleEquals demonstrates comparing two maps for equality.
func ExampleEquals() {
	// Create two maps
	map1 := map[string]int{"a": 1, "b": 2, "c": 3}
	map2 := map[string]int{"a": 1, "b": 2, "c": 3}
	map3 := map[string]int{"a": 1, "b": 2, "c": 4}

	// Compare using equality function
	equal1 := Equals(map1, map2, func(x, y int) bool { return x == y })
	equal2 := Equals(map1, map3, func(x, y int) bool { return x == y })

	fmt.Printf("map1 == map2: %t\n", equal1)
	fmt.Printf("map1 == map3: %t\n", equal2)
	// Output:
	// map1 == map2: true
	// map1 == map3: false
}

// ExampleFilterMapTuple demonstrates filtering and transforming using tuple returns.
func ExampleFilterMapTuple() {
	// Create a map of scores
	scores := map[string]int{
		"Alice": 85,
		"Bob":   70,
		"Carol": 95,
		"Dave":  60,
	}

	// Keep high scores and convert to grade format
	grades := FilterMapTuple(scores, func(name string, score int) (string, string, bool) {
		if score >= 80 {
			var grade string
			if score >= 90 {
				grade = "A"
			} else {
				grade = "B"
			}
			return name, grade, true
		}
		return "", "", false
	})

	fmt.Printf("High performers: %d\n", len(grades))
	fmt.Printf("Alice's grade: %s\n", grades["Alice"])
	// Output:
	// High performers: 2
	// Alice's grade: B
}

// ExampleFold demonstrates folding a map with an initial value.
func ExampleFold() {
	// Create a map of item prices
	prices := map[string]float64{
		"apple":  1.20,
		"banana": 0.80,
		"cherry": 2.50,
	}

	// Calculate total with initial tax
	totalWithTax := Fold(prices, func(acc float64, item string, price float64) float64 {
		return acc + price*1.1 // Add 10% tax
	}, 5.0) // Start with 5.0 base fee

	fmt.Printf("Total with tax: %.2f\n", totalWithTax)
	// Output:
	// Total with tax: 9.95
}

// ExampleSlice demonstrates converting a map to a slice.
func ExampleSlice() {
	// Create a map of user data
	users := map[int]string{
		1: "Alice",
		2: "Bob",
		3: "Carol",
	}

	// Convert to slice of formatted strings
	userList := Slice(users, func(id int, name string) string {
		return fmt.Sprintf("ID:%d Name:%s", id, name)
	})

	fmt.Printf("Users count: %d\n", len(userList))
	// Note: map iteration order is not guaranteed
	// Output:
	// Users count: 3
}

func TestKeys(t *testing.T) {
	type testCase struct {
		name     string
		payload  map[string]int
		expected int // expected length
	}

	tests := []testCase{
		{
			name:     "nil map returns nil slice",
			payload:  nil,
			expected: 0,
		},
		{
			name:     "empty map returns empty slice",
			payload:  map[string]int{},
			expected: 0,
		},
		{
			name:     "filled map returns all keys",
			payload:  map[string]int{"apple": 1, "banana": 2, "cherry": 3},
			expected: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Keys(test.payload)

			if test.payload == nil && actual != nil {
				t.Errorf("expected nil slice for nil map, got %v", actual)
				return
			}

			if len(actual) != test.expected {
				t.Errorf("unexpected slice length\nwant %d\nhave %d",
					test.expected, len(actual))
				return
			}

			// Verify all keys are present (order may vary)
			if test.payload != nil {
				for expectedKey := range test.payload {
					found := false
					for _, actualKey := range actual {
						if actualKey == expectedKey {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("missing key %v in result", expectedKey)
					}
				}
			}
		})
	}
}

func TestValues(t *testing.T) {
	type testCase struct {
		name     string
		payload  map[string]int
		expected int // expected length
	}

	tests := []testCase{
		{
			name:     "nil map returns nil slice",
			payload:  nil,
			expected: 0,
		},
		{
			name:     "empty map returns empty slice",
			payload:  map[string]int{},
			expected: 0,
		},
		{
			name:     "filled map returns all values",
			payload:  map[string]int{"apple": 1, "banana": 2, "cherry": 3},
			expected: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Values(test.payload)

			if test.payload == nil && actual != nil {
				t.Errorf("expected nil slice for nil map, got %v", actual)
				return
			}

			if len(actual) != test.expected {
				t.Errorf("unexpected slice length\nwant %d\nhave %d",
					test.expected, len(actual))
				return
			}

			// Verify all values are present (order may vary)
			if test.payload != nil {
				for _, expectedValue := range test.payload {
					found := false
					for _, actualValue := range actual {
						if actualValue == expectedValue {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("missing value %v in result", expectedValue)
					}
				}
			}
		})
	}
}

func TestSeqKeys(t *testing.T) {
	type testCase struct {
		name     string
		payload  map[string]int
		expected int // expected count
	}

	tests := []testCase{
		{
			name:     "nil map yields no keys",
			payload:  nil,
			expected: 0,
		},
		{
			name:     "empty map yields no keys",
			payload:  map[string]int{},
			expected: 0,
		},
		{
			name:     "filled map yields all keys",
			payload:  map[string]int{"apple": 1, "banana": 2, "cherry": 3},
			expected: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			seq := SeqKeys(test.payload)

			var actualKeys []string
			for key := range seq {
				actualKeys = append(actualKeys, key)
			}

			if len(actualKeys) != test.expected {
				t.Errorf("unexpected key count\nwant %d\nhave %d",
					test.expected, len(actualKeys))
				return
			}

			// Verify all keys are present (order may vary)
			if test.payload != nil {
				for expectedKey := range test.payload {
					found := false
					for _, actualKey := range actualKeys {
						if actualKey == expectedKey {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("missing key %v in result", expectedKey)
					}
				}
			}
		})
	}
}

func TestSeqValues(t *testing.T) {
	type testCase struct {
		name     string
		payload  map[string]int
		expected int // expected count
	}

	tests := []testCase{
		{
			name:     "nil map yields no values",
			payload:  nil,
			expected: 0,
		},
		{
			name:     "empty map yields no values",
			payload:  map[string]int{},
			expected: 0,
		},
		{
			name:     "filled map yields all values",
			payload:  map[string]int{"apple": 1, "banana": 2, "cherry": 3},
			expected: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			seq := SeqValues(test.payload)

			var actualValues []int
			for value := range seq {
				actualValues = append(actualValues, value)
			}

			if len(actualValues) != test.expected {
				t.Errorf("unexpected value count\nwant %d\nhave %d",
					test.expected, len(actualValues))
				return
			}

			// Verify all values are present (order may vary)
			if test.payload != nil {
				for _, expectedValue := range test.payload {
					found := false
					for _, actualValue := range actualValues {
						if actualValue == expectedValue {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("missing value %v in result", expectedValue)
					}
				}
			}
		})
	}
}

// ExampleKeys demonstrates extracting all keys from a map.
func ExampleKeys() {
	// Create a map of products to prices
	prices := map[string]int{
		"apple":  100,
		"banana": 50,
		"cherry": 200,
	}

	// Extract all product names
	products := Keys(prices)

	fmt.Printf("Product count: %d\n", len(products))
	// Note: map iteration order is not guaranteed
	// Output:
	// Product count: 3
}

// ExampleValues demonstrates extracting all values from a map.
func ExampleValues() {
	// Create a map of products to prices
	prices := map[string]int{
		"apple":  100,
		"banana": 50,
		"cherry": 200,
	}

	// Extract all prices
	allPrices := Values(prices)

	fmt.Printf("Price count: %d\n", len(allPrices))
	// Note: map iteration order is not guaranteed
	// Output:
	// Price count: 3
}

// ExampleSeqKeys demonstrates iterating over map keys using iterator pattern.
func ExampleSeqKeys() {
	// Create a map of products to prices
	prices := map[string]int{
		"apple":  100,
		"banana": 50,
		"cherry": 200,
	}

	// Iterate over keys lazily
	keyCount := 0
	for range SeqKeys(prices) {
		keyCount++
	}

	fmt.Printf("Key count: %d\n", keyCount)
	// Output:
	// Key count: 3
}

// ExampleSeqValues demonstrates iterating over map values using iterator pattern.
func ExampleSeqValues() {
	// Create a map of products to prices
	prices := map[string]int{
		"apple":  100,
		"banana": 50,
		"cherry": 200,
	}

	// Calculate total using iterator
	total := 0
	for price := range SeqValues(prices) {
		total += price
	}

	fmt.Printf("Total price: %d\n", total)
	// Output:
	// Total price: 350
}
