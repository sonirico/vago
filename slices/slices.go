// Package slices provides a comprehensive set of generic utility functions for working with slices.
// It offers a functional approach to common slice operations such as transforming, filtering,
// searching, and manipulating elements in a type-safe manner.
package slices

import (
	"bytes"
	"fmt"

	"github.com/sonirico/vago/fp"
)

type (
	// Slice is a generic slice type that provides a rich set of operations.
	// It wraps a standard Go slice and extends it with methods for common operations.
	Slice[T any] []T
)

// String returns a string representation of the slice, with each element on a new line.
// Useful for debugging and displaying slice contents.
func (s Slice[T]) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("[\n")
	s.Range(func(x T, i int) bool {
		buf.WriteString(fmt.Sprintf("\t%d -> %v\n", i, x))
		return true
	})
	buf.WriteString("]\n")
	return buf.String()
}

// Len returns the number of elements in the slice.
func (s Slice[T]) Len() int {
	return len(s)
}

// Range iterates over each element in the slice, calling the provided function with
// each element and its index. Iteration stops if the function returns false.
func (s Slice[T]) Range(fn func(t T, i int) bool) {
	Range(s, fn)
}

// ForEach applies the provided function to each element in the slice.
// It iterates through all elements, allowing side effects or processing without
// returning a new slice.
func (s Slice[T]) ForEach(fn func(t T)) {
	ForEach(s, fn)
}

// Get safely retrieves the element at the specified index.
// Returns the element and true if the index is valid, otherwise returns
// the zero value and false.
func (s Slice[T]) Get(i int) (res T, ok bool) {
	ok = i >= 0 && i < len(s)
	if !ok {
		return
	}
	res = s[i]
	return
}

// Contains checks if the slice contains an element that satisfies the predicate.
// Returns true if any element matches the predicate, false otherwise.
func (s Slice[T]) Contains(fn func(t T) bool) bool {
	return Contains(s, fn)
}

// Equals compares this slice with another slice using the provided equality function.
// Returns true if both slices have the same length and corresponding elements
// satisfy the equality function.
func (s Slice[T]) Equals(other Slice[T], predicate func(x, y T) bool) (res bool) {
	return Equals(s, other, predicate)
}

// Clone creates a new slice with the same elements as this slice.
func (s Slice[T]) Clone() Slice[T] {
	res := make([]T, len(s))
	copy(res, s)
	return res
}

// Delete removes the element at the specified index without preserving order.
// Modifies the slice in place and returns it.
func (s *Slice[T]) Delete(idx int) Slice[T] {
	*s = Delete(*s, idx)
	return *s
}

// Push adds an element to the end of the slice.
// Modifies the slice in place and returns it.
func (s *Slice[T]) Push(item T) Slice[T] {
	return s.Append(item)
}

// Append adds an element to the end of the slice.
// Modifies the slice in place and returns it.
func (s *Slice[T]) Append(item T) Slice[T] {
	*s = append(*s, item)
	return *s
}

// AppendVector adds all elements from another slice to the end of this slice.
// Modifies the slice in place and returns it.
func (s *Slice[T]) AppendVector(items []T) Slice[T] {
	*s = append(*s, items...)
	return *s
}

// Map creates a new slice by applying the transformation function to each element.
func (s Slice[T]) Map(predicate func(T) T) Slice[T] {
	return Map(s, predicate)
}

// MapInPlace transforms each element in the slice using the provided function.
// Modifies the slice in place and returns it.
func (s Slice[T]) MapInPlace(predicate func(T) T) Slice[T] {
	return MapInPlace(s, predicate)
}

// Filter creates a new slice containing only the elements that satisfy the predicate.
func (s Slice[T]) Filter(predicate func(x T) bool) Slice[T] {
	return Filter(s, predicate)
}

// FilterMapTuple creates a new slice by applying a transformation function that
// also filters elements. The function should return the transformed value and
// a boolean indicating whether to include the element.
func (s Slice[T]) FilterMapTuple(predicate func(x T) (T, bool)) Slice[T] {
	return FilterMapTuple(s, predicate)
}

// FilterMap creates a new slice by applying a transformation function that
// returns an Option. Elements with Some options are included in the result,
// while None options are excluded.
func (s Slice[T]) FilterMap(predicate func(x T) fp.Option[T]) Slice[T] {
	return FilterMap(s, predicate)
}

// FilterInPlace modifies the slice in place to contain only elements that
// satisfy the predicate.
func (s Slice[T]) FilterInPlace(predicate func(x T) bool) Slice[T] {
	return FilterInPlace(s, predicate)
}

// Reduce compacts the slice into a single value by iteratively applying
// the reduction function to each element.
func (s Slice[T]) Reduce(predicate func(x, y T) T) T {
	return ReduceSame(s, predicate)
}

// Fold compacts the slice into a single value by iteratively applying
// the reduction function, starting with the provided initial value.
func (s Slice[T]) Fold(predicate func(x, y T) T, initial T) T {
	return FoldSame(s, predicate, initial)
}

func (s Slice[T]) Copy() Slice[T] {
	// Copy creates a shallow copy of the slice.
	// Returns a new slice with the same elements.
	if s == nil {
		return nil
	}

	res := make([]T, len(s))
	copy(res, s)
	return res
}
