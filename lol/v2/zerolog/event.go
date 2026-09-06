package zerolog

import (
	"time"

	rs "github.com/rs/zerolog"

	lol "github.com/sonirico/vago/lol/v2"
)

// event wraps one *rs.Event. It is a single pointer, so storing it in a
// lol.Event interface does not allocate, and rs.Event's methods accept a
// nil receiver, which is what a disabled level returns.
type event struct {
	e *rs.Event
}

func (ev event) Str(k, v string) lol.Event               { return event{ev.e.Str(k, v)} }
func (ev event) Int(k string, v int) lol.Event           { return event{ev.e.Int(k, v)} }
func (ev event) Int64(k string, v int64) lol.Event       { return event{ev.e.Int64(k, v)} }
func (ev event) Uint64(k string, v uint64) lol.Event     { return event{ev.e.Uint64(k, v)} }
func (ev event) Float64(k string, v float64) lol.Event   { return event{ev.e.Float64(k, v)} }
func (ev event) Bool(k string, v bool) lol.Event         { return event{ev.e.Bool(k, v)} }
func (ev event) Dur(k string, v time.Duration) lol.Event { return event{ev.e.Dur(k, v)} }
func (ev event) Time(k string, v time.Time) lol.Event    { return event{ev.e.Time(k, v)} }
func (ev event) Err(err error) lol.Event                 { return event{ev.e.Err(err)} }
func (ev event) Any(k string, v any) lol.Event           { return event{ev.e.Interface(k, v)} }
func (ev event) Msg(msg string)                          { ev.e.Msg(msg) }
func (ev event) Msgf(format string, args ...any)         { ev.e.Msgf(format, args...) }
