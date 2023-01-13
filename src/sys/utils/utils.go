package utils

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func ValueOrDefault[T any](input *T, defaultValue T) T {
	if input == nil {
		return defaultValue
	}
	return *input
}

func EnvOrDefault[T any](envKey string, defaultValue T) T {
	if viper.IsSet(envKey) {
		result, ok := viper.Get(envKey).(T)
		if ok {
			return result
		}
	}
	return defaultValue
}

func MakeValidationError(field, rule string) error {
	return errors.Errorf("%s failed on validation '%s'", field, rule)
}
