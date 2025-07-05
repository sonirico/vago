package lol

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"go.elastic.co/apm/module/apmzerolog/v2"
)

type zerologWrapper struct {
	log        zerolog.Logger
	apmEnabled bool
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
	if !z.apmEnabled {
		return z
	}
	hook := apmzerolog.TraceContextHook(ctx)

	return &zerologWrapper{
		log: z.log.With().Logger().Hook(hook),
	}
}

func zerologParseLevel(level Level) zerolog.Level {
	switch level {
	case LevelTrace:
		return zerolog.TraceLevel
	case LevelDebug:
		return zerolog.DebugLevel
	case LevelInfo:
		return zerolog.InfoLevel
	case LevelWarn:
		return zerolog.WarnLevel
	case LevelError:
		return zerolog.ErrorLevel
	case LevelFatal:
		return zerolog.FatalLevel
	case LevelPanic:
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

// NewZerolog creates a new Logger using zerolog as the backend
func NewZerolog(opts ...Opt) Logger {
	// Default configuration
	config := newDefaultConfig()
	// Apply options
	for _, opt := range opts {
		opt(&config)
	}

	// Parse log level
	logLevel := zerologParseLevel(config.level)

	// Configure zerolog
	zerolog.TimeFieldFormat = config.timeFieldFormat

	// Create base logger
	var logger zerolog.Logger
	if config.env == EnvLocal {
		// Pretty logging for local development
		logger = zerolog.New(zerolog.ConsoleWriter{Out: config.writer}).With().Timestamp().Logger()
	} else {
		// JSON logging for production
		logger = zerolog.New(config.writer).With().Timestamp().Logger()
	}

	// If APM is enabled, set the error stack marshaler and
	// create a multi-writer that includes the APM writer
	if config.apm {
		zerolog.ErrorStackMarshaler = apmzerolog.MarshalErrorStack
		logger = zerolog.New(zerolog.MultiLevelWriter(
			config.writer, new(apmzerolog.Writer)))
	}

	// Set log level
	logger = logger.Level(logLevel)

	log := &zerologWrapper{
		log: logger,
	}

	// Add initial fields
	return log.WithFields(config.fields)
}
