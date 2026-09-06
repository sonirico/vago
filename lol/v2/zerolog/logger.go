// Package zerolog is the lol.Logger backend over github.com/rs/zerolog.
// It carries no tracing dependency; lol/v2/zerolog/apm adds Elastic APM
// through the Options this package exposes.
package zerolog

import (
	"context"
	"io"
	"os"
	"time"

	rs "github.com/rs/zerolog"

	lol "github.com/sonirico/vago/lol/v2"
)

type logger struct {
	log       rs.Logger
	traceHook TraceHook
}

// New builds a lol.Logger from cfg. EnvLocal gets zerolog's console writer,
// every other Env writes JSON. Options extend the writer set and the hooks
// the backend alone cannot know about.
func New(cfg lol.Config, opts ...Option) lol.Logger {
	c := config{Config: cfg}
	for _, opt := range opts {
		opt(&c)
	}

	var base io.Writer = c.Writer
	if c.Env == lol.EnvLocal {
		base = rs.ConsoleWriter{Out: c.Writer, NoColor: !isTerminal(c.Writer)}
	}
	out := rs.MultiLevelWriter(append([]io.Writer{base}, c.extraWriters...)...)

	ctx := rs.New(out).Level(level(c.Level)).With()
	for k, v := range c.Fields {
		ctx = ctx.Interface(k, v)
	}
	log := ctx.Logger()
	// rs.TimeFieldFormat is a package global; a layout per logger has to be
	// a hook instead, so two loggers in one process never race on it.
	if c.TimeFormat == "" {
		log = log.With().Timestamp().Logger()
	} else {
		log = log.Hook(timestampHook{layout: c.TimeFormat})
	}
	return logger{log: log, traceHook: c.traceHook}
}

func (l logger) Trace() lol.Event { return event{l.log.Trace()} }
func (l logger) Debug() lol.Event { return event{l.log.Debug()} }
func (l logger) Info() lol.Event  { return event{l.log.Info()} }
func (l logger) Warn() lol.Event  { return event{l.log.Warn()} }
func (l logger) Error() lol.Event { return event{l.log.Error()} }
func (l logger) Fatal() lol.Event { return event{l.log.Fatal()} }
func (l logger) Panic() lol.Event { return event{l.log.Panic()} }

func (l logger) With(fields lol.Fields) lol.Logger {
	ctx := l.log.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return logger{log: ctx.Logger(), traceHook: l.traceHook}
}

func (l logger) WithTrace(ctx context.Context) lol.Logger {
	if l.traceHook == nil {
		return l
	}
	return logger{log: l.log.Hook(l.traceHook(ctx)), traceHook: l.traceHook}
}

func level(l lol.Level) rs.Level {
	switch l {
	case lol.LevelTrace:
		return rs.TraceLevel
	case lol.LevelDebug:
		return rs.DebugLevel
	case lol.LevelInfo:
		return rs.InfoLevel
	case lol.LevelWarn:
		return rs.WarnLevel
	case lol.LevelError:
		return rs.ErrorLevel
	case lol.LevelFatal:
		return rs.FatalLevel
	case lol.LevelPanic:
		return rs.PanicLevel
	default:
		return rs.InfoLevel
	}
}

type timestampHook struct {
	layout string
}

func (h timestampHook) Run(e *rs.Event, _ rs.Level, _ string) {
	e.Str(rs.TimestampFieldName, time.Now().Format(h.layout))
}

// isTerminal reports whether w is a character device, so colour goes to a
// person at a terminal and never into a file or a collector's pipe.
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	st, err := f.Stat()
	if err != nil {
		return false
	}
	return st.Mode()&os.ModeCharDevice != 0
}
