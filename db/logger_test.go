package db

import "github.com/sonirico/vago/lol"

// lol.Logger satisfies Logger with no adapter: Logger's method set is a
// strict, non-recursive subset of lol.Logger's (WithField/WithTrace
// deliberately excluded - see logger.go), so this assignment needs nothing
// beyond what both interfaces already declare.
var _ Logger = (lol.Logger)(nil)
