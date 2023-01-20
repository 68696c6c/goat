package goat

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ValueOrDefault returns the value of input or defaultValue if input is nil.
func ValueOrDefault[T any](input *T, defaultValue T) T {
	if input == nil {
		return defaultValue
	}
	return *input
}

// EnvOrDefault returns the value of the specified env var or defaultValue if the env var is not set.
func EnvOrDefault[T any](envKey string, defaultValue T) T {
	if viper.IsSet(envKey) {
		// Goat sets viper.SetTypeByDefaultValue to true, so viper should have already performed the typecast.
		// Unfortunately, since viper.Get returns interface{}, we need to cast again in order for this function to return a
		// typed value, but in practice, this typecast should never fail.
		result, ok := viper.Get(envKey).(T)
		if ok {
			return result
		}
	}
	return defaultValue
}

// ValueToString attempts to cast the provided value to a string and returns the result.
func ValueToString(value any) (string, error) {
	s, ok := value.(string)
	if !ok {
		b, ok := value.([]byte)
		if !ok {
			return "", errors.Errorf("failed to parse value to string: %+v", value)
		}
		s = string(b)
	}
	return s, nil
}
