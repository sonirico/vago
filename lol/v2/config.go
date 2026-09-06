package lol

import "io"

// Config is what every backend needs to build a Logger. Build one with
// NewConfig and the With* options, or as a literal.
type Config struct {
	Level  Level
	Env    Env
	Writer io.Writer
	Fields Fields
	// TimeFormat is the layout for the timestamp field; empty means the
	// backend's default.
	TimeFormat string
}

// Opt mutates a Config.
type Opt func(*Config)

// NewConfig returns the default Config (info level, EnvLocal, io.Discard)
// with opts applied.
func NewConfig(opts ...Opt) Config {
	c := Config{Level: LevelInfo, Env: EnvLocal, Writer: io.Discard}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

func WithLevel(level Level) Opt {
	return func(c *Config) { c.Level = level }
}

func WithEnv(env Env) Opt {
	return func(c *Config) { c.Env = env }
}

func WithWriter(w io.Writer) Opt {
	return func(c *Config) { c.Writer = w }
}

func WithFields(fields Fields) Opt {
	return func(c *Config) { c.Fields = fields }
}

func WithTimeFormat(layout string) Opt {
	return func(c *Config) { c.TimeFormat = layout }
}
