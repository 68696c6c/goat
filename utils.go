package goat

import "github.com/pkg/errors"

// ValueOrDefault returns the value of input or defaultValue if input is nil.
func ValueOrDefault[T any](input *T, defaultValue T) T {
	if input == nil {
		return defaultValue
	}
	return *input
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
