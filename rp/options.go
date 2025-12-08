package rp

import "github.com/twmb/franz-go/pkg/kgo"

type (
	APMConfig struct {
		TxName string
		TxType string
	}

	Option interface {
		Apply(*Transport)
	}

	ConfigureFunc func(*Transport)

	LogLevel kgo.LogLevel
)

const (
	// LogLevelNone disables logging.
	LogLevelNone LogLevel = iota
	// LogLevelError logs all errors. Generally, these should not happen.
	LogLevelError
	// LogLevelWarn logs all warnings, such as request failures.
	LogLevelWarn
	// LogLevelInfo logs informational messages, such as requests. This is
	// usually the default log level.
	LogLevelInfo
	// LogLevelDebug logs verbose information, and is usually not used in
	// production.
	LogLevelDebug
)

func (c ConfigureFunc) Apply(rp *Transport) {
	c(rp)
}

func WithAPM(cfg *APMConfig) Option {
	return ConfigureFunc(func(t *Transport) {
		t.apm = cfg
	})
}

func WithPublishSyncEnabled() Option {
	return ConfigureFunc(func(t *Transport) {
		t.cfg.producerPublishSync = true
	})
}

func WithOnPublishAsync(fn func(Msg, error)) Option {
	return ConfigureFunc(func(t *Transport) {
		t.cfg.producerOnPublishAsync = fn
	})
}

func WithInternalLogger() Option {
	return ConfigureFunc(func(t *Transport) {
		t.cfg.WithInternalLogger = true
	})
}

func WithInternalLogLevel(level LogLevel) Option {
	return ConfigureFunc(func(t *Transport) {
		t.cfg.InternalLogLevel = level
	})
}
