// Package lol (lots of logs) is a structured logging interface with typed
// fields and swappable backends. The package holds only the contract; a
// backend lives in its own module so a consumer imports the one it uses:
// lol/v2/zerolog, lol/v2/slog.
//
//	log := zerolog.New(lol.NewConfig(lol.WithLevel(lol.LevelInfo)))
//	log.Info().Str("user", id).Int("attempts", n).Msg("login")
//
// An Event on a disabled level is a shared no-op: the chain above allocates
// nothing when info is off.
package lol

import (
	"context"
	"time"
)

// Logger opens one Event per severity. Fatal exits and Panic panics after
// the event's Msg, as the backends define them.
type Logger interface {
	Trace() Event
	Debug() Event
	Info() Event
	Warn() Event
	Error() Event
	Fatal() Event
	Panic() Event

	// With returns a Logger that stamps fields on every event.
	With(fields Fields) Logger
	// WithTrace returns a Logger that stamps ctx's trace context, when the
	// backend knows how to read one; otherwise it returns the receiver.
	WithTrace(ctx context.Context) Logger
}

// Event is one log record under construction. Every setter returns the
// event so calls chain, and Msg or Msgf writes it. An event on a disabled
// level discards everything and allocates nothing, with one exception the
// language imposes: Msgf's variadic args are boxed at the call site before
// the level is known, exactly as with fmt.
type Event interface {
	Str(key, value string) Event
	Int(key string, value int) Event
	Int64(key string, value int64) Event
	Uint64(key string, value uint64) Event
	Float64(key string, value float64) Event
	Bool(key string, value bool) Event
	Dur(key string, value time.Duration) Event
	Time(key string, value time.Time) Event
	Err(err error) Event
	Any(key string, value any) Event
	Msg(msg string)
	Msgf(format string, args ...any)
}
