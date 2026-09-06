package lol

// Env is where the process runs. Backends may choose their output format
// from it: EnvLocal is for a human at a terminal, the rest for a collector.
type Env uint8

const (
	EnvTest Env = iota
	EnvLocal
	EnvDev
	EnvProd
)
