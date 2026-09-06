package lol

import (
	"context"
	"time"
)

// Discard is a Logger that drops everything and allocates nothing.
var Discard Logger = NewNoopLogger()

type noopLogger struct{}

// NewNoopLogger returns a Logger that drops everything.
func NewNoopLogger() Logger {
	return noopLogger{}
}

func (noopLogger) Trace() Event                       { return noopEvent{} }
func (noopLogger) Debug() Event                       { return noopEvent{} }
func (noopLogger) Info() Event                        { return noopEvent{} }
func (noopLogger) Warn() Event                        { return noopEvent{} }
func (noopLogger) Error() Event                       { return noopEvent{} }
func (noopLogger) Fatal() Event                       { return noopEvent{} }
func (noopLogger) Panic() Event                       { return noopEvent{} }
func (l noopLogger) With(Fields) Logger               { return l }
func (l noopLogger) WithTrace(context.Context) Logger { return l }

type noopEvent struct{}

func (e noopEvent) Str(string, string) Event        { return e }
func (e noopEvent) Int(string, int) Event           { return e }
func (e noopEvent) Int64(string, int64) Event       { return e }
func (e noopEvent) Uint64(string, uint64) Event     { return e }
func (e noopEvent) Float64(string, float64) Event   { return e }
func (e noopEvent) Bool(string, bool) Event         { return e }
func (e noopEvent) Dur(string, time.Duration) Event { return e }
func (e noopEvent) Time(string, time.Time) Event    { return e }
func (e noopEvent) Err(error) Event                 { return e }
func (e noopEvent) Any(string, any) Event           { return e }
func (noopEvent) Msg(string)                        {}
func (noopEvent) Msgf(string, ...any)               {}
