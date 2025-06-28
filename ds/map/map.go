// Package _map provides generic map implementations and interfaces
package _map

import (
	"github.com/sonirico/vago/fp"
	"github.com/sonirico/vago/slices"
	"github.com/sonirico/vago/tuples"
)

type (
	// Map is a generic interface for map operations with keys of type K and values of type V.
	// It defines common operations like getting, setting, checking existence, iteration, and transformations.
	Map[K comparable, V any] interface {
		// Get retrieves a value by its key. Returns the value and whether the key exists.
		Get(K) (V, bool)

		// Has checks if the map contains the specified key.
		Has(K) bool

		// Set adds or replaces a key-value pair in the map.
		Set(K, V)

		// Range iterates over each key-value pair in the map, with its index.
		// The iteration stops if the provided function returns false.
		Range(fn func(K, V, int) bool)

		// Delete removes a key-value pair from the map.
		Delete(K)

		// GetOrSet retrieves a value or sets a default if the key doesn't exist.
		// Returns the value (existing or default) and whether the key already existed.
		GetOrSet(K, V) (V, bool)

		// Map applies a transformation function to each key-value pair and returns a new map.
		Map(fn func(K, V) (K, V)) Map[K, V]

		// FilterMap applies a function that may filter out or transform key-value pairs.
		// Returns a new map containing only the key-value pairs for which the function
		// returns a Some result.
		FilterMap(fn func(K, V) fp.Option[tuples.Tuple2[K, V]]) Map[K, V]

		// Filter returns a new map containing only the key-value pairs that satisfy the predicate.
		Filter(fn func(K, V) bool) Map[K, V]

		// Keys returns a slice of all keys in the map.
		Keys() slices.Slice[K]

		// Values returns a slice of all values in the map.
		Values() slices.Slice[V]

		// Entries returns a slice of all key-value pairs in the map as Entry structs.
		Entries() slices.Slice[Entry[K, V]]
	}
)
