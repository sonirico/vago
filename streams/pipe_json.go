package streams

import (
	"encoding/json"
	"fmt"
	"io"

	"errors"
)

type (
	TransformJSON[T any] struct {
		stream ReadStream[T]
	}
)

func (p *TransformJSON[T]) WriteTo(w io.Writer) (written int64, err error) {
	n, err := w.Write([]byte("["))
	if err != nil {
		return int64(n), err
	} else {
		written += int64(n)
	}

	writeComma := false

	for p.stream.Next() {
		if err = p.stream.Err(); err != nil {
			if !errors.Is(err, io.EOF) {
				err = fmt.Errorf("stream err: %w", err)
			}
			return
		}

		if !writeComma {
			writeComma = true
		} else {
			n, err = w.Write([]byte(","))
			if err != nil {
				return int64(n), err
			} else {
				written += int64(n)
			}
		}

		data, err := json.Marshal(p.stream.Data())

		if err != nil {
			return written, err
		}

		n, err := w.Write(data)
		if err != nil {
			return written, err
		}

		written += int64(n)
	}

	n, err = w.Write([]byte("]"))
	if err != nil && !errors.Is(err, io.EOF) {
		return int64(n), err
	} else {
		written += int64(n)
	}

	if err = p.stream.Err(); err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
		} else {
			err = fmt.Errorf("final stream error: %w", err)
		}
	}

	return
}

// JSONTransform creates a Transform that converts a ReadStream to JSON format.
// The resulting Transform can be used with WriteTo to output the stream data
// as a JSON array to any io.Writer.
//
// This is useful for converting structured data to JSON for APIs, file output,
// or network transmission.
func JSONTransform[T any](r ReadStream[T]) Transform[T] {
	return &TransformJSON[T]{
		stream: r,
	}
}

// PipeJSON writes the JSON representation of each item in the stream to the provided writer.
func PipeJSON[T any](stream ReadStream[T], w io.Writer) (int64, error) {
	return JSONTransform(stream).WriteTo(w)
}
