package streams

import (
	"strconv"
	"testing"

	"github.com/sonirico/gozo/fp"
	"github.com/stretchr/testify/assert"
)

func TestFilterMapStream(t *testing.T) {
	// Test filtering integers
	sourceData := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	source := MemReader(sourceData, nil)

	// return even numbers as string
	evenFilterMap := func(n int) (string, bool) {
		if n%2 == 1 {
			return "", false // Filter out odd numbers
		}

		return strconv.FormatInt(int64(n), 10), true // Convert even numbers to string
	}
	filtered := FilterMap(source, evenFilterMap)

	result, err := Consume(filtered)
	if err != nil {
		t.Fatalf("Failed to consume filtered stream: %v", err)
	}

	expected := []string{"2", "4", "6", "8", "10"}

	assert.Lenf(t, result, len(expected), "Expected %d items, got %d", len(expected), len(result))
	assert.EqualValues(t, expected, result, "Filtered items do not match expected values")
}

func TestFilterMapOptStream(t *testing.T) {
	// Test filtering integers
	sourceData := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	source := MemReader(sourceData, nil)

	// return even numbers as string
	evenFilterMap := func(n int) fp.Option[string] {
		if n%2 == 1 {
			return fp.None[string]()
		}

		return fp.Some(strconv.FormatInt(int64(n), 10))
	}
	filtered := FilterMapOpt(source, evenFilterMap)

	result, err := Consume(filtered)
	if err != nil {
		t.Fatalf("Failed to consume filtered stream: %v", err)
	}

	expected := []string{"2", "4", "6", "8", "10"}

	assert.Lenf(t, result, len(expected), "Expected %d items, got %d", len(expected), len(result))
	assert.EqualValues(t, expected, result, "Filtered items do not match expected values")
}

func TestFilterMapOptStreamEmpty(t *testing.T) {
	// Test filtering integers
	sourceData := []int{1, 3, 5, 7, 9}
	source := MemReader(sourceData, nil)

	// return even numbers as string
	evenFilterMap := func(n int) fp.Option[string] {
		if n%2 == 1 {
			return fp.None[string]()
		}

		return fp.Some(strconv.FormatInt(int64(n), 10))
	}
	filtered := FilterMapOpt(source, evenFilterMap)

	result, err := Consume(filtered)
	if err != nil {
		t.Fatalf("Failed to consume filtered stream: %v", err)
	}

	assert.Len(t, result, 0, "Expected no items, got %d", len(result))
}
