package lol

import "fmt"

type Environment uint8

const (
	TestEnvironment Environment = iota
	LocalEnvironment
	DevelopmentEnvironment
	StagingEnvironment
	ProductionEnvironment
)

func ParseEnv(env string) Environment {
	switch env {
	case "test":
		return TestEnvironment
	case "local":
		return LocalEnvironment
	case "development", "dev":
		return DevelopmentEnvironment
	case "staging":
		return StagingEnvironment
	case "production", "prod":
		return ProductionEnvironment
	default:
		panic(fmt.Sprintf("invalid environment: %s", env))
	}
}
