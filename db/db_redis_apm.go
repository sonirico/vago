package db

import (
	"context"

	"github.com/go-redis/redis/v8"
	apmgoredis "go.elastic.co/apm/module/apmgoredisv8/v2"
	"go.elastic.co/apm/v2"
)

type hook struct {
	redis.Hook
}

func NewApmRedisHook() redis.Hook {
	return hook{apmgoredis.NewHook()}
}

func (r hook) BeforeProcess(c context.Context, cmd redis.Cmder) (ctx context.Context, err error) {
	ctx, err = r.Hook.BeforeProcess(c, cmd)

	if span := apm.SpanFromContext(ctx); span != nil {
		args := cmd.Args()
		if len(args) >= 1 {
			span.Context.SetLabel("method", args[0])
		}
		if len(args) > 1 {
			span.Context.SetLabel("key", args[1])
		}
	}

	return ctx, err
}
