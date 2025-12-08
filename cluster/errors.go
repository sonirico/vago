package cluster

import "errors"

var (
	ErrCannotScan = errors.New("cannot scan")
	ErrNoSessions = errors.New("no sessions")
)
