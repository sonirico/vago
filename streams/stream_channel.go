package streams

import "iter"

type StreamChannel[T any] struct {
	ch      <-chan T
	current T
}

func Channel[T any](ch <-chan T) ReadStream[T] {
	return &StreamChannel[T]{
		ch: ch,
	}
}

func (s *StreamChannel[T]) Next() bool {
	var ok bool
	s.current, ok = <-s.ch
	return ok
}

func (s *StreamChannel[T]) Data() T {
	return s.current
}

func (s *StreamChannel[T]) Err() error {
	return nil
}

func (s *StreamChannel[T]) Close() error {
	return nil
}

func (s *StreamChannel[T]) Iter() iter.Seq[T] {
	return Iter(s)
}

var _ ReadStream[any] = new(StreamChannel[any])
