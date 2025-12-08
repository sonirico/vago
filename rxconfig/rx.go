package rxconfig

import (
	"context"
)

// Change represents a change in rx config
type Change[T any] struct {
	Prev    Config[T]
	Next    Config[T]
	Version int64
}

type Observable[T any] interface {
	Watch(ctx context.Context)
	Changes() <-chan Change[T]
	Get(ctx context.Context) (Config[T], error)
	Put(ctx context.Context, x T) (Config[T], error)
}
