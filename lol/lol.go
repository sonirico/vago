// Package lol (lots of logs) provides a unified logging interface with multiple backends.
//
// This package offers a simple, structured logging interface that can be backed by
// different logging implementations. Currently it supports zerolog as the primary backend.
//
// Key features:
// - Structured logging with fields
// - Multiple log levels (trace, debug, info, warn, error, fatal, panic)
// - APM trace context integration
// - Environment-aware configuration
// - Testing utilities
//
// Basic usage:
//
//	logger := lol.NewZerologLogger(
//		lol.Fields{"service": "myapp"},
//		"production",
//		"info",
//		os.Stdout,
//		lol.APMConfig{Enabled: true},
//	)
//
//	logger.Info("Application started")
//	logger.WithField("user_id", 123).Warn("User action")
//
// For testing:
//
//	testLogger := lol.NewTest()
//	testLogger.Error("This won't be printed")
package lol

import (
	"io"
	"os"
)

var (
	Log = NewZerologLogger(
		Fields{"type": "default"},
		"local",
		"info",
		os.Stderr,
		APMConfig{Enabled: true},
	)

	TestLog = NewZerologLogger(
		nil,
		"test",
		"error",
		os.Stderr,
		NoAPM,
	)
)

func NewTest() Logger {
	logLevel := "debug"
	apmCfg := APMConfig{Enabled: false}
	logger := NewZerologLogger(Fields{}, "test", logLevel, io.Discard, apmCfg)

	return logger
}

type option struct {
	level   string
	env     string
	apm     bool
	writeTo io.Writer
}

type Opt func(*option)

func WithLevel(level string) Opt {
	return func(o *option) {
		o.level = level
	}
}

func WithApm() Opt {
	return func(o *option) {
		o.apm = true
	}
}

func WithWriter(w io.Writer) Opt {
	return func(o *option) {
		o.writeTo = w
	}
}

func WithEnv(env string) Opt {
	return func(o *option) {
		o.env = env
	}
}

func NewTestOpts(opts ...Opt) Logger {
	o := &option{
		level:   "debug",
		env:     "test",
		apm:     false,
		writeTo: io.Discard,
	}

	for _, opt := range opts {
		opt(o)
	}

	logLevel := o.level
	logger := NewZerologLogger(Fields{}, o.env, logLevel, o.writeTo, APMConfig{Enabled: o.apm})

	return logger
}
