package lol

import "context"

type Fields map[string]any

type Logger interface {
	Trace(args ...any)
	Debug(args ...any)
	Print(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Warning(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Panic(args ...any)
	// printf
	Tracef(format string, args ...any)
	Debugf(format string, args ...any)
	Printf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Warningf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Panicf(format string, args ...any)
	// ln
	Traceln(args ...any)
	Debugln(args ...any)
	Println(args ...any)
	Infoln(args ...any)
	Warnln(args ...any)
	Warningln(args ...any)
	Errorln(args ...any)
	Fatalln(args ...any)
	Panicln(args ...any)

	WithField(key string, value any) Logger
	WithFields(fields Fields) Logger

	WithTrace(ctx context.Context) Logger
}
