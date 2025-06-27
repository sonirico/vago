package streams

import "io"

type (
	ReadStream[T any] interface {
		Next() bool
		Data() T
		Err() error
		Close() error
	}

	WriteStream[T any] interface {
		Write(T) (int64, error)
		Flush() error
		Err() error
		Close() error
	}

	ReadStreamFactory[T any] func(closer io.ReadCloser) ReadStream[T]

	WriteStreamFactory[T any] func(closer io.WriteCloser) WriteStream[T]

	TransformFunc[T, V any] func(T) V
)
