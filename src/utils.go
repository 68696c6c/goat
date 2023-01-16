package goat

import "github.com/spf13/viper"

// ValueOrDefault returns the value of input or defaultValue if input is nil.
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
