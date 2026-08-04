package streams

import (
	"context"
	"iter"
)

type StreamChannel[T any] struct {
	ch      <-chan T
	current T
	err     error
}

func Channel[T any](ch <-chan T) ReadStream[T] {
	return &StreamChannel[T]{
		ch: ch,
	}
}

func (s *StreamChannel[T]) Next(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		s.err = ctx.Err()
		return false
	case v, ok := <-s.ch:
		s.current = v
		return ok
	}
}

func (s *StreamChannel[T]) Data() T {
	return s.current
}

func (s *StreamChannel[T]) Err() error {
	return s.err
}

func (s *StreamChannel[T]) Close() error {
	return nil
}

func (s *StreamChannel[T]) Iter(ctx context.Context) iter.Seq[T] {
	return Iter(ctx, s)
}

var _ ReadStream[any] = new(StreamChannel[any])
