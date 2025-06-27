package streams

import (
	"io"
	"iter"
)

type (
	FlattenerStream[T any] struct {
		inner  ReadStream[[]T]
		data   []T
		cursor int
		next   bool
	}
)

func (s *FlattenerStream[T]) Next() bool {
	if len(s.data) == 0 || s.cursor == len(s.data)-1 {
		if s.next = s.inner.Next(); s.next {
			s.data = s.inner.Data()
			s.cursor = 0

			if len(s.data) == 0 {
				return s.Next()
			}
		} else {
			s.data = s.data[:0]
		}
	} else {
		s.cursor++
	}

	if !s.next && len(s.data) == 0 {
		return false
	}

	return true
}

func (s *FlattenerStream[T]) Data() T {
	return s.data[s.cursor]
}

func (s *FlattenerStream[T]) Err() error {
	return s.inner.Err()
}

func (s *FlattenerStream[T]) Close() error {
	return s.inner.Close()
}

func (s *FlattenerStream[T]) Iter() iter.Seq[T] {
	return Iter(s)
}

func (s *FlattenerStream[T]) Iter2() iter.Seq2[T, error] {
	return Iter2(s)
}

// NewFlattenerStream creates a new ReadStream that flattens slices from the inner stream.
// It takes a ReadStream[[]T] and converts it to ReadStream[T] by emitting each
// element from the inner slices individually.
//
// This is useful when you have a stream of slices and want to process each
// individual element, such as flattening batched data or expanding grouped results.
func NewFlattenerStream[T any](inner ReadStream[[]T]) ReadStream[T] {
	return &FlattenerStream[T]{
		inner: inner,
	}
}

func NewFlattenerStreamFactory[T any](
	inner ReadStreamFactory[[]T],
) ReadStreamFactory[T] {
	return func(readCloser io.ReadCloser) ReadStream[T] {
		return NewFlattenerStream[T](inner(readCloser))
	}
}

var _ ReadStream[any] = new(FlattenerStream[any])
