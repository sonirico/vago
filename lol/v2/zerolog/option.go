package zerolog

import (
	"context"
	"io"

	rs "github.com/rs/zerolog"

	lol "github.com/sonirico/vago/lol/v2"
)

// TraceHook builds the hook WithTrace installs for one context.
type TraceHook func(ctx context.Context) rs.Hook

type config struct {
	lol.Config
	extraWriters []io.Writer
	traceHook    TraceHook
}

// Option extends what lol.Config cannot say about this backend.
type Option func(*config)

// WithExtraWriter tees every record to w beside lol.Config.Writer. A w
// that implements rs.LevelWriter gets the level, as rs.MultiLevelWriter
// does.
func WithExtraWriter(w io.Writer) Option {
	return func(c *config) { c.extraWriters = append(c.extraWriters, w) }
}

// WithTraceHook makes Logger.WithTrace install hook(ctx).
func WithTraceHook(hook TraceHook) Option {
	return func(c *config) { c.traceHook = hook }
}
