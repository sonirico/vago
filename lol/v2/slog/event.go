package slog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	lol "github.com/sonirico/vago/lol/v2"
)

// afterMsg is what an event does once written: nothing, exit, or panic.
type afterMsg uint8

const (
	noExit afterMsg = iota
	exitProcess
	panicMsg
)

// event accumulates attrs until Msg. A disabled level hands out a nil
// *event, which every method accepts and ignores, so the chain allocates
// nothing; enabled events come from a pool and go back on Msg.
type event struct {
	log   *slog.Logger
	lvl   slog.Level
	after afterMsg
	attrs []slog.Attr
}

var pool = sync.Pool{New: func() any { return &event{attrs: make([]slog.Attr, 0, 8)} }}

func newEvent(log *slog.Logger, lvl slog.Level, after afterMsg) *event {
	ev := pool.Get().(*event)
	ev.log, ev.lvl, ev.after = log, lvl, after
	return ev
}

func (ev *event) add(a slog.Attr) lol.Event {
	if ev == nil {
		return ev
	}
	ev.attrs = append(ev.attrs, a)
	return ev
}

func (ev *event) Str(k, v string) lol.Event               { return ev.add(slog.String(k, v)) }
func (ev *event) Int(k string, v int) lol.Event           { return ev.add(slog.Int(k, v)) }
func (ev *event) Int64(k string, v int64) lol.Event       { return ev.add(slog.Int64(k, v)) }
func (ev *event) Uint64(k string, v uint64) lol.Event     { return ev.add(slog.Uint64(k, v)) }
func (ev *event) Float64(k string, v float64) lol.Event   { return ev.add(slog.Float64(k, v)) }
func (ev *event) Bool(k string, v bool) lol.Event         { return ev.add(slog.Bool(k, v)) }
func (ev *event) Dur(k string, v time.Duration) lol.Event { return ev.add(slog.Duration(k, v)) }
func (ev *event) Time(k string, v time.Time) lol.Event    { return ev.add(slog.Time(k, v)) }
func (ev *event) Any(k string, v any) lol.Event           { return ev.add(slog.Any(k, v)) }

func (ev *event) Err(err error) lol.Event {
	if err == nil {
		return ev
	}
	return ev.add(slog.String("error", err.Error()))
}

func (ev *event) Msg(msg string) {
	if ev == nil {
		return
	}
	ev.log.LogAttrs(context.Background(), ev.lvl, msg, ev.attrs...)
	after := ev.after
	ev.log, ev.attrs = nil, ev.attrs[:0]
	pool.Put(ev)
	switch after {
	case exitProcess:
		os.Exit(1)
	case panicMsg:
		panic(msg)
	}
}

func (ev *event) Msgf(format string, args ...any) {
	if ev == nil {
		return
	}
	ev.Msg(fmt.Sprintf(format, args...))
}
