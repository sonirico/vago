package streams

import (
	"io"
	"iter"
)

type (
	FilterStream[T any] struct {
		inner     ReadStream[T]
		predicate func(T) bool
		current   T
		hasData   bool
	}
)

func (s *FilterStream[T]) Next() bool {
	for s.inner.Next() {
		data := s.inner.Data()
		if s.predicate(data) {
			s.current = data
			s.hasData = true
			return true
		}
	}
	s.hasData = false
	return false
}

func (s *FilterStream[T]) Data() T {
	if !s.hasData {
		var zero T
		return zero
	}
	return s.current
}

func (s *FilterStream[T]) Err() error {
	return s.inner.Err()
}

func (s *FilterStream[T]) Close() error {
	return s.inner.Close()
}

func (s *FilterStream[T]) Iter() iter.Seq[T] {
	return Iter(s)
}

// Filter creates a new ReadStream that filters elements from the inner stream
// using the provided predicate function. Only elements that satisfy the predicate
// (return true) will be included in the resulting stream.
//
// This is useful for creating data processing pipelines where you need to exclude
// certain elements based on custom criteria.
func Filter[T any](inner ReadStream[T], predicate func(T) bool) ReadStream[T] {
	return &FilterStream[T]{
		inner:     inner,
		predicate: predicate,
	}
}

func FilterFactory[T any](
	inner ReadStreamFactory[T],
	predicate func(T) bool,
) ReadStreamFactory[T] {
	return func(readCloser io.ReadCloser) ReadStream[T] {
		return Filter[T](inner(readCloser), predicate)
	}
}

var _ ReadStream[any] = new(FilterStream[any])
