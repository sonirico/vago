package streams

import (
	"encoding/json"
	"fmt"
	"io"
	"iter"
)

type JSONEachRowStream[T any] struct {
	current T

	r io.ReadCloser

	decoder *json.Decoder

	err error
}

func (s *JSONEachRowStream[T]) Next() bool {
	s.err = s.decoder.Decode(&s.current)

	if s.err != nil {
		if errClose := s.r.Close(); errClose != nil {
			s.err = fmt.Errorf("%w: %s", s.err, errClose.Error())
		}

		return false
	}

	return true
}

func (s *JSONEachRowStream[T]) Data() T {
	return s.current
}

func (s *JSONEachRowStream[T]) Err() error {
	return s.err
}

func (s *JSONEachRowStream[T]) Close() error {
	return s.r.Close()
}

func (s *JSONEachRowStream[T]) Iter() iter.Seq[T] {
	return Iter(s)
}

func (s *JSONEachRowStream[T]) Iter2() iter.Seq2[T, error] {
	return Iter2(s)
}

// JSON creates a new JSONEachRowStream that reads from the provided io.ReadCloser.
// The stream decodes JSON objects from the input, where each object represents a row.
// This is useful for processing JSON data in a row-oriented manner, such as reading
func JSON[T any](r io.ReadCloser) *JSONEachRowStream[T] {
	return &JSONEachRowStream[T]{
		r:       r,
		decoder: json.NewDecoder(r),
	}
}

var _ ReadStream[any] = new(JSONEachRowStream[any])
