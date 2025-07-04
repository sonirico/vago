package streams

import (
	"iter"
	"slices"
)

// Iter converts a ReadStream into an iter.Seq.
func Iter[T any](stream ReadStream[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		defer stream.Close()

		for stream.Next() {
			if !yield(stream.Data()) {
				break
			}
		}
	}
}

// Iter2 converts a ReadStream into an iter.Seq2, yielding both data and error.
func Iter2[T any](stream ReadStream[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		defer stream.Close()

		for stream.Next() {
			if !yield(stream.Data(), nil) {
				break
			}
		}
		if err := stream.Err(); err != nil {
			var x T
			yield(x, err)
		}
		if err := stream.Close(); err != nil {
			var x T
			yield(x, err)
		}
	}
}

// Collect collects all elements from an iter.Seq into a slice.
func Collect[T any](stream iter.Seq[T]) []T {
	return slices.Collect(stream)
}

// SeqKeys collects all keys from an iter.Seq2 into another Seq.
func SeqKeys[K, V any](iter iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range iter {
			if !yield(k) {
				break
			}
		}
	}
}

// SeqValues collects all values from an iter.Seq2 into another Seq.
func SeqValues[K, V any](iter iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range iter {
			if !yield(v) {
				break
			}
		}
	}
}
