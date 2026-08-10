package testit

// Logger is the minimal logging contract this package needs from a
// caller-supplied logger: exactly the methods the resource pool and
// RunSafe call. It is a strict subset of
// github.com/sonirico/vago/lol.Logger's method set, so lol.Logger
// satisfies it with no adapter.
//
// WithField-style scoped loggers are deliberately not part of this
// interface: Go cannot satisfy a self-referential method
// (WithField(...) Logger) across two differently named interface types,
// so reproducing that chaining here would force the lol coupling this
// interface exists to drop from core. Backend modules that want scoped
// logging keep their own lol dependency instead.
type Logger interface {
	Info(args ...any)
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}
