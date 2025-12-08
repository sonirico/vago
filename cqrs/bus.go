package cqrs

import (
	"context"

	"github.com/sonirico/vago/rp"
)

type bus interface {
	id() string
	topic() string
	codec() Codec
	publish(ctx context.Context, msg rp.Msg) error
	subscribe(ctx context.Context, handler rp.ConsumerHandler) error
	close()
	hasHandlers() bool
}
