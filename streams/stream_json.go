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

func NewJSONEachRowStream[T any](r io.ReadCloser) *JSONEachRowStream[T] {
	return &JSONEachRowStream[T]{
		r:       r,
		decoder: json.NewDecoder(r),
	}
}

func NewJSONEachRowStreamFactory[T any]() ReadStreamFactory[T] {
	return func(rc io.ReadCloser) ReadStream[T] {
		return NewJSONEachRowStream[T](rc)
	}
}

var _ ReadStream[any] = new(JSONEachRowStream[any])
