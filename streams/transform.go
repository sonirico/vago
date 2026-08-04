package streams

import (
	"context"
	"io"
)

type (
	// Transform writes the elements read from a ReadStream to an io.Writer in a
	// specific wire format. WriteTo takes a ctx (rather than satisfying
	// io.WriterTo) so it can honour cancellation of the underlying read stream.
	Transform[T any] interface {
		WriteTo(ctx context.Context, w io.Writer) (int64, error)
	}
)
