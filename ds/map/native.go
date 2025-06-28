package _map

import (
	"github.com/sonirico/vago/fp"
	"github.com/sonirico/vago/maps"
	"github.com/sonirico/vago/slices"
	"github.com/sonirico/vago/tuples"
)

type (
	// Entry represents a key-value pair in a map.
	// It's used for returning map entries as structured data.
	Entry[K comparable, V any] struct {
		K K // The key
		V V // The value
	}

	// Native implements the Map interface using Go's native map type.
	// It provides generic map operations for keys of type K and values of type V.
	Native[K comparable, V any] struct {
		data map[K]V
	}
)

// NewNative creates a new Native map with keys of type K and values of type V.
// It initializes the underlying Go map and returns the Native wrapper.
func NewNative[K comparable, V any]() Native[K, V] {
	return Native[K, V]{data: make(map[K]V)}
}

// Get retrieves a value by its key.
// Returns the value and a boolean indicating whether the key exists.
func (m Native[K, V]) Get(k K) (v V, ok bool) {
	v, ok = m.data[k]
	return
}

// Has checks if the map contains the specified key.
func (m Native[K, V]) Has(k K) (ok bool) {
	_, ok = m.data[k]
	return
}

// Set adds or replaces a key-value pair in the map.
func (m Native[K, V]) Set(k K, v V) {
	m.data[k] = v
	return
}

// Range iterates over each key-value pair in the map, with its index.
// The iteration stops if the provided function returns false.
func (m Native[K, V]) Range(fn func(K, V, int) bool) {
	i := 0
	for k, v := range m.data {
		if !fn(k, v, i) {
			return
		}
		i++
	}
}

// Delete removes a key-value pair from the map.
func (m Native[K, V]) Delete(k K) {
	delete(m.data, k)
}

// GetOrSet retrieves a value or sets a default if the key doesn't exist.
// Returns the value (existing or default) and whether the key already existed.
func (m Native[K, V]) GetOrSet(k K, def V) (v V, ok bool) {
	if v, ok = m.data[k]; ok {
		return
	}

	m.data[k] = def
	v = def
	ok = true
	return
}

// Map applies a transformation function to each key-value pair and returns a new map.
// The function receives each key and value and should return the transformed key and value.
func (m Native[K, V]) Map(fn func(K, V) (K, V)) Map[K, V] {
	return Native[K, V]{data: maps.Map(m.data, fn)}
}

// FilterMap applies a function that may filter out or transform key-value pairs.
// Returns a new map containing only the key-value pairs for which the function
// returns a Some result. Useful for simultaneously filtering and transforming.
func (m Native[K, V]) FilterMap(fn func(K, V) fp.Option[tuples.Tuple2[K, V]]) Map[K, V] {
	return Native[K, V]{data: maps.FilterMap(m.data, fn)}
}

// Filter returns a new map containing only the key-value pairs that satisfy the predicate.
func (m Native[K, V]) Filter(fn func(K, V) bool) Map[K, V] {
	return Native[K, V]{data: maps.Filter(m.data, fn)}
}

// Values returns a slice of all values in the map.
func (m Native[K, V]) Values() slices.Slice[V] {
	res := make([]V, len(m.data))
	i := 0
	for _, v := range m.data {
		res[i] = v
		i++
	}
	return res
}

// Keys returns a slice of all keys in the map.
func (m Native[K, V]) Keys() slices.Slice[K] {
	res := make([]K, len(m.data))
	i := 0
	for k := range m.data {
		res[i] = k
		i++
	}
	return res
}

// Entries returns a slice of all key-value pairs in the map as Entry structs.
func (m Native[K, V]) Entries() slices.Slice[Entry[K, V]] {
	res := make([]Entry[K, V], len(m.data))
	i := 0
	for k, v := range m.data {
		res[i] = Entry[K, V]{K: k, V: v}
		i++
	}
	return res
}
