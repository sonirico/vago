package streams

import (
	"io"
)

type (
	Transform[T any] interface {
		io.WriterTo
	}
)
