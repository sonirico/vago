package _map

import (
	"sync"

	"github.com/sonirico/gozo/fp"
	"github.com/sonirico/gozo/slices"
	"github.com/sonirico/gozo/tuples"
)

type (
	// MapInner is an alias for the Map interface, used in the Concurrent implementation.
	MapInner[K comparable, V any] Map[K, V]

	// Concurrent provides a thread-safe wrapper around any Map implementation.
	// It implements the Map interface by delegating to the inner map with proper synchronization.
	Concurrent[K comparable, V any] struct {
		L              sync.RWMutex // Mutex for thread-safe access
		MapInner[K, V]              // The underlying map implementation
	}
)

// NewConcurrent creates a new thread-safe map wrapper around the provided map implementation.
// This makes any Map implementation safe for concurrent use.
func NewConcurrent[K comparable, V any](inner Map[K, V]) *Concurrent[K, V] {
	return &Concurrent[K, V]{MapInner: inner}
}

// Get retrieves a value by its key with read lock protection.
// Returns the value and a boolean indicating whether the key exists.
func (m *Concurrent[K, V]) Get(k K) (v V, ok bool) {
	m.L.RLock()
	v, ok = m.MapInner.Get(k)
	m.L.RUnlock()
	return
}

// Has checks if the map contains the specified key with read lock protection.
func (m *Concurrent[K, V]) Has(k K) (ok bool) {
	m.L.RLock()
	_, ok = m.MapInner.Get(k)
	m.L.RUnlock()
	return
}

// Set adds or replaces a key-value pair in the map with write lock protection.
func (m *Concurrent[K, V]) Set(k K, v V) {
	m.L.Lock()
	m.MapInner.Set(k, v)
	m.L.Unlock()
	return
}

// Range iterates over each key-value pair in the map with read lock protection.
// The iteration stops if the provided function returns false.
func (m *Concurrent[K, V]) Range(fn func(K, V, int) bool) {
	m.L.RLock()
	defer m.L.RUnlock()
	m.MapInner.Range(fn)
}

// Delete removes a key-value pair from the map with write lock protection.
func (m *Concurrent[K, V]) Delete(k K) {
	m.L.Lock()
	m.MapInner.Delete(k)
	m.L.Unlock()
}

// GetOrSet retrieves a value or sets a default if the key doesn't exist, with write lock protection.
// Returns the value (existing or default) and whether the key already existed.
func (m *Concurrent[K, V]) GetOrSet(k K, def V) (v V, ok bool) {
	m.L.Lock()
	v, ok = m.MapInner.GetOrSet(k, def)
	m.L.Unlock()
	return
}

// Map applies a transformation function to each key-value pair and returns a new concurrent map.
// Uses read lock protection during the transformation.
func (m *Concurrent[K, V]) Map(fn func(K, V) (K, V)) Map[K, V] {
	m.L.RLock()
	defer m.L.RUnlock()
	return &Concurrent[K, V]{MapInner: m.MapInner.Map(fn)}
}

// FilterMap applies a function that may filter out or transform key-value pairs with read lock protection.
// Returns a new concurrent map containing only the key-value pairs for which the function
// returns a Some result.
func (m *Concurrent[K, V]) FilterMap(
	fn func(K, V) fp.Option[tuples.Tuple2[K, V]],
) Map[K, V] {
	m.L.RLock()
	defer m.L.RUnlock()
	return &Concurrent[K, V]{MapInner: m.MapInner.FilterMap(fn)}
}

// Filter returns a new concurrent map containing only the key-value pairs that satisfy
// the predicate, with read lock protection.
func (m *Concurrent[K, V]) Filter(fn func(K, V) bool) Map[K, V] {
	m.L.RLock()
	defer m.L.RUnlock()
	return &Concurrent[K, V]{MapInner: m.MapInner.Filter(fn)}
}

// Values returns a slice of all values in the map with read lock protection.
func (m *Concurrent[K, V]) Values() slices.Slice[V] {
	m.L.RLock()
	res := m.MapInner.Values()
	m.L.RUnlock()
	return res
}

// Keys returns a slice of all keys in the map with read lock protection.
func (m *Concurrent[K, V]) Keys() slices.Slice[K] {
	m.L.RLock()
	res := m.MapInner.Keys()
	m.L.RUnlock()
	return res
}

// Entries returns a slice of all key-value pairs in the map as Entry structs
// with read lock protection.
func (m *Concurrent[K, V]) Entries() slices.Slice[Entry[K, V]] {
	m.L.RLock()
	res := m.MapInner.Entries()
	m.L.RUnlock()
	return res
}
