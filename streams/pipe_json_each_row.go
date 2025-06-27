package streams

import (
	"encoding/json"
	"fmt"
	"io"

	"errors"
)

type (
	TransformJSONEachRow[T any] struct {
		stream ReadStream[T]
	}
)

func (p *TransformJSONEachRow[T]) WriteTo(w io.Writer) (written int64, err error) {
	for p.stream.Next() {
		if err = p.stream.Err(); err != nil {
			if !errors.Is(err, io.EOF) {
				err = fmt.Errorf("stream err: %w", err)
			}
			return
		}

		var (
			data []byte
			n    int
		)

		data, err = json.Marshal(p.stream.Data())

		if err != nil {
			err = fmt.Errorf("json marshal: %w", err)
			return
		}

		n, err = w.Write(data)
		if err != nil {
			err = fmt.Errorf("json write: %w", err)
			return
		}

		written += int64(n)

		n, err = w.Write([]byte("\n"))
		if err != nil {
			err = fmt.Errorf("json write carriage return: %w", err)
			return
		}

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

// JSONEachRowTransform creates a Transform that converts a ReadStream to JSON-lines format.
// Each element in the stream becomes a separate JSON object on its own line,
// which is useful for streaming JSON processing and log formats.
//
// This format is also known as NDJSON (Newline Delimited JSON) and is commonly
// used for streaming APIs and log processing systems.
func JSONEachRowTransform[T any](stream ReadStream[T]) Transform[T] {
	return &TransformJSONEachRow[T]{
		stream: stream,
	}
}

func PipeJSONEachRow[T any](stream ReadStream[T], w io.Writer) (int64, error) {
	return JSONEachRowTransform(stream).WriteTo(w)
}
