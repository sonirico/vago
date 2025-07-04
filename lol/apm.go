package lol

type (
	APMConfig struct {
		Enabled bool
		URL     string
	}
)

var NoAPM = APMConfig{}
