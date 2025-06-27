package streams

import (
	"testing"
)

func TestFilterStream(t *testing.T) {
	// Test filtering integers
	sourceData := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	source := NewMemory(sourceData, nil)

	// Filter even numbers
	evenFilter := func(n int) bool { return n%2 == 0 }
	filtered := NewFilterStream(source, evenFilter)

	result, err := Consume(filtered)
	if err != nil {
		t.Fatalf("Failed to consume filtered stream: %v", err)
	}

	expected := []int{2, 4, 6, 8, 10}
	if len(result) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(result))
	}

	for i, expectedVal := range expected {
		if result[i] != expectedVal {
			t.Errorf("Expected item %d to be %d, got %d", i, expectedVal, result[i])
		}
	}
}

func TestFilterStreamStrings(t *testing.T) {
	// Test filtering strings
	sourceData := []string{"hello", "world", "test", "filter", "stream"}
	source := NewMemory(sourceData, nil)

	// Filter strings with length > 4
	lengthFilter := func(s string) bool { return len(s) > 4 }
	filtered := NewFilterStream(source, lengthFilter)

	result, err := Consume(filtered)
	if err != nil {
		t.Fatalf("Failed to consume filtered stream: %v", err)
	}

	expected := []string{"hello", "world", "filter", "stream"}
	if len(result) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(result))
	}

	for i, expectedVal := range expected {
		if result[i] != expectedVal {
			t.Errorf("Expected item %d to be %s, got %s", i, expectedVal, result[i])
		}
	}
}

func TestFilterStreamEmpty(t *testing.T) {
	// Test filtering that results in empty stream
	sourceData := []int{1, 3, 5, 7, 9}
	source := NewMemory(sourceData, nil)

	// Filter even numbers (none exist)
	evenFilter := func(n int) bool { return n%2 == 0 }
	filtered := NewFilterStream(source, evenFilter)

	result, err := Consume(filtered)
	if err != nil {
		t.Fatalf("Failed to consume filtered stream: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d items", len(result))
	}
}

func TestFilterStreamAll(t *testing.T) {
	// Test filtering that keeps all items
	sourceData := []int{2, 4, 6, 8, 10}
	source := NewMemory(sourceData, nil)

	// Filter even numbers (all are even)
	evenFilter := func(n int) bool { return n%2 == 0 }
	filtered := NewFilterStream(source, evenFilter)

	result, err := Consume(filtered)
	if err != nil {
		t.Fatalf("Failed to consume filtered stream: %v", err)
	}

	if len(result) != len(sourceData) {
		t.Fatalf("Expected %d items, got %d", len(sourceData), len(result))
	}

	for i, expectedVal := range sourceData {
		if result[i] != expectedVal {
			t.Errorf("Expected item %d to be %d, got %d", i, expectedVal, result[i])
		}
	}
}

func TestFilterStreamChaining(t *testing.T) {
	// Test chaining filter with mapper
	sourceData := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	source := NewMemory(sourceData, nil)

	// First filter even numbers
	evenFilter := func(n int) bool { return n%2 == 0 }
	filtered := NewFilterStream(source, evenFilter)

	// Then map to double the values
	doubler := func(n int) int { return n * 2 }
	mapped := NewMapperStream(filtered, doubler)

	result, err := Consume(mapped)
	if err != nil {
		t.Fatalf("Failed to consume chained stream: %v", err)
	}

	expected := []int{4, 8, 12, 16, 20} // Even numbers doubled
	if len(result) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(result))
	}

	for i, expectedVal := range expected {
		if result[i] != expectedVal {
			t.Errorf("Expected item %d to be %d, got %d", i, expectedVal, result[i])
		}
	}
}

func TestFilterStreamInterface(t *testing.T) {
	// Test that FilterStream implements ReadStream interface
	var _ ReadStream[int] = &FilterStream[int]{}
}

func TestFilterStreamIterator(t *testing.T) {
	// Test using FilterStream with iterator
	sourceData := []int{1, 2, 3, 4, 5}
	source := NewMemory(sourceData, nil)

	// Filter odd numbers
	oddFilter := func(n int) bool { return n%2 == 1 }
	filtered := NewFilterStream(source, oddFilter).(*FilterStream[int])

	var result []int
	for value := range filtered.Iter() {
		result = append(result, value)
	}

	expected := []int{1, 3, 5}
	if len(result) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(result))
	}

	for i, expectedVal := range expected {
		if result[i] != expectedVal {
			t.Errorf("Expected item %d to be %d, got %d", i, expectedVal, result[i])
		}
	}
}
