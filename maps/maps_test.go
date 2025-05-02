package maps

import (
	"strconv"
	"testing"

	"github.com/sonirico/gozo/fp"
	"github.com/sonirico/gozo/tuples"
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
			strconv.FormatInt(int64(k), 10),
			strconv.FormatInt(int64(v*v), 10),
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
