package slices

import (
	"fmt"
	"testing"

	"github.com/sonirico/gozo/fp"
)

func TestSlice_Len(t *testing.T) {
	type testCase struct {
		name           string
		payload        Slice[int]
		expectedLength int
	}

	tests := []testCase{
		{
			name:           "zero length slice",
			payload:        Slice[int]([]int{}),
			expectedLength: 0,
		},
		{
			name:           "nil slice",
			payload:        Slice[int](nil),
			expectedLength: 0,
		},
		{
			name:           "slice with more than one",
			payload:        Slice[int]([]int{1, 2, 3}),
			expectedLength: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectedLength != test.payload.Len() {
				t.Errorf("unexpected slice length. want %d, have %d",
					test.expectedLength, test.payload.Len())
			}
		})
	}
}

func TestSlice_Range(t *testing.T) {
	type testCase struct {
		name           string
		payload        Slice[int]
		expectedLength int
	}

	tests := []testCase{
		{
			name:           "zero length slice",
			payload:        Slice[int]([]int{}),
			expectedLength: 0,
		},
		{
			name:           "nil slice",
			payload:        Slice[int](nil),
			expectedLength: 0,
		},
		{
			name:           "slice with more than one",
			payload:        Slice[int]([]int{1, 2, 3}),
			expectedLength: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualLen := 0
			test.payload.Range(func(x int, _ int) bool {
				actualLen += 1
				return true
			})
			if test.expectedLength != actualLen {
				t.Errorf("unexpected slice length. want %d, have %d",
					test.expectedLength, actualLen)
			}
		})
	}
}

func TestSlice_Range_EarlyReturn(t *testing.T) {
	slice := Slice[int]([]int{1, 2, 3})
	actualLen := 0
	expectedLength := 2
	slice.Range(func(x int, i int) bool {
		actualLen += 1
		return i%2 == 0
	})
	if actualLen != expectedLength {
		t.Errorf("unexpected length. want %d, have %d", expectedLength, actualLen)
	}
}

func TestSlice_Get(t *testing.T) {
	type testCase struct {
		name        string
		payload     Slice[int]
		index       int
		expectedOk  bool
		expectedRes int
	}

	tests := []testCase{
		{
			name:        "negative index",
			payload:     Slice[int]([]int{}),
			index:       -1,
			expectedRes: 0,
			expectedOk:  false,
		},
		{
			name:        "zero index",
			payload:     Slice[int]([]int{1, 2, 3}),
			index:       0,
			expectedRes: 1,
			expectedOk:  true,
		},
		{
			name:        "zero index for nil slice",
			payload:     Slice[int](nil),
			index:       0,
			expectedRes: 0,
			expectedOk:  false,
		},
		{
			name:        "index of last item",
			payload:     Slice[int]([]int{1, 2, 3}),
			index:       2,
			expectedRes: 3,
			expectedOk:  true,
		},
		{
			name:        "out of bounds index",
			payload:     Slice[int]([]int{1, 2, 3}),
			index:       3,
			expectedRes: 0,
			expectedOk:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			actualRes, actualOk := test.payload.Get(test.index)

			if test.expectedOk != actualOk {
				t.Errorf("unexpected ok, want %t, have %t", test.expectedOk, actualOk)
			}
			if test.expectedRes != actualRes {
				t.Errorf("unexpected value, want %d, have %d", test.expectedRes, actualRes)
			}
		})
	}
}

func TestSlice_Append(t *testing.T) {
	numbers := Slice[int]([]int{1, 2, 3})
	numbers.Append(4)
	expectedLength := 4
	if numbers.Len() != expectedLength {
		t.Errorf("unexpected slice length, want %d, have %d",
			expectedLength, numbers.Len())
	}

	numbers.AppendVector([]int{5, 6})
	expectedLength = 6
	if numbers.Len() != expectedLength {
		t.Errorf("unexpected slice length, want %d, have %d",
			expectedLength, numbers.Len())
	}

	numbers.Push(7)
	expectedLength = 7
	if numbers.Len() != expectedLength {
		t.Errorf("unexpected slice length, want %d, have %d",
			expectedLength, numbers.Len())
	}
}

func TestSlice_IndexOf(t *testing.T) {
	type testCase struct {
		name        string
		payload     Slice[int]
		predicate   func(i int) bool
		expectedIdx int
	}

	tests := []testCase{
		{
			name:    "nil slice should return -1",
			payload: Slice[int]([]int{}),
			predicate: func(i int) bool {
				return true
			},
			expectedIdx: -1,
		},
		{
			name:    "item at the first position",
			payload: Slice[int]([]int{1, 2, 3}),
			predicate: func(i int) bool {
				return i == 1
			},
			expectedIdx: 0,
		},
		{
			name:    "item at the last position",
			payload: Slice[int]([]int{1, 2, 3}),
			predicate: func(i int) bool {
				return 3 == i
			},
			expectedIdx: 2,
		},
		{
			name:    "item not found",
			payload: Slice[int]([]int{73, 30, 5}),
			predicate: func(i int) bool {
				return 42 == i
			},
			expectedIdx: -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualIdx := test.payload.IndexOf(test.predicate)

			if test.expectedIdx != actualIdx {
				t.Errorf("unexpected value, want %d, have %d", test.expectedIdx, actualIdx)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type testCase struct {
		name     string
		payload  Slice[int]
		target   int
		expected bool
	}

	tests := []testCase{
		{
			name:     "nil slice should return false",
			payload:  Slice[int]([]int{}),
			target:   1,
			expected: false,
		},
		{
			name:     "item at the first position",
			payload:  Slice[int]([]int{1, 2, 3}),
			target:   1,
			expected: true,
		},
		{
			name:     "item at the last position",
			payload:  Slice[int]([]int{1, 2, 3}),
			target:   3,
			expected: true,
		},
		{
			name:     "item not found",
			payload:  Slice[int]([]int{73, 30, 5}),
			target:   3,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualContains := Contains(test.payload, func(x int) bool { return x == test.target })

			if test.expected != actualContains {
				t.Errorf("unexpected value, want %t, have %t", test.expected, actualContains)
			}
		})
	}
}

func TestIncludes(t *testing.T) {
	type testCase struct {
		name     string
		payload  Slice[int]
		target   int
		expected bool
	}

	tests := []testCase{
		{
			name:     "nil slice should return false",
			payload:  Slice[int]([]int{}),
			target:   1,
			expected: false,
		},
		{
			name:     "item at the first position",
			payload:  Slice[int]([]int{1, 2, 3}),
			target:   1,
			expected: true,
		},
		{
			name:     "item at the last position",
			payload:  Slice[int]([]int{1, 2, 3}),
			target:   3,
			expected: true,
		},
		{
			name:     "item not found",
			payload:  Slice[int]([]int{73, 30, 5}),
			target:   3,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualIncludes := Includes(test.payload, test.target)

			if test.expected != actualIncludes {
				t.Errorf("unexpected value, want %t, have %t", test.expected, actualIncludes)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	type testCase struct {
		name      string
		payload   Slice[int]
		expected  Slice[int]
		predicate func(int) bool
	}

	tests := []testCase{
		{
			name:     "nil slice should return nil slice",
			payload:  Slice[int]([]int{}),
			expected: Slice[int]([]int{}),
			predicate: func(i int) bool {
				return true
			},
		},
		{
			name:     "elements are filtered leaving some",
			payload:  Slice[int]([]int{1, 2, 3}),
			expected: Slice[int]([]int{2}),
			predicate: func(i int) bool {
				return i%2 == 0
			},
		},
		{
			name:     "elements are filtered leaving none",
			payload:  Slice[int]([]int{1, 2, 3}),
			expected: Slice[int]([]int{}),
			predicate: func(i int) bool {
				return i > 10
			},
		},
	}

	for _, test := range tests {
		t.Run("[Filter] "+test.name, func(t *testing.T) {
			actual := Filter(test.payload.Clone(), test.predicate)

			if !test.expected.Equals(actual, func(x, y int) bool {
				return x == y
			}) {
				t.Errorf("unexpected value, want %v, have %v", test.expected, actual)
			}
		})

		t.Run("[FilterInPlace] "+test.name, func(t *testing.T) {
			actual := FilterInPlace(test.payload.Clone(), test.predicate)

			if !test.expected.Equals(actual, func(x, y int) bool {
				return x == y
			}) {
				t.Errorf("unexpected value, want %v, have %v", test.expected, actual)
			}
		})
	}
}

func TestFilterMap(t *testing.T) {
	type testCase struct {
		name      string
		payload   Slice[int]
		expected  Slice[int]
		predicate func(int) fp.Option[int]
	}

	tests := []testCase{
		{
			name:     "nil slice should return nil slice",
			payload:  Slice[int]([]int{}),
			expected: Slice[int]([]int{}),
			predicate: func(i int) fp.Option[int] {
				return fp.None[int]()
			},
		},
		{
			name:     "elements are filtered leaving some",
			payload:  Slice[int]([]int{1, 2, 3}),
			expected: Slice[int]([]int{4}),
			predicate: func(i int) fp.Option[int] {
				if i%2 == 0 {
					return fp.Some(i * i)
				}
				return fp.None[int]()
			},
		},
		{
			name:     "elements are filtered leaving none",
			payload:  Slice[int]([]int{1, 2, 3}),
			expected: Slice[int]([]int{}),
			predicate: func(i int) fp.Option[int] {
				return fp.None[int]()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := FilterMap(test.payload, test.predicate)

			if !test.expected.Equals(actual, func(x, y int) bool {
				return x == y
			}) {
				t.Errorf("unexpected value, want %v, have %v", test.expected, actual)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	type testCase struct {
		name     string
		payload  Slice[int]
		expected int
	}

	predicate := func(x, y int) int { return x + y }

	tests := []testCase{
		{
			name:     "nil slice should return zero value",
			payload:  Slice[int]([]int{}),
			expected: 0,
		},
		{
			name:     "slice with only 1 element",
			payload:  Slice[int]([]int{1}),
			expected: 1,
		},
		{
			name:     "slice with several elements",
			payload:  Slice[int]([]int{1, 2, 3}),
			expected: 6,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := ReduceSame(test.payload, predicate)

			if test.expected != actual {
				t.Errorf("unexpected value, want %d, have %d", test.expected, actual)
			}
		})
	}
}

func TestCut(t *testing.T) {
	type testCase struct {
		name     string
		payload  Slice[int]
		from     int
		to       int
		expected Slice[int]
	}

	tests := []testCase{
		{
			name:     "nil slice should be noop",
			payload:  Slice[int]([]int{}),
			from:     0,
			to:       0,
			expected: Slice[int]([]int{}),
		},
		{
			name:     "slice with one item",
			payload:  Slice[int]([]int{1}),
			from:     0,
			to:       0,
			expected: Slice[int]([]int{}),
		},
		{
			name:     "slice with two items cut first one",
			payload:  Slice[int]([]int{1, 2}),
			from:     0,
			to:       0,
			expected: Slice[int]([]int{2}),
		},
		{
			name:     "slice with two items cut last one",
			payload:  Slice[int]([]int{1, 2}),
			from:     1,
			to:       1,
			expected: Slice[int]([]int{1}),
		},
		{
			name:     "slice with two items cut all",
			payload:  Slice[int]([]int{1, 2}),
			from:     0,
			to:       1,
			expected: Slice[int]([]int{}),
		},
		{
			name:     "`from` greater than to should consider `to` to be amount",
			payload:  Slice[int]([]int{1, 2}),
			from:     1,
			to:       0,
			expected: Slice[int]([]int{1}),
		},
		{
			name:     "`from` greater than slice length is moved to end",
			payload:  Slice[int]([]int{1, 2}),
			from:     3,
			to:       0,
			expected: Slice[int]([]int{1}),
		},
		{
			name:     "`to` greater than slice length is moved to end",
			payload:  Slice[int]([]int{1, 2}),
			from:     0,
			to:       3,
			expected: Slice[int]([]int{}),
		},
		{
			name:     "`from` lower than zero is moved to zero",
			payload:  Slice[int]([]int{1, 2}),
			from:     -1,
			to:       0,
			expected: Slice[int]([]int{2}),
		},
		{
			name:     "`to` lower than zero is moved to zero",
			payload:  Slice[int]([]int{1, 2}),
			from:     0,
			to:       -1,
			expected: Slice[int]([]int{2}),
		},
		{
			name:     "cut with more than two items cut all",
			payload:  Slice[int]([]int{1, 2, 3, 4}),
			from:     0,
			to:       4,
			expected: Slice[int]([]int{}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Cut(test.payload, test.from, test.to)

			if !test.expected.Equals(actual, func(x, y int) bool { return x == y }) {
				t.Errorf("unexpected value, want %v, have %v", test.expected, actual)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type testCase struct {
		name     string
		payload  Slice[int]
		idx      int
		expected Slice[int]
	}

	tests := []testCase{
		{
			name:     "nil slice should be noop",
			payload:  Slice[int]([]int{}),
			idx:      0,
			expected: Slice[int]([]int{}),
		},
		{
			name:     "slice with one item",
			payload:  Slice[int]([]int{1}),
			idx:      0,
			expected: Slice[int]([]int{}),
		},
		{
			name:     "slice with one item (idx lower than 0)",
			payload:  Slice[int]([]int{1}),
			idx:      -1,
			expected: Slice[int]([]int{1}),
		},
		{
			name:     "slice with one item (idx greater than length)",
			payload:  Slice[int]([]int{1}),
			idx:      3,
			expected: Slice[int]([]int{1}),
		},
		{
			name:     "slice with two elements",
			payload:  Slice[int]([]int{1, 2}),
			idx:      0,
			expected: Slice[int]([]int{2}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Delete(test.payload, test.idx)

			if !test.expected.Equals(actual, testArrEq) {
				t.Errorf("unexpected value, want %v, have %v", test.expected, actual)
			}
		})
	}
}

func TestDeleteOrder(t *testing.T) {
	type testCase struct {
		name     string
		payload  Slice[int]
		idx      int
		expected Slice[int]
	}

	tests := []testCase{
		{
			name:     "nil slice should be noop",
			payload:  Slice[int]([]int{}),
			idx:      0,
			expected: Slice[int]([]int{}),
		},
		{
			name:     "slice with one item",
			payload:  Slice[int]([]int{1}),
			idx:      0,
			expected: Slice[int]([]int{}),
		},
		{
			name:     "slice with one item (idx lower than 0)",
			payload:  Slice[int]([]int{1}),
			idx:      -1,
			expected: Slice[int]([]int{1}),
		},
		{
			name:     "slice with one item (idx greater than length)",
			payload:  Slice[int]([]int{1}),
			idx:      3,
			expected: Slice[int]([]int{1}),
		},
		{
			name:     "slice with two elements",
			payload:  Slice[int]([]int{1, 2}),
			idx:      0,
			expected: Slice[int]([]int{2}),
		},
		{
			name:     "delete keeps order",
			payload:  Slice[int]([]int{1, 2, 3, 4}),
			idx:      0,
			expected: Slice[int]([]int{2, 3, 4}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := DeleteOrder(test.payload, test.idx)

			if !test.expected.Equals(actual, testArrEq) {
				t.Errorf("unexpected value, want %v, have %v", test.expected, actual)
			}
		})
	}
}

func TestFind(t *testing.T) {
	type testCase struct {
		name       string
		payload    Slice[int]
		expected   int
		expectedOk bool
	}

	tests := []testCase{
		{
			name:       "nil slice should be noop",
			payload:    Slice[int]([]int{}),
			expected:   0,
			expectedOk: false,
		},
		{
			name:       "matched item is unique",
			payload:    Slice[int]([]int{2}),
			expected:   2,
			expectedOk: true,
		},
		{
			name:       "item in the first position",
			payload:    Slice[int]([]int{2, 1, 1, 1, 1}),
			expected:   2,
			expectedOk: true,
		},
		{
			name:       "item in the last position",
			payload:    Slice[int]([]int{1, 1, 1, 1, 2}),
			expected:   2,
			expectedOk: true,
		},
		{
			name:       "item in the middle",
			payload:    Slice[int]([]int{1, 1, 2, 1, 1, 4}),
			expected:   2,
			expectedOk: true,
		},
	}

	search := func(x int) bool {
		return x%2 == 0
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, ok := Find(test.payload, search)

			if test.expectedOk != ok || test.expected != actual {
				t.Errorf("unexpected value, want (%v, %t), have (%v, %t)",
					test.expected, test.expectedOk, actual, ok)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	type testCase struct {
		name        string
		payload     Slice[int]
		expected    int
		expectedArr Slice[int]
		expectedOk  bool
	}

	tests := []testCase{
		{
			name:        "nil slice should be noop",
			payload:     Slice[int]([]int{}),
			expected:    0,
			expectedOk:  false,
			expectedArr: Slice[int]([]int{}),
		},
		{
			name:        "matched item is unique",
			payload:     Slice[int]([]int{2}),
			expected:    2,
			expectedOk:  true,
			expectedArr: Slice[int]([]int{}),
		},
		{
			name:        "item in the first position",
			payload:     Slice[int]([]int{2, 1, 1, 1, 1}),
			expected:    2,
			expectedOk:  true,
			expectedArr: Slice[int]([]int{1, 1, 1, 1}),
		},
		{
			name:        "item in the last position",
			payload:     Slice[int]([]int{1, 1, 1, 1, 2}),
			expected:    2,
			expectedOk:  true,
			expectedArr: Slice[int]([]int{1, 1, 1, 1}),
		},
		{
			name:        "item in the middle",
			payload:     Slice[int]([]int{1, 1, 2, 1, 1, 4}),
			expected:    2,
			expectedOk:  true,
			expectedArr: Slice[int]([]int{1, 1, 4, 1, 1}),
		},
	}

	search := func(x int) bool {
		return x%2 == 0
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			arr, actual, ok := Extract(test.payload, search)

			if test.expectedOk != ok || test.expected != actual ||
				!test.expectedArr.Equals(arr, testArrEq) {
				t.Errorf("unexpected value, want (%v, %v, %t), have (%v, %v, %t)",
					test.expectedArr, test.expected, test.expectedOk,
					arr, actual, ok)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	type testCase struct {
		name     string
		payload  Slice[int]
		item     int
		idx      int
		expected Slice[int]
	}

	tests := []testCase{
		{
			name:     "nil slice should create a new one",
			payload:  nil,
			item:     1,
			idx:      0,
			expected: Slice[int]([]int{1}),
		},
		{
			name:     "empty slice should insert at first position",
			payload:  Slice[int]([]int{}),
			item:     1,
			idx:      0,
			expected: Slice[int]([]int{1}),
		},
		{
			name:     "insert at first position",
			payload:  Slice[int]([]int{2}),
			item:     1,
			idx:      0,
			expected: Slice[int]([]int{1, 2}),
		},
		{
			name:     "insert at last position",
			payload:  Slice[int]([]int{2}),
			item:     1,
			idx:      1,
			expected: Slice[int]([]int{2, 1}),
		},
		{
			name:     "insert middle position",
			payload:  Slice[int]([]int{1, 3}),
			item:     2,
			idx:      1,
			expected: Slice[int]([]int{1, 2, 3}),
		},
		{
			name:     "out of bounds from left is noop",
			payload:  Slice[int]([]int{1, 3}),
			item:     2,
			idx:      -1,
			expected: Slice[int]([]int{1, 3}),
		},
		{
			name:     "out of bounds from right is noop",
			payload:  Slice[int]([]int{1, 3}),
			item:     2,
			idx:      3,
			expected: Slice[int]([]int{1, 3}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Insert(test.payload, test.item, test.idx)

			if !test.expected.Equals(actual, testArrEq) {
				t.Errorf("unexpected value, want %v, have %v",
					test.expected, actual)
			}
		})
	}
}

func TestInsertVector(t *testing.T) {
	type testCase struct {
		name     string
		payload  Slice[int]
		items    []int
		idx      int
		expected Slice[int]
	}

	tests := []testCase{
		{
			name:     "nil slice should create a new one",
			payload:  nil,
			items:    []int{1},
			idx:      0,
			expected: Slice[int]([]int{1}),
		},
		{
			name:     "empty slice should insert at first position",
			payload:  Slice[int]([]int{}),
			items:    []int{1, 2},
			idx:      0,
			expected: Slice[int]([]int{1, 2}),
		},
		{
			name:     "insert at first position",
			payload:  Slice[int]([]int{2}),
			items:    []int{1, 2},
			idx:      0,
			expected: Slice[int]([]int{1, 2, 2}),
		},
		{
			name:     "insert at last position",
			payload:  Slice[int]([]int{2}),
			items:    []int{3, 5},
			idx:      1,
			expected: Slice[int]([]int{2, 3, 5}),
		},
		{
			name:     "insert middle position",
			payload:  Slice[int]([]int{1, 3}),
			items:    []int{2, 4},
			idx:      1,
			expected: Slice[int]([]int{1, 2, 4, 3}),
		},
		{
			name:     "insert empty is noop",
			payload:  Slice[int]([]int{1, 3}),
			items:    []int{},
			idx:      1,
			expected: Slice[int]([]int{1, 3}),
		},
		{
			name:     "out of bounds from left is noop",
			payload:  Slice[int]([]int{1, 3}),
			items:    []int{},
			idx:      -1,
			expected: Slice[int]([]int{1, 3}),
		},
		{
			name:     "out of bounds from right is noop",
			payload:  Slice[int]([]int{1, 3}),
			items:    []int{},
			idx:      3,
			expected: Slice[int]([]int{1, 3}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := InsertVector(test.payload, test.items, test.idx)

			if !test.expected.Equals(actual, testArrEq) {
				t.Errorf("unexpected value, want %v, have %v",
					test.expected, actual)
			}
		})
	}
}

func TestPop(t *testing.T) {
	var (
		payload = []int{1, 2}
		item    int
		ok      bool
	)

	payload, item, ok = Pop(payload)

	if item != 2 || !ok {
		t.Errorf("unexpected values, want (%d, %t), have (%d, %t)",
			2, true,
			item, ok,
		)
	}
	payload, item, ok = Pop(payload)

	if item != 1 || !ok {
		t.Errorf("unexpected values, want (%d, %t), have (%d, %t)",
			1, true,
			item, ok,
		)
	}

	payload, item, ok = Pop(payload)

	if item != 0 || ok {
		t.Errorf("unexpected values, want (%d, %t), have (%d, %t)",
			0, false,
			item, ok,
		)
	}
}

func testArrEq(x, y int) bool { return x == y }

// ExampleFilter demonstrates filtering a slice to keep only elements that satisfy a condition.
func ExampleFilter() {
	// Create a slice of numbers
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Filter to keep only even numbers
	evenNumbers := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})

	fmt.Println(evenNumbers)
	// Output: [2 4 6 8 10]
}

// ExampleMap demonstrates transforming elements in a slice.
func ExampleMap() {
	// Create a slice of numbers
	numbers := []int{1, 2, 3, 4, 5}

	// Transform each number by squaring it
	squares := Map(numbers, func(n int) int {
		return n * n
	})

	fmt.Println(squares)
	// Output: [1 4 9 16 25]
}

// ExampleFilterMap demonstrates filtering and transforming in a single operation.
func ExampleFilterMap() {
	// Create a slice of numbers
	numbers := []int{1, 2, 3, 4, 5, 6}

	// Keep only even numbers and square them
	evenSquares := FilterMap(numbers, func(n int) fp.Option[int] {
		if n%2 == 0 {
			return fp.Some(n * n)
		}
		return fp.None[int]()
	})

	fmt.Println(evenSquares)
	// Output: [4 16 36]
}

// ExampleReduce demonstrates combining all elements into a single value.
func ExampleReduce() {
	// Create a slice of numbers
	numbers := []int{1, 2, 3, 4, 5}

	// Sum all numbers
	sum := Reduce[int, int](numbers, func(acc, curr int) int {
		return acc + curr
	})

	fmt.Println(sum)
	// Output: 15
}

// ExampleFind demonstrates finding the first element that matches a condition.
func ExampleFind() {
	// Create a slice of names
	names := []string{"Alice", "Bob", "Charlie", "David"}

	// Find the first name that starts with 'C'
	result, found := Find(names, func(name string) bool {
		return len(name) > 0 && name[0] == 'C'
	})

	fmt.Printf("Found: %t, Name: %s\n", found, result)
	// Output: Found: true, Name: Charlie
}

// ExampleContains demonstrates checking if any element satisfies a condition.
func ExampleContains() {
	// Create a slice of numbers
	numbers := []int{1, 2, 3, 4, 5}

	// Check if any number is greater than 3
	hasLarge := Contains(numbers, func(n int) bool {
		return n > 3
	})

	// Check if any number is negative
	hasNegative := Contains(numbers, func(n int) bool {
		return n < 0
	})

	fmt.Printf("Has number > 3: %t\n", hasLarge)
	fmt.Printf("Has negative: %t\n", hasNegative)
	// Output:
	// Has number > 3: true
	// Has negative: false
}

// ExampleSome demonstrates checking if some elements satisfy a condition.
func ExampleSome() {
	// Create a slice of words
	words := []string{"hello", "world", "go", "programming"}

	// Check if some words are short (< 4 characters)
	hasShort := Some(words, func(word string) bool {
		return len(word) < 4
	})

	fmt.Printf("Has short words: %t\n", hasShort)
	// Output:
	// Has short words: true
}

// ExampleAll demonstrates checking if all elements satisfy a condition.
func ExampleAll() {
	// Create a slice of positive numbers
	numbers := []int{1, 2, 3, 4, 5}

	// Check if all numbers are positive
	allPositive := All(numbers, func(n int) bool {
		return n > 0
	})

	// Check if all numbers are even
	allEven := All(numbers, func(n int) bool {
		return n%2 == 0
	})

	fmt.Printf("All positive: %t\n", allPositive)
	fmt.Printf("All even: %t\n", allEven)
	// Output:
	// All positive: true
	// All even: false
}

// ExampleFold demonstrates folding a slice with an initial value.
func ExampleFold() {
	// Create a slice of strings
	words := []string{"Hello", "World", "from", "Go"}

	// Join with custom separator and prefix
	result := Fold(words, func(acc, word string) string {
		if acc == "" {
			return "Greeting: " + word
		}
		return acc + " " + word
	}, "")

	fmt.Printf("Result: %s\n", result)
	// Output:
	// Result: Greeting: Hello World from Go
}

// ExampleToMap demonstrates converting a slice to a map using a key function.
func ExampleToMap() {
	// Create a slice of words
	words := []string{"apple", "banana", "cherry"}

	// Convert to map with first letter as key
	wordMap := ToMap(words, func(word string) rune {
		return rune(word[0])
	})

	fmt.Printf("'a' word: %s\n", wordMap['a'])
	fmt.Printf("'b' word: %s\n", wordMap['b'])
	// Output:
	// 'a' word: apple
	// 'b' word: banana
}

func ExampleEquals() {
	numbers1 := []int{1, 2, 3}
	numbers2 := []int{1, 2, 3}
	numbers3 := []int{1, 2, 4}

	// Compare two slices for equality
	equal := Equals(numbers1, numbers2, func(a, b int) bool { return a == b })
	fmt.Println("numbers1 equals numbers2:", equal)

	equal = Equals(numbers1, numbers3, func(a, b int) bool { return a == b })
	fmt.Println("numbers1 equals numbers3:", equal)

	// Output:
	// numbers1 equals numbers2: true
	// numbers1 equals numbers3: false
}

func ExampleToMapIdx() {
	fruits := []string{"apple", "banana", "cherry"}

	// Map fruits to their indices
	result := ToMapIdx(fruits, func(fruit string) string { return fruit })

	fmt.Println("apple:", result["apple"])
	fmt.Println("banana:", result["banana"])

	// Output:
	// apple: {apple 0}
	// banana: {banana 1}
}

func ExampleIndexOf() {
	numbers := []int{10, 20, 30, 40, 50}

	// Find index of first element greater than 25
	index := IndexOf(numbers, func(n int) bool { return n > 25 })
	fmt.Println("Index of first number > 25:", index)

	// Find index of non-existent element
	index = IndexOf(numbers, func(n int) bool { return n > 100 })
	fmt.Println("Index of first number > 100:", index)

	// Output:
	// Index of first number > 25: 2
	// Index of first number > 100: -1
}

func ExampleIncludes() {
	numbers := []int{1, 2, 3, 4, 5}

	fmt.Println("Contains 3:", Includes(numbers, 3))
	fmt.Println("Contains 6:", Includes(numbers, 6))

	// Output:
	// Contains 3: true
	// Contains 6: false
}

func ExampleAny() {
	numbers := []int{1, 2, 3, 4, 5}

	// Check if any number is even
	hasEven := Any(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Println("Has even numbers:", hasEven)

	// Check if any number is greater than 10
	hasLarge := Any(numbers, func(n int) bool { return n > 10 })
	fmt.Println("Has numbers > 10:", hasLarge)

	// Output:
	// Has even numbers: true
	// Has numbers > 10: false
}

func ExampleMapInPlace() {
	numbers := []int{1, 2, 3, 4, 5}

	// Double each number in place
	result := MapInPlace(numbers, func(n int) int { return n * 2 })
	fmt.Println("Doubled:", result)

	// Output:
	// Doubled: [2 4 6 8 10]
}

func ExampleFilterMapTuple() {
	numbers := []int{1, 2, 3, 4, 5}

	// Square only even numbers
	squares := FilterMapTuple(numbers, func(n int) (int, bool) {
		if n%2 == 0 {
			return n * n, true
		}
		return 0, false
	})
	fmt.Println("Squares of even numbers:", squares)

	// Output:
	// Squares of even numbers: [4 16]
}

func ExampleFilterInPlace() {
	numbers := []int{1, 2, 3, 4, 5, 6}

	// Keep only even numbers
	result := FilterInPlace(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Println("Even numbers:", result)

	// Output:
	// Even numbers: [2 4 6]
}

func ExampleReduceSame() {
	numbers := []int{1, 2, 3, 4, 5}

	// Sum all numbers
	sum := ReduceSame(numbers, func(a, b int) int { return a + b })
	fmt.Println("Sum:", sum)

	// Output:
	// Sum: 15
}

func ExampleFoldSame() {
	numbers := []int{1, 2, 3, 4, 5}

	// Sum with initial value
	sum := FoldSame(numbers, func(a, b int) int { return a + b }, 10)
	fmt.Println("Sum with initial 10:", sum)

	// Output:
	// Sum with initial 10: 25
}

func ExampleCut() {
	numbers := []int{1, 2, 3, 4, 5}

	// Remove elements from index 1 to 3
	result := Cut(numbers, 1, 3)
	fmt.Println("After cutting [1:3]:", result)

	// Output:
	// After cutting [1:3]: [1 5]
}

func ExampleAppend() {
	numbers := []int{1, 2, 3}

	// Add element to the end
	result := Append(numbers, 4)
	fmt.Println("After appending 4:", result)

	// Output:
	// After appending 4: [1 2 3 4]
}

func ExampleAppendVector() {
	numbers := []int{1, 2, 3}
	moreNumbers := []int{4, 5, 6}

	// Concatenate slices
	result := AppendVector(numbers, moreNumbers)
	fmt.Println("After appending vector:", result)

	// Output:
	// After appending vector: [1 2 3 4 5 6]
}

func ExampleDelete() {
	numbers := []int{1, 2, 3, 4, 5}

	// Delete element at index 2 (value 3)
	result := Delete(numbers, 2)
	fmt.Println("After deleting index 2:", result)

	// Output:
	// After deleting index 2: [1 2 5 4]
}

func ExampleDeleteOrder() {
	numbers := []int{1, 2, 3, 4, 5}

	// Delete element at index 2 preserving order
	result := DeleteOrder(numbers, 2)
	fmt.Println("After deleting index 2 (order preserved):", result)

	// Output:
	// After deleting index 2 (order preserved): [1 2 4 5]
}

func ExampleFindIdx() {
	numbers := []int{10, 20, 30, 40, 50}

	// Find first element greater than 25
	value, index := FindIdx(numbers, func(n int) bool { return n > 25 })
	fmt.Printf("Found %d at index %d\n", value, index)

	// Find non-existent element
	value, index = FindIdx(numbers, func(n int) bool { return n > 100 })
	fmt.Printf("Found %d at index %d\n", value, index)

	// Output:
	// Found 30 at index 2
	// Found 0 at index -1
}

func ExampleExtractIdx() {
	numbers := []int{1, 2, 3, 4, 5}

	// Extract element at index 2
	remaining, extracted, ok := ExtractIdx(numbers, 2)
	fmt.Printf("Extracted: %d, OK: %t\n", extracted, ok)
	fmt.Println("Remaining:", remaining)

	// Output:
	// Extracted: 3, OK: true
	// Remaining: [1 2 5 4]
}

func ExampleExtract() {
	numbers := []int{1, 2, 3, 4, 5}

	// Extract first even number
	remaining, extracted, ok := Extract(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Printf("Extracted: %d, OK: %t\n", extracted, ok)
	fmt.Println("Remaining:", remaining)

	// Output:
	// Extracted: 2, OK: true
	// Remaining: [1 5 3 4]
}

func ExamplePop() {
	numbers := []int{1, 2, 3, 4, 5}

	// Remove last element
	remaining, popped, ok := Pop(numbers)
	fmt.Printf("Popped: %d, OK: %t\n", popped, ok)
	fmt.Println("Remaining:", remaining)

	// Output:
	// Popped: 5, OK: true
	// Remaining: [1 2 3 4]
}

func ExamplePeek() {
	numbers := []int{1, 2, 3, 4, 5}

	// Peek at element at index 2
	value, ok := Peek(numbers, 2)
	fmt.Printf("Peeked: %d, OK: %t\n", value, ok)
	fmt.Println("Original:", numbers)

	// Output:
	// Peeked: 3, OK: true
	// Original: [1 2 3 4 5]
}

func ExamplePushFront() {
	numbers := []int{2, 3, 4}

	// Add element to the beginning
	result := PushFront(numbers, 1)
	fmt.Println("After pushing 1 to front:", result)

	// Output:
	// After pushing 1 to front: [1 2 3 4]
}

func ExampleUnshift() {
	numbers := []int{2, 3, 4}

	// Add element to the beginning (alias for PushFront)
	result := Unshift(numbers, 1)
	fmt.Println("After unshifting 1:", result)

	// Output:
	// After unshifting 1: [1 2 3 4]
}

func ExamplePopFront() {
	numbers := []int{1, 2, 3, 4, 5}

	// Remove first element
	remaining, popped, ok := PopFront(numbers)
	fmt.Printf("Popped: %d, OK: %t\n", popped, ok)
	fmt.Println("Remaining:", remaining)

	// Output:
	// Popped: 1, OK: true
	// Remaining: [2 3 4 5]
}

func ExampleShift() {
	numbers := []int{1, 2, 3, 4, 5}

	// Remove first element (alias for PopFront)
	remaining, shifted, ok := Shift(numbers)
	fmt.Printf("Shifted: %d, OK: %t\n", shifted, ok)
	fmt.Println("Remaining:", remaining)

	// Output:
	// Shifted: 1, OK: true
	// Remaining: [2 3 4 5]
}

func ExampleInsert() {
	numbers := []int{1, 2, 4, 5}

	// Insert 3 at index 2
	result := Insert(numbers, 3, 2)
	fmt.Println("After inserting 3 at index 2:", result)

	// Output:
	// After inserting 3 at index 2: [1 2 3 4 5]
}

func ExampleInsertVector() {
	numbers := []int{1, 2, 5, 6}
	toInsert := []int{3, 4}

	// Insert slice at index 2
	result := InsertVector(numbers, toInsert, 2)
	fmt.Println("After inserting [3, 4] at index 2:", result)

	// Output:
	// After inserting [3, 4] at index 2: [1 2 3 4 5 6]
}
