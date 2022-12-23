package sys

import (
	"fmt"
)

// Environment represents the environment in which a Goat app is running.
// For example, 'local' would be a developers local machine, 'test' would be the environment used while running unit
// tests, while 'dev', 'staging', and 'prod' are reserved for deployment environments.
//
// In 'local' and 'test' environments, Gin and GORM are run in debug mode.
type Environment string

const (
	typeNameEnvironment string      = "environment"
	EnvironmentLocal    Environment = "local"
	EnvironmentTest     Environment = "test"
	EnvironmentDev      Environment = "dev"
	EnvironmentStaging  Environment = "staging"
	EnvironmentProd     Environment = "prod"
)

func EnvironmentFromString(s string) (Environment, error) {
	if s == string(EnvironmentLocal) ||
		s == string(EnvironmentTest) ||
		s == string(EnvironmentDev) ||
		s == string(EnvironmentStaging) ||
		s == string(EnvironmentProd) {
		return Environment(s), nil
	}
	return Environment(""), fmt.Errorf("%s not a valid %s", s, typeNameEnvironment)
}

func (t Environment) String() string {
	return string(t)
}

// Gin has two log modes, expressed by string constants.
// Gin also has a third mode for testing, but that is reserved for its own internal use.
//
// GORM has two log modes, expressed by an integer: 2 enables detailed logs, 1 disables logging, and 0 prints only
// errors. GORM defaults to 0; calling gorm.LogMode(true) set the mode to 2; calling gorm.LogMode(false) will set the
// mode to 1.  Goat will only call gorm.LogMode to enable detailed logs and will never disable GORM logging completely.
func DebugFromEnvironment(env Environment) bool {
	switch env {
	case EnvironmentLocal:
		fallthrough
	case EnvironmentTest:
		return true
	default:
		return false
	}
}
