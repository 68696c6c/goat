package goat

import (
	"os"
	"strconv"
)

func EnvString(s string, defaultValue string) string {
	val, ok := os.LookupEnv(s)
	if !ok {
		return defaultValue
	}
	return val
}

func MustEnvInt(s string, defaultValue int) int {
	val, ok := os.LookupEnv(s)
	if !ok {
		return defaultValue
	}
	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
}

func EnvBool(s string, defaultValue bool) bool {
	val, ok := os.LookupEnv(s)
	if !ok {
		return defaultValue
	}
	return val == "1"
}
