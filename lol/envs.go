package lol

type Env uint8

const (
	EnvTest Env = iota
	EnvLocal
	EnvDev
	EnvProd
)
