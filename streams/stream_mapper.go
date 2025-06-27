package streams

import (
	"io"
	"iter"
)

type (
	MapperStream[T, V any] struct {
		inner  ReadStream[T]
		mapper func(T) V
	}
)

func (s *MapperStream[T, V]) Next() bool {
	return s.inner.Next()
}

func (s *MapperStream[T, V]) Data() V {
	return s.mapper(s.inner.Data())
}

func (s *MapperStream[T, V]) Err() error {
	return s.inner.Err()
}

func (s *MapperStream[T, V]) Close() error {
	return s.inner.Close()
}

func (s *MapperStream[T, V]) Iter() iter.Seq[V] {
	return Iter(s)
}

// NewMapperStream creates a new ReadStream that transforms elements from the inner stream
// using the provided mapper function. Each element of type T is converted to type V
// using the mapper function.
//
// This is useful for creating data processing pipelines where you need to transform
// data from one type to another, such as converting strings to uppercase or
// extracting specific fields from structs.
func NewMapperStream[T, V any](inner ReadStream[T], mapper func(T) V) ReadStream[V] {
	return &MapperStream[T, V]{
		inner:  inner,
		mapper: mapper,
	}
}

func NewMapperStreamFactory[T, V any](
	inner ReadStreamFactory[T],
	mapper func(T) V,
) ReadStreamFactory[V] {
	return func(readCloser io.ReadCloser) ReadStream[V] {
		return NewMapperStream[T, V](inner(readCloser), mapper)
	}
}

var _ ReadStream[any] = new(MapperStream[any, any])
