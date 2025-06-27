package streams

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"slices"

	"errors"
)

func ReadAllBytes[T any](transform Transform[T]) ([]byte, error) {
	buf := new(bytes.Buffer)

	if _, err := transform.WriteTo(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Consume[T any](stream ReadStream[T]) ([]T, error) {
	var res []T

	for stream.Next() {
		res = append(res, stream.Data())
	}

	err := stream.Err()

	if err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	return res, nil
}

func ConsumeErrSkip[T any](stream ReadStream[T]) []T {
	var res []T

	for stream.Next() {
		if err := stream.Err(); err == nil {
			res = append(res, stream.Data())
		}
	}

	return res
}

// WriteAll writes all items from a slice to a WriteStream
// Returns the total number of bytes written and any error
func WriteAll[T any](stream WriteStream[T], items []T) (int64, error) {
	return WriteSeq(stream, slices.Values(items))
}

// WriteSeq writes all items from an iter.Seq to a WriteStream
// Returns the total number of bytes written and any error
func WriteSeq[T any](stream WriteStream[T], items iter.Seq[T]) (int64, error) {
	bytesWritten := int64(0)

	for v := range items {
		n, err := stream.Write(v)
		if err != nil {
			return 0, fmt.Errorf("write error: %w", err)
		}
		if n == 0 {
			continue // Skip if nothing was written
		}
		bytesWritten += n
	}
	if err := stream.Flush(); err != nil {
		return 0, fmt.Errorf("flush error: %w", err)
	}
	return bytesWritten, nil
}

// WriteSeqKeys writes all keys from an iter.Seq2 to a WriteStream
// Returns the total number of bytes written and any error
func WriteSeqKeys[K, V any](stream WriteStream[K], items iter.Seq2[K, V]) (int64, error) {
	return WriteSeq(stream, SeqKeys(items))
}

// WriteSeqValues writes all values from an iter.Seq2 to a WriteStream
// Returns the total number of bytes written and any error
func WriteSeqValues[K, V any](stream WriteStream[V], items iter.Seq2[K, V]) (int64, error) {
	return WriteSeq(stream, SeqValues(items))
}

// Pipe copies all items from a ReadStream to a WriteStream
// Returns the total number of bytes written and any error
func Pipe[T any](src ReadStream[T], dst WriteStream[T]) (int64, error) {
	var totalBytes int64

	for src.Next() {
		if err := src.Err(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return totalBytes, fmt.Errorf("read error: %w", err)
		}

		n, err := dst.Write(src.Data())
		if err != nil {
			return totalBytes, fmt.Errorf("write error: %w", err)
		}
		totalBytes += n
	}

	if err := dst.Flush(); err != nil {
		return totalBytes, fmt.Errorf("flush error: %w", err)
	}

	return totalBytes, nil
}

// Multiplex copies all items from a ReadStream to multiple WriteStreams
// Returns a slice with bytes written to each destination and any error
func Multiplex[T any](src ReadStream[T], destinations ...WriteStream[T]) ([]int64, error) {
	if len(destinations) == 0 {
		return []int64{}, nil
	}

	bytesWritten := make([]int64, len(destinations))

	for src.Next() {
		if err := src.Err(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return bytesWritten, fmt.Errorf("read error: %w", err)
		}

		data := src.Data()
		for i, dst := range destinations {
			n, err := dst.Write(data)
			if err != nil {
				return bytesWritten, fmt.Errorf("write error to destination %d: %w", i, err)
			}
			bytesWritten[i] += n
		}
	}

	for i, dst := range destinations {
		if err := dst.Flush(); err != nil {
			return bytesWritten, fmt.Errorf("flush error for destination %d: %w", i, err)
		}
	}

	return bytesWritten, nil
}
