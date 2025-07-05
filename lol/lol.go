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
//	logger := lol.NewZerolog(
//		lol.WithFields(lol.Fields{"service": "myapp"}),
//		lol.WithEnv("production"),
//		lol.WithLevel("info"),
//		lol.WithWriter(os.Stdout),
//		lol.WithAPM(lol.APMConfig{Enabled: true}),
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
	ZeroLogger = NewZerolog(
		WithFields(Fields{"type": "default"}),
		WithEnv(EnvLocal),
		WithLevel(LevelInfo),
		WithWriter(os.Stderr),
	)

	ZeroTestLogger = NewZerolog(
		WithEnv(EnvTest),
		WithLevel(LevelError),
		WithWriter(os.Stderr),
	)

	ZeroDiscardLogger = NewZerolog(
		WithWriter(io.Discard),
	)
)
