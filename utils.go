package goat

import (
	"github.com/pkg/errors"
)

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

func ValueInArray[T comparable](value T, items []T) bool {
	for _, t := range items {
		if t == value {
			return true
		}
	}
	return false
}

func Spread[T any](slices ...[]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

func Ref[T any](value T) *T {
	return &value
}

func DeRef[T any](value *T) T {
	if value == nil {
		var empty T
		return empty
	}
	return *value
}

type UniqueItems[T comparable] map[T]struct{}

func (u UniqueItems[T]) Add(key T) {
	u[key] = struct{}{}
}

func (u UniqueItems[T]) Slice() []T {
	var result []T
	for key := range u {
		result = append(result, key)
	}
	return result
}
