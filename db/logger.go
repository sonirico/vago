package db

// Logger is the minimal logging contract this package needs from a
// caller-supplied logger: exactly the methods the executor and bulk-DML
// helpers call. It is a strict subset of github.com/sonirico/vago/lol.Logger's
// method set, so lol.Logger satisfies it with no adapter.
//
// WithField/WithTrace-style scoped loggers are deliberately not part of this
// interface. Go cannot satisfy a self-referential method (WithField(...) Logger)
// across two differently named interface types unless the return type is
// spelled identically at both ends - a shape-for-shape recursive interface in
// this package still fails to be satisfied by lol.Logger, whose own WithField
// returns lol.Logger, not db.Logger. Reproducing that chaining here would need
// an adapter, which is exactly the lol.Logger coupling this interface exists
// to drop from core. See the split report for the two call sites this
// affects (executor.go, executor_rwo.go), where a "mode"/trace tag that used
// to be attached via WithField/WithTrace is no longer available.
type Logger interface {
	Debug(args ...any)
	Debugf(format string, args ...any)
	Debugln(args ...any)
	Info(args ...any)
	Infoln(args ...any)
	Warnln(args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}
