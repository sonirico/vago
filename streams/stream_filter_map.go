package streams

import (
	"iter"

	"github.com/sonirico/vago/fp"
)

type (
	FilterMapStream[T, R any] struct {
		inner     ReadStream[T]
		predicate func(T) (R, bool)
		current   R
		hasData   bool
	}
)

func (s *FilterMapStream[T, R]) Next() bool {
	for s.inner.Next() {
		data := s.inner.Data()
		if x, ok := s.predicate(data); ok {
			s.current = x
			s.hasData = true
			return true
		}
	}
	s.hasData = false
	return false
}

func (s *FilterMapStream[T, R]) Data() R {
	if !s.hasData {
		var zero R
		return zero
	}
	return s.current
}

func (s *FilterMapStream[T, R]) Err() error {
	return s.inner.Err()
}

func (s *FilterMapStream[T, R]) Close() error {
	return s.inner.Close()
}

func (s *FilterMapStream[T, R]) Iter() iter.Seq[R] {
	return Iter(s)
}

// FilterMap creates a new ReadStream that filters elements from the inner stream
// using the provided predicate function that returns a value of type R and a boolean.
// Only elements that return true from the predicate will be included in the resulting stream.
// This is useful for creating data processing pipelines where you need to filter and transform
// data in a single step, such as filtering out invalid entries while also transforming valid ones.
func FilterMap[T, R any](inner ReadStream[T], predicate func(T) (R, bool)) ReadStream[R] {
	return &FilterMapStream[T, R]{
		inner:     inner,
		predicate: predicate,
	}
}

// FilterMapOpt creates a new ReadStream that filters elements from the inner stream
// using the provided predicate function that returns an Option type. Only elements
// that return Some from the predicate will be included in the resulting stream.
func FilterMapOpt[T, R any](inner ReadStream[T], predicate func(T) fp.Option[R]) ReadStream[R] {
	return &FilterMapStream[T, R]{
		inner: inner,
		predicate: func(t T) (R, bool) {
			return predicate(t).Unwrap()
		},
	}
}

var _ ReadStream[any] = new(FilterStream[any])
