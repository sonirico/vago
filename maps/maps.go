// Package maps provides generic utility functions to work with Go maps.
// It offers a functional approach to common map operations like filtering, mapping,
// reducing, and comparing maps.
package maps

import (
	"github.com/sonirico/gozo/fp"
	"github.com/sonirico/gozo/slices"
	"github.com/sonirico/gozo/tuples"
)

// Equals compares two maps and returns whether they are equal in values.
// Two maps are considered equal if:
// - They have the same length
// - They contain the same keys
// - For each key, the values in both maps satisfy the equality function
//
// Maps are compared using the provided equality function for values.
// This allows for deep equality checks on complex value types.
func Equals[K comparable, V any](m1, m2 map[K]V, eq func(V, V) bool) bool {
	if len(m1) != len(m2) {
		return false
	}

	if m1 == nil && m2 != nil {
		return false
	}

	if m1 != nil && m2 == nil {
		return false
	}

	for k1, v1 := range m1 {
		v2, ok := m2[k1]
		if !ok {
			return false
		}

		if !eq(v1, v2) {
			return false
		}
	}

	return true
}

// Map transforms a map into another map, with potentially different key and value types.
// The transformation is applied to each key-value pair by the provided function,
// which returns the new key and value for the resulting map.
//
// This function preserves nil semantics: if the input map is nil, the output will also be nil.
// Otherwise, a new map is created with the transformed key-value pairs.
func Map[K1 comparable, V1 any, K2 comparable, V2 any](
	m map[K1]V1,
	p func(K1, V1) (K2, V2),
) map[K2]V2 {
	if m == nil {
		return nil
	}

	res := make(map[K2]V2, len(m))

	for k1, v1 := range m {
		k2, v2 := p(k1, v1)
		res[k2] = v2
	}

	return res
}

// FilterMap both filters and maps a map into a new map, potentially with different key and value types.
// The predicate function should return an fp.Option monad containing a tuple of the new key and value:
// - fp.Some to include the entry in the result (with transformed key and value)
// - fp.None to exclude the entry from the result
//
// This provides a powerful way to simultaneously transform and filter map entries
// while leveraging the Option monad for expressing presence/absence.
func FilterMap[K1 comparable, V1 any, K2 comparable, V2 any](
	m map[K1]V1,
	p func(K1, V1) fp.Option[tuples.Tuple2[K2, V2]],
) map[K2]V2 {
	if m == nil {
		return nil
	}

	res := make(map[K2]V2, len(m))

	for k1, v1 := range m {
		tpl := p(k1, v1)
		if tpl.IsSome() {
			v := tpl.UnwrapUnsafe()
			res[v.V1] = v.V2
		}
	}

	return res
}

// FilterMapTuple both filters and maps the given map into a new map, potentially with different key and value types.
// The predicate function returns three values:
// - The new key (K2)
// - The new value (V2)
// - A boolean indicating whether to include this entry in the result
//
// This function is an alternative to FilterMap that uses Go's native boolean return
// instead of the Option monad for expressing presence/absence.
func FilterMapTuple[K1 comparable, V1 any, K2 comparable, V2 any](
	m map[K1]V1,
	p func(K1, V1) (K2, V2, bool),
) map[K2]V2 {
	if m == nil {
		return nil
	}

	res := make(map[K2]V2, len(m))

	for k1, v1 := range m {
		if k2, v2, ok := p(k1, v1); ok {
			res[k2] = v2
		}
	}

	return res
}

// Filter creates a new map containing only the key-value pairs that satisfy the predicate.
// The predicate function takes a key and value and returns a boolean indicating
// whether to include the entry in the result.
//
// Unlike FilterInPlace, this function creates a new map and does not modify the input map.
func Filter[K comparable, V any](
	m map[K]V,
	p func(K, V) bool,
) map[K]V {
	if m == nil {
		return nil
	}

	res := make(map[K]V, len(m))

	for k, v := range m {
		if p(k, v) {
			res[k] = v
		}
	}

	return res
}

// FilterInPlace modifies the given map by removing entries that do not satisfy the predicate.
// The predicate function takes a key and value and returns a boolean indicating
// whether to keep the entry in the map.
//
// This function directly modifies the input map for better performance when
// creating a new map is not necessary.
// It returns the modified map for convenience in chaining operations.
func FilterInPlace[K comparable, V any](
	m map[K]V,
	p func(K, V) bool,
) map[K]V {
	if m == nil {
		return nil
	}

	for k, v := range m {
		if !p(k, v) {
			delete(m, k)
		}
	}

	return m
}

// Reduce compacts a map into a single value by iteratively applying a reduction function.
// The reduction function takes the accumulator, a key, and a value, and returns
// the updated accumulator.
//
// The initial value for the accumulator is the zero value of type R.
// If you need a different initial value, use Fold instead.
func Reduce[K comparable, V any, R any](
	m map[K]V,
	p func(R, K, V) R,
) R {
	var r R

	if m == nil {
		return r
	}

	for k, v := range m {
		r = p(r, k, v)
	}

	return r
}

// Fold compacts a map into a single value by iteratively applying a reduction function
// with an explicit initial value.
// The reduction function takes the accumulator, a key, and a value, and returns
// the updated accumulator.
//
// Unlike Reduce, Fold takes an explicit initial value for the accumulator.
// This is useful when the zero value of the result type is not appropriate
// as the starting value.
func Fold[K comparable, V any, R any](
	m map[K]V,
	p func(R, K, V) R,
	initial R,
) R {
	if m == nil {
		return initial
	}

	r := initial

	for k, v := range m {
		r = p(r, k, v)
	}

	return r
}

// Slice converts a map into a slice by applying a transformation function to each key-value pair.
// The transformation function takes a key and value and returns an element
// for the resulting slice.
//
// The order of elements in the resulting slice is not guaranteed, as map iteration
// in Go is not deterministic.
func Slice[K comparable, V, R any](
	m map[K]V,
	p func(K, V) R,
) slices.Slice[R] {
	res := make([]R, len(m))
	i := 0

	for k, v := range m {
		res[i] = p(k, v)
		i++
	}

	return res
}
