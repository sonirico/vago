package testit

type noopLogger struct{}

// NewNoopLogger returns a Logger that discards everything, for callers
// that do not care about resource lifecycle logs.
func NewNoopLogger() Logger {
	return noopLogger{}
}

func (noopLogger) Info(args ...any)                  {}
func (noopLogger) Infof(format string, args ...any)  {}
func (noopLogger) Errorf(format string, args ...any) {}
