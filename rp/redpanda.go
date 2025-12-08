package rp

import (
	"context"
)

type (
	common interface {
		Ping(ctx context.Context) error
		Close()
	}

	Producer interface {
		common

		Publish(ctx context.Context, msg Msg) error
		PublishAsync(ctx context.Context, msg Msg, fn func(Msg, error)) error
		Flush(ctx context.Context) error
	}

	Consumer interface {
		common

		Subscribe(ctx context.Context, h ConsumerHandler) error
	}
)
