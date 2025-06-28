package streams

import (
	"errors"
	"io"
)

// Reduce applies a reduction function to a ReadStream, accumulating results in a map.
// The function takes the current map and an item from the stream, returning a new map.
// It returns the final accumulated map or an error if the stream encounters one.
func Reduce[T, R any](
	s ReadStream[T],
	fn func(R, T) R,
	initial R,
) (R, error) {
	res := initial
	for s.Next() {
		if err := s.Err(); err != nil {
			var x R
			return x, err
		}

		res = fn(res, s.Data())
	}

	if err := s.Err(); err != nil && !errors.Is(err, io.EOF) {
		var x R
		return x, err
	}

	return res, nil
}

// ReduceSlice applies a reduction function to a ReadStream, accumulating results in a slice.
// The function takes the current slice and an item from the stream, returning a new slice.
// It returns the final accumulated slice or an error if the stream encounters one.
func ReduceSlice[T any](
	s ReadStream[T],
	fn func([]T, T) []T,
) ([]T, error) {
	return Reduce(
		s,
		func(acc []T, item T) []T {
			acc = fn(acc, item)
			return acc
		},
		make([]T, 0),
	)
}

// ReduceMap applies a reduction function to a ReadStream, accumulating results in a map.
// The function takes the current map and an item from the stream, returning a new map.
// It returns the final accumulated map or an error if the stream encounters one.
func ReduceMap[T any, K comparable, V any](
	s ReadStream[T],
	fn func(map[K]V, T) map[K]V,
) (map[K]V, error) {
	return Reduce(
		s,
		func(acc map[K]V, item T) map[K]V {
			acc = fn(acc, item)
			return acc
		},
		make(map[K]V),
	)
}
