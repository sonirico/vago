// Package slog is the lol.Logger backend over the standard library's
// log/slog. New builds a handler from lol.Config; Wrap adopts a
// *slog.Logger the process already has.
package slog

import (
	"context"
	"log/slog"

	lol "github.com/sonirico/vago/lol/v2"
)

type logger struct {
	log *slog.Logger
}

// New builds a lol.Logger from cfg. EnvLocal writes text, every other Env
// writes JSON. Trace maps to slog's debug minus 4 and fatal, panic to
// error plus 4 and plus 8; a level guard applies to Fatal and Panic like
// any other event, and both still exit or panic after writing.
func New(cfg lol.Config) lol.Logger {
	opts := &slog.HandlerOptions{Level: level(cfg.Level)}
	if cfg.TimeFormat != "" {
		opts.ReplaceAttr = timeFormat(cfg.TimeFormat)
	}
	var h slog.Handler
	if cfg.Env == lol.EnvLocal {
		h = slog.NewTextHandler(cfg.Writer, opts)
	} else {
		h = slog.NewJSONHandler(cfg.Writer, opts)
	}
	return Wrap(slog.New(h).With(attrs(cfg.Fields)...))
}

// Wrap adopts l as a lol.Logger.
func Wrap(l *slog.Logger) lol.Logger {
	return logger{log: l}
}

func (l logger) Trace() lol.Event { return l.event(levelTrace, noExit) }
func (l logger) Debug() lol.Event { return l.event(slog.LevelDebug, noExit) }
func (l logger) Info() lol.Event  { return l.event(slog.LevelInfo, noExit) }
func (l logger) Warn() lol.Event  { return l.event(slog.LevelWarn, noExit) }
func (l logger) Error() lol.Event { return l.event(slog.LevelError, noExit) }
func (l logger) Fatal() lol.Event { return l.event(levelFatal, exitProcess) }
func (l logger) Panic() lol.Event { return l.event(levelPanic, panicMsg) }

func (l logger) With(fields lol.Fields) lol.Logger {
	return logger{log: l.log.With(attrs(fields)...)}
}

func (l logger) WithTrace(context.Context) lol.Logger {
	return l
}

func (l logger) event(lvl slog.Level, after afterMsg) lol.Event {
	if !l.log.Enabled(context.Background(), lvl) {
		return (*event)(nil)
	}
	return newEvent(l.log, lvl, after)
}

const (
	levelTrace = slog.LevelDebug - 4
	levelFatal = slog.LevelError + 4
	levelPanic = slog.LevelError + 8
)

func level(l lol.Level) slog.Level {
	switch l {
	case lol.LevelTrace:
		return levelTrace
	case lol.LevelDebug:
		return slog.LevelDebug
	case lol.LevelInfo:
		return slog.LevelInfo
	case lol.LevelWarn:
		return slog.LevelWarn
	case lol.LevelError:
		return slog.LevelError
	case lol.LevelFatal:
		return levelFatal
	case lol.LevelPanic:
		return levelPanic
	default:
		return slog.LevelInfo
	}
}

func attrs(fields lol.Fields) []any {
	out := make([]any, 0, len(fields))
	for k, v := range fields {
		out = append(out, slog.Any(k, v))
	}
	return out
}

func timeFormat(layout string) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if len(groups) == 0 && a.Key == slog.TimeKey {
			return slog.String(slog.TimeKey, a.Value.Time().Format(layout))
		}
		return a
	}
}
