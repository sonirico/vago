package streams

import (
	"io"
	"iter"
)

type (
	// GroupStream groups consecutive items with the same key into slices.
	// It uses a key extraction function to determine when items should be grouped together.
	GroupStream[T any, K comparable] struct {
		inner      ReadStream[T]
		keyFunc    func(T) K
		buffer     []T
		currentKey K
		err        error
		done       bool
		hasNext    bool
	}
)

// Group creates a new stream that groups consecutive items with the same key.
// The keyFunc is used to extract the grouping key from each item (like Python's itemgetter).
func Group[T any, K comparable](inner ReadStream[T], keyFunc func(T) K) ReadStream[[]T] {
	return &GroupStream[T, K]{
		inner:   inner,
		keyFunc: keyFunc,
	}
}

func (s *GroupStream[T, K]) Next() bool {
	if s.done {
		return false
	}

	// Create a new buffer for the next group instead of reusing
	s.buffer = make([]T, 0)

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
	s.currentKey = groupKey

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

func (s *GroupStream[T, K]) Data() []T {
	return s.buffer
}

func (s *GroupStream[T, K]) Err() error {
	return s.err
}

func (s *GroupStream[T, K]) Close() error {
	return s.inner.Close()
}

func (s *GroupStream[T, K]) Iter() iter.Seq[[]T] {
	return Iter(s)
}

// GroupFactory creates a factory for GroupStream instances.
func GroupFactory[T any, K comparable](
	innerFactory ReadStreamFactory[T],
	keyFunc func(T) K,
) ReadStreamFactory[[]T] {
	return func(rc io.ReadCloser) ReadStream[[]T] {
		return Group(innerFactory(rc), keyFunc)
	}
}

var _ ReadStream[[]any] = new(GroupStream[any, string])
