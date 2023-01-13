package goat

import "github.com/68696c6c/goat/sys/utils"

func ValueOrDefault[T any](input *T, defaultValue T) T {
	return utils.ValueOrDefault[T](input, defaultValue)
}

func EnvOrDefault[T any](envKey string, defaultValue T) T {
	return utils.EnvOrDefault[T](envKey, defaultValue)
}
