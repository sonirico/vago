package streams

import (
	"context"
	"iter"
)

type (
	MapperStreamErr[T, V any] struct {
		err    error
		cur    V
		inner  ReadStream[T]
		mapper func(T) (V, error)
	}
)

func (s *MapperStreamErr[T, V]) Next(ctx context.Context) bool {
	if !s.inner.Next(ctx) {
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

func (s *MapperStreamErr[T, V]) Iter(ctx context.Context) iter.Seq[V] {
	return Iter(ctx, s)
}

func (s *MapperStreamErr[T, V]) Iter2(ctx context.Context) iter.Seq2[V, error] {
	return Iter2(ctx, s)
}

func MapErr[T, V any](inner ReadStream[T], mapper func(T) (V, error)) ReadStream[V] {
	return &MapperStreamErr[T, V]{
		inner:  inner,
		mapper: mapper,
	}
}

var _ ReadStream[any] = new(MapperStreamErr[any, any])
