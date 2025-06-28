// Package slices provides a comprehensive set of generic utility functions for working with slices.
// It offers a functional approach to common slice operations such as transforming, filtering,
// searching, and manipulating elements in a type-safe manner.
package slices

import (
	"bytes"
	"fmt"

	"github.com/sonirico/gozo/fp"
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
	for i, x := range s {
		if !fn(x, i) {
			return
		}
	}
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

// Equals compares two slices and returns whether they contain equal elements.
// Two slices are considered equal if they have the same length and corresponding
// elements satisfy the equality function.
func Equals[T any](one, other []T, predicate func(x, y T) bool) (res bool) {
	if len(one) != len(other) {
		return
	}

	res = true

	for idx, otherItem := range other {
		res = predicate(one[idx], otherItem)
		if !res {
			return
		}
	}
	return
}

// IndexOf returns the index of the first element that satisfies the predicate.
// Returns the index where the element was found, or -1 if not found.
func (s Slice[T]) IndexOf(fn func(t T) bool) int {
	return IndexOf(s, fn)
}

// ToMap creates a map from a slice, using the provided function to determine the key
// for each element. The element itself becomes the value in the map.
func ToMap[V any, K comparable](arr []V, predicate func(x V) K) map[K]V {
	res := make(map[K]V, len(arr))

	for _, x := range arr {
		res[predicate(x)] = x
	}

	return res
}

type (
	// WrappedIdx stores an element along with its original index in the slice.
	WrappedIdx[T any] struct {
		value T   // The element value
		idx   int // The original index
	}
)

// ToMapIdx creates a map from a slice, preserving each element's original index.
// Uses the provided function to determine the key for each element.
// The value in the map is a WrappedIdx containing both the element and its original index.
func ToMapIdx[V any, K comparable](arr []V, predicate func(x V) K) map[K]WrappedIdx[V] {
	res := make(map[K]WrappedIdx[V], len(arr))

	for i, x := range arr {
		res[predicate(x)] = WrappedIdx[V]{value: x, idx: i}
	}

	return res
}

// IndexOf returns the index of the first element that satisfies the predicate.
// Returns the index where the element was found, or -1 if not found.
func IndexOf[T any](arr []T, predicate func(t T) bool) (pos int) {
	pos = -1
	for i, x := range arr {
		if predicate(x) {
			pos = i
			return
		}
	}
	return
}

// Contains checks if the slice contains an element that satisfies the predicate.
// Returns true if any element matches the predicate, false otherwise.
func Contains[T any](arr []T, predicate func(t T) bool) bool {
	return IndexOf(arr, predicate) >= 0
}

// Includes checks if the slice contains a specific element using the equality operator.
// Returns true if the element is found, false otherwise.
func Includes[T comparable](arr []T, target T) bool {
	return Contains(arr, func(t T) bool {
		return t == target
	})
}

// Some checks if at least one element in the slice satisfies the predicate.
// Returns true if any element matches the predicate, false otherwise.
// Alias for Contains.
func Some[T any](arr []T, predicate func(t T) bool) bool {
	return Contains(arr, predicate)
}

// Any checks if at least one element in the slice satisfies the predicate.
// Returns true if any element matches the predicate, false otherwise.
// Alias for Contains.
func Any[T any](arr []T, predicate func(t T) bool) bool {
	return Contains(arr, predicate)
}

// All checks if all elements in the slice satisfy the predicate.
// Returns true if all elements match the predicate, false otherwise.
func All[T any](arr []T, predicate func(t T) bool) bool {
	for _, x := range arr {
		if !predicate(x) {
			return false
		}
	}
	return true
}

// Map creates a new slice by applying the transformation function to each element.
// The transformation can change the type of the elements.
func Map[T, U any](arr []T, predicate func(t T) U) []U {
	res := make([]U, 0, len(arr))

	for _, x := range arr {
		res = append(res, predicate(x))
	}

	return res
}

// MapInPlace transforms each element in the slice using the provided function.
// Modifies the slice in place and returns it.
func MapInPlace[T any](arr []T, predicate func(t T) T) []T {
	for i, x := range arr {
		arr[i] = predicate(x)
	}

	return arr
}

// Filter creates a new slice containing only the elements that satisfy the predicate.
func Filter[T any](arr []T, predicate func(t T) bool) []T {
	res := make([]T, 0, len(arr))

	for _, x := range arr {
		if predicate(x) {
			res = append(res, x)
		}
	}

	return res
}

// FilterMapTuple creates a new slice by applying a transformation function that
// also filters elements. The function should return the transformed value and
// a boolean indicating whether to include the element.
func FilterMapTuple[T, U any](arr []T, predicate func(t T) (U, bool)) []U {
	res := make([]U, 0, len(arr))

	for _, x := range arr {
		if mapped, ok := predicate(x); ok {
			res = append(res, mapped)
		}
	}

	return res
}

// FilterMap creates a new slice by applying a transformation function that
// returns an Option. Elements with Some options are included in the result,
// while None options are excluded.
func FilterMap[T, U any](arr []T, predicate func(t T) fp.Option[U]) []U {
	pre := func(t T) (U, bool) {
		return predicate(t).Unwrap()
	}

	return FilterMapTuple[T, U](arr, pre)
}

// FilterInPlace modifies the slice in place to contain only elements that
// satisfy the predicate. This is more efficient than Filter when creating
// a new slice is not necessary.
func FilterInPlace[T any](arr []T, predicate func(t T) bool) []T {
	n := 0
	for i, x := range arr {
		if predicate(x) {
			if n != i {
				arr[n] = x
			}
			n++
		}
	}

	arr = arr[:n]

	return arr
}

// Reduce compacts the slice into a single value by iteratively applying
// the reduction function to each element. Starts with the zero value.
func Reduce[T, U any](arr []T, p func(T, T) T) (res T) {
	return Fold(arr, p, res)
}

// ReduceSame is a convenience wrapper around Reduce for when the accumulator
// and element types are the same.
func ReduceSame[T any](arr []T, p func(T, T) T) T {
	return Reduce[T, T](arr, p)
}

// FoldSame is a convenience wrapper around Fold for when the accumulator
// and element types are the same.
func FoldSame[T any](arr []T, p func(T, T) T, initial T) T {
	return Fold[T, T](arr, p, initial)
}

// Fold compacts the slice into a single value by iteratively applying
// the reduction function, starting with the provided initial value.
// The accumulator type can be different from the element type.
func Fold[T, U any](arr []T, p func(U, T) U, initial U) U {
	if len(arr) < 1 {
		return initial
	}

	initial = p(initial, arr[0])

	if len(arr) < 2 {
		return initial
	}

	i := 1

	for i < len(arr) {
		initial = p(initial, arr[i])

		i++
	}

	return initial
}

// Cut removes a sector from slice given lower and upper bounds. Bounds are
// represented as indices of the slice. E.g:
// Cut([1, 2, 3, 4], 1, 2) -> [1, 4]
// Cut([4], 0, 0) -> []
// Cut will returned the original slice without the cut subslice.
func Cut[T any](arr []T, from, to int) []T {
	if len(arr) < 1 {
		return arr
	}

	if from < 0 {
		from = 0
	}

	if from >= len(arr) {
		from = len(arr) - 1
	}

	if to < 0 {
		to = 0
	}

	if to >= len(arr) {
		to = len(arr) - 1
	}

	if len(arr) == 1 {
		return arr[:0]
	}

	if from > to {
		// In this case, consider `to` to be the amount to remove from `from`.
		return append(arr[:from], arr[from+to+1:]...)
	}

	return append(arr[:from], arr[to+1:]...)
}

// Append adds an element to the end of the slice and returns the result.
// Unlike the builtin append, this function is explicitly named for clarity.
func Append[T any](arr []T, item T) []T {
	return append(arr, item)
}

// AppendVector adds all elements from another slice to the end of this slice.
// Returns the resulting concatenated slice.
func AppendVector[T any](arr, items []T) []T {
	return append(arr, items...)
}

// Delete removes the element at the specified index without preserving order.
// This provides better performance than DeleteOrder but changes the order of elements.
// If the index is out of bounds, returns the original slice unchanged.
func Delete[T any](arr []T, idx int) []T {
	le := len(arr) - 1
	if le < 0 || idx > le || idx < 0 {
		return arr
	}
	var t T
	arr[idx] = arr[le]
	arr[le] = t
	arr = arr[:le]
	return arr
}

// DeleteOrder removes the element at the specified index while preserving order.
// This is slower than Delete but maintains the relative order of the remaining elements.
// If the index is out of bounds, returns the original slice unchanged.
func DeleteOrder[T any](arr []T, idx int) []T {
	le := len(arr) - 1
	if le < 0 || idx > le || idx < 0 {
		return arr
	}
	var t T

	if le > 0 {
		copy(arr[idx:], arr[idx+1:])
	}

	arr[le] = t
	arr = arr[:le]
	return arr
}

// Find returns the first element that satisfies the predicate.
// Returns the element and true if found, otherwise the zero value and false.
func Find[T any](arr []T, predicate func(t T) bool) (res T, ok bool) {
	var idx int
	res, idx = FindIdx(arr, predicate)
	ok = idx > -1
	return
}

// FindIdx returns the first element that satisfies the predicate and its index.
// Returns the element and its index if found, otherwise the zero value and -1.
func FindIdx[T any](arr []T, predicate func(t T) bool) (res T, idx int) {
	idx = IndexOf(arr, predicate)
	if idx < 0 {
		return
	}

	res = arr[idx]
	return
}

// ExtractIdx gets and deletes the element at the given position.
// Returns the modified slice, the extracted element, and a success flag.
// If the index is out of bounds, returns the original slice, zero value, and false.
func ExtractIdx[T any](arr []T, idx int) (res []T, item T, ok bool) {
	if idx >= len(arr) || idx < 0 {
		return
	}

	ok = true
	item = arr[idx]
	res = Delete(arr, idx)

	return
}

// Extract gets and deletes the first element that matches the predicate.
// Returns the modified slice, the extracted element, and a success flag.
// If no element matches, returns the original slice, zero value, and false.
func Extract[T any](arr []T, predicate func(t T) bool) ([]T, T, bool) {
	res, idx := FindIdx(arr, predicate)
	if idx < 0 {
		return arr, res, false
	}

	arr = Delete(arr, idx)
	return arr, res, true
}

// Pop deletes and returns the last item from the slice.
// Returns the modified slice, the popped element, and a success flag.
// If the slice is empty, returns the original slice, zero value, and false.
func Pop[T any](arr []T) (res []T, item T, ok bool) {
	if len(arr) < 1 {
		return
	}

	var t T
	le := len(arr) - 1
	res = arr[:le]
	item = arr[le]
	ok = true

	arr[le] = t // GC

	return
}

// Peek returns the item at the specified index without modifying the slice.
// Returns the element and true if the index is valid, otherwise the zero value and false.
func Peek[T any](arr []T, idx int) (item T, ok bool) {
	if len(arr) < 1 || idx >= len(arr) {
		return
	}

	item = arr[idx]
	ok = true

	return
}

// PushFront inserts an element at the beginning of the slice.
// Returns the resulting slice with the new element at the front.
func PushFront[T any](arr []T, item T) []T {
	return append([]T{item}, arr...)
}

// Unshift inserts an element at the beginning of the slice.
// Alias for PushFront, following JavaScript array method naming conventions.
func Unshift[T any](arr []T, item T) []T {
	return PushFront(arr, item)
}

// PopFront removes and returns the first element of the slice.
// Returns the modified slice (without the first element), the removed element, and a success flag.
// If the slice is empty, returns the original slice, zero value, and false.
func PopFront[T any](arr []T) (res []T, item T, ok bool) {
	if len(arr) < 1 {
		res = arr
		return
	}

	item, res = arr[0], arr[1:]
	ok = true
	return
}

// Shift removes and returns the first element of the slice.
// Alias for PopFront, following JavaScript array method naming conventions.
func Shift[T any](arr []T) ([]T, T, bool) {
	return PopFront(arr)
}

// Insert places an element at the specified index in the slice.
// Elements at or after the index are shifted to the right.
// Returns the resulting slice with the new element inserted.
// If the index is out of bounds, returns the original slice unchanged.
func Insert[T any](arr []T, item T, idx int) []T {
	if arr == nil {
		return []T{item}
	}

	if idx < 0 || idx > len(arr) {
		return arr
	}

	return append(arr[:idx], append([]T{item}, arr[idx:]...)...)
}

// InsertVector places a slice of elements at the specified index in the slice.
// Elements at or after the index are shifted to the right.
// Returns the resulting slice with the new elements inserted.
// If the index is out of bounds, returns the original slice unchanged.
func InsertVector[T any](arr, items []T, idx int) (res []T) {
	if arr == nil {
		res = items[:]
		return
	}

	if items == nil || len(items) == 0 {
		res = arr
		return
	}

	if idx < 0 || idx > len(arr) {
		return arr
	}

	return append(arr[:idx], append(items, arr[idx:]...)...)
}
