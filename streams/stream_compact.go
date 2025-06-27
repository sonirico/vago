package streams

import (
	"io"
	"iter"
)

type (
	// CompactStream groups consecutive items with the same key into slices.
	// It uses a key extraction function to determine when items should be grouped together.
	CompactStream[T any, K comparable] struct {
		inner      ReadStream[T]
		keyFunc    func(T) K
		buffer     []T
		currentKey *K
		err        error
		done       bool
		hasNext    bool
	}
)

// NewCompactStream creates a new stream that groups consecutive items with the same key.
// The keyFunc is used to extract the grouping key from each item (like Python's itemgetter).
func NewCompactStream[T any, K comparable](inner ReadStream[T], keyFunc func(T) K) ReadStream[[]T] {
	return &CompactStream[T, K]{
		inner:   inner,
		keyFunc: keyFunc,
	}
}

func (s *CompactStream[T, K]) Next() bool {
	if s.done {
		return false
	}

	// Clear the buffer for the next group
	s.buffer = s.buffer[:0]

	// If we don't have a next item, we're done
	if !s.hasNext && !s.inner.Next() {
		if err := s.inner.Err(); err != nil {
			s.err = err
		}
		s.done = true
		return false
	}

	// Get the first item of the group
	var firstItem T
	var groupKey K

	if s.hasNext {
		// We already have the first item from the previous iteration
		firstItem = s.inner.Data()
		groupKey = s.keyFunc(firstItem)
		s.hasNext = false
	} else {
		firstItem = s.inner.Data()
		groupKey = s.keyFunc(firstItem)
	}

	s.buffer = append(s.buffer, firstItem)
	s.currentKey = &groupKey

	// Collect all consecutive items with the same key
	for s.inner.Next() {
		item := s.inner.Data()
		itemKey := s.keyFunc(item)

		if itemKey != groupKey {
			// Different key found, this item belongs to the next group
			s.hasNext = true
			break
		}

		s.buffer = append(s.buffer, item)
	}

	// Check if inner stream ended with an error
	if !s.hasNext {
		if err := s.inner.Err(); err != nil {
			s.err = err
		}
	}

	return len(s.buffer) > 0
}

func (s *CompactStream[T, K]) Data() []T {
	return s.buffer
}

func (s *CompactStream[T, K]) Err() error {
	return s.err
}

func (s *CompactStream[T, K]) Close() error {
	return s.inner.Close()
}

func (s *CompactStream[T, K]) Iter() iter.Seq[[]T] {
	return Iter(s)
}

// NewCompactStreamFactory creates a factory for CompactStream instances.
func NewCompactStreamFactory[T any, K comparable](
	innerFactory ReadStreamFactory[T],
	keyFunc func(T) K,
) ReadStreamFactory[[]T] {
	return func(rc io.ReadCloser) ReadStream[[]T] {
		return NewCompactStream(innerFactory(rc), keyFunc)
	}
}

var _ ReadStream[[]any] = new(CompactStream[any, string])
