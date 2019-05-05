package sys

import "fmt"

type Environment string

const (
	typeNameEnvironment string      = "environment"
	EnvironmentDev      Environment = "dev"
	EnvironmentTest     Environment = "test"
	EnvironmentStaging  Environment = "staging"
	EnvironmentProd     Environment = "prod"
)

func EnvironmentFromString(s string) (Environment, error) {
	if s == string(EnvironmentDev) ||
		s == string(EnvironmentTest) ||
		s == string(EnvironmentStaging) ||
		s == string(EnvironmentProd) {
		return Environment(s), nil
	}
	return Environment(""), fmt.Errorf("%s not a valid %s", s, typeNameEnvironment)
}

func (t Environment) String() string {
	return string(t)
}
