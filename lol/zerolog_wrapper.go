package lol

import (
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog"
	"go.elastic.co/apm/v2"
)

type zerologWrapper struct {
	log zerolog.Logger
}

func (z *zerologWrapper) Trace(args ...any) {
	z.log.Trace().Msg(fmt.Sprint(args...))
}

func (z *zerologWrapper) Debug(args ...any) {
	z.log.Debug().Msg(fmt.Sprint(args...))
}

func (z *zerologWrapper) Print(args ...any) {
	z.log.Info().Msg(fmt.Sprint(args...))
}

func (z *zerologWrapper) Info(args ...any) {
	z.log.Info().Msg(fmt.Sprint(args...))
}

func (z *zerologWrapper) Warn(args ...any) {
	z.log.Warn().Msg(fmt.Sprint(args...))
}

func (z *zerologWrapper) Warning(args ...any) {
	z.log.Warn().Msg(fmt.Sprint(args...))
}

func (z *zerologWrapper) Error(args ...any) {
	z.log.Error().Msg(fmt.Sprint(args...))
}

func (z *zerologWrapper) Fatal(args ...any) {
	z.log.Fatal().Msg(fmt.Sprint(args...))
}

func (z *zerologWrapper) Panic(args ...any) {
	z.log.Panic().Msg(fmt.Sprint(args...))
}

// printf methods
func (z *zerologWrapper) Tracef(format string, args ...any) {
	z.log.Trace().Msgf(format, args...)
}

func (z *zerologWrapper) Debugf(format string, args ...any) {
	z.log.Debug().Msgf(format, args...)
}

func (z *zerologWrapper) Printf(format string, args ...any) {
	z.log.Info().Msgf(format, args...)
}

func (z *zerologWrapper) Infof(format string, args ...any) {
	z.log.Info().Msgf(format, args...)
}

func (z *zerologWrapper) Warnf(format string, args ...any) {
	z.log.Warn().Msgf(format, args...)
}

func (z *zerologWrapper) Warningf(format string, args ...any) {
	z.log.Warn().Msgf(format, args...)
}

func (z *zerologWrapper) Errorf(format string, args ...any) {
	z.log.Error().Msgf(format, args...)
}

func (z *zerologWrapper) Fatalf(format string, args ...any) {
	z.log.Fatal().Msgf(format, args...)
}

func (z *zerologWrapper) Panicf(format string, args ...any) {
	z.log.Panic().Msgf(format, args...)
}

// ln methods
func (z *zerologWrapper) Traceln(args ...any) {
	z.log.Trace().Msg(fmt.Sprintln(args...))
}

func (z *zerologWrapper) Debugln(args ...any) {
	z.log.Debug().Msg(fmt.Sprintln(args...))
}

func (z *zerologWrapper) Println(args ...any) {
	z.log.Info().Msg(fmt.Sprintln(args...))
}

func (z *zerologWrapper) Infoln(args ...any) {
	z.log.Info().Msg(fmt.Sprintln(args...))
}

func (z *zerologWrapper) Warnln(args ...any) {
	z.log.Warn().Msg(fmt.Sprintln(args...))
}

func (z *zerologWrapper) Warningln(args ...any) {
	z.log.Warn().Msg(fmt.Sprintln(args...))
}

func (z *zerologWrapper) Errorln(args ...any) {
	z.log.Error().Msg(fmt.Sprintln(args...))
}

func (z *zerologWrapper) Fatalln(args ...any) {
	z.log.Fatal().Msg(fmt.Sprintln(args...))
}

func (z *zerologWrapper) Panicln(args ...any) {
	z.log.Panic().Msg(fmt.Sprintln(args...))
}

func (z *zerologWrapper) WithField(key string, value any) Logger {
	return &zerologWrapper{
		log: z.log.With().Interface(key, value).Logger(),
	}
}

func (z *zerologWrapper) WithFields(fields Fields) Logger {
	logger := z.log.With()
	for k, v := range fields {
		logger = logger.Interface(k, v)
	}
	return &zerologWrapper{
		log: logger.Logger(),
	}
}

func (z *zerologWrapper) WithTrace(ctx context.Context) Logger {
	logger := z.log.With()

	// Add APM trace context if available
	if tx := apm.TransactionFromContext(ctx); tx != nil {
		logger = logger.Str("trace.id", tx.TraceContext().Trace.String())
		logger = logger.Str("transaction.id", tx.TraceContext().Transaction.String())
	}

	if span := apm.SpanFromContext(ctx); span != nil {
		logger = logger.Str("span.id", span.TraceContext().Span.String())
	}

	return &zerologWrapper{
		log: logger.Logger(),
	}
}

// NewZerologLogger creates a new Logger using zerolog as the backend
func NewZerologLogger(
	fields Fields,
	env, level string,
	writer io.Writer,
	apmConfig APMConfig,
) Logger {
	// Parse log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	// Configure zerolog
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z07:00"

	// Create base logger
	var logger zerolog.Logger
	if env == "local" || env == "development" {
		// Pretty logging for local development
		logger = zerolog.New(zerolog.ConsoleWriter{Out: writer}).With().Timestamp().Logger()
	} else {
		// JSON logging for production
		logger = zerolog.New(writer).With().Timestamp().Logger()
	}

	// Set log level
	logger = logger.Level(logLevel)

	// Add environment
	logger = logger.With().Str("env", env).Logger()

	// Add initial fields
	if fields != nil {
		loggerWith := logger.With()
		for k, v := range fields {
			loggerWith = loggerWith.Interface(k, v)
		}
		logger = loggerWith.Logger()
	}

	return &zerologWrapper{
		log: logger,
	}
}
