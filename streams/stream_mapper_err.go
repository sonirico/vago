package streams

import "iter"

type (
	MapperStreamErr[T, V any] struct {
		err    error
		cur    V
		inner  ReadStream[T]
		mapper func(T) (V, error)
	}
)

func (s *MapperStreamErr[T, V]) Next() bool {
	if !s.inner.Next() {
		return false
	}

	s.cur, s.err = s.mapper(s.inner.Data())
	return true
}

func (s *MapperStreamErr[T, V]) Data() V {
	return s.cur
}

func (s *MapperStreamErr[T, V]) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.inner.Err()
}

func (s *MapperStreamErr[T, V]) Close() error {
	return s.inner.Close()
}

func (s *MapperStreamErr[T, V]) Iter() iter.Seq[V] {
	return Iter(s)
}

func (s *MapperStreamErr[T, V]) Iter2() iter.Seq2[V, error] {
	return Iter2(s)
}

func NewMapperStreamErr[T, V any](inner ReadStream[T], mapper func(T) (V, error)) ReadStream[V] {
	return &MapperStreamErr[T, V]{
		inner:  inner,
		mapper: mapper,
	}
}

var _ ReadStream[any] = new(MapperStreamErr[any, any])
