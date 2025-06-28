package streams

import (
	"io"
	"iter"
)

type (
	BatchStream[T any] struct {
		inner     ReadStream[T]
		batchSize int
		buffer    []T
		err       error
		done      bool
	}
)

// Batch creates a new batch-oriented stream that reads items from
// `inner` in chunks of `batchSize` and returns them as []T.
func Batch[T any](inner ReadStream[T], batchSize int) ReadStream[[]T] {
	return &BatchStream[T]{
		inner:     inner,
		batchSize: batchSize,
	}
}

func (s *BatchStream[T]) Next() bool {
	if s.done {
		return false
	}

	// Create a new buffer for each batch instead of reusing
	s.buffer = make([]T, 0, s.batchSize)

	for len(s.buffer) < s.batchSize {
		if !s.inner.Next() {
			// If the inner stream ended due to an error, capture it.
			if err := s.inner.Err(); err != nil {
				s.err = err
			}
			break
		}
		s.buffer = append(s.buffer, s.inner.Data())
	}

	// If we didn't get any items, we're done.
	if len(s.buffer) == 0 {
		s.done = true
		return false
	}

	return true
}

func (s *BatchStream[T]) Data() []T {
	return s.buffer
}

func (s *BatchStream[T]) Err() error {
	return s.err
}

func (s *BatchStream[T]) Close() error {
	return s.inner.Close()
}

func (s *BatchStream[T]) Iter() iter.Seq[[]T] {
	return Iter(s)
}

func BatchFactory[T any](
	innerFactory ReadStreamFactory[T],
	batchSize int,
) ReadStreamFactory[[]T] {
	return func(rc io.ReadCloser) ReadStream[[]T] {
		return Batch(innerFactory(rc), batchSize)
	}
}

var _ ReadStream[[]any] = new(BatchStream[any])
