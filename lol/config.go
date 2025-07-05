package lol

import (
	"io"

	"github.com/rs/zerolog"
)

// config holds the configuration for logger
type config struct {
	level           Level
	env             Env
	writer          io.Writer
	fields          Fields
	apm             bool
	timeFieldFormat string
}

func newDefaultConfig() config {
	return config{
		level:           LevelInfo,
		env:             EnvLocal,
		writer:          io.Discard,
		fields:          nil,
		apm:             false,
		timeFieldFormat: zerolog.TimeFormatUnix,
	}
}

type Opt func(*config)

// WithLevel sets the logging level
func WithLevel(level Level) Opt {
	return func(o *config) {
		o.level = level
	}
}

// WithApm enables APM tracing
func WithApm() Opt {
	return func(o *config) {
		o.apm = true
	}
}

// WithWriter sets the output writer for the logger. By default it uses os.Stderr
func WithWriter(w io.Writer) Opt {
	return func(o *config) {
		o.writer = w
	}
}

func WithEnv(env Env) Opt {
	return func(o *config) {
		o.env = env
	}
}

// WithFields sets initial fields
func WithFields(fields Fields) Opt {
	return func(c *config) {
		c.fields = fields
	}
}

// WithTimeFieldFormat sets the format for time fields
func WithTimeFieldFormat(format string) Opt {
	return func(c *config) {
		c.timeFieldFormat = format
	}
}
