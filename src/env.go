package goat

import (
	"time"

	"github.com/spf13/viper"
)

// Viper conveniently supports reading configuration from a variety of sources.  For our env helpers to work seamlessly
// with Viper, we need to use Viper to read env vars rather than os.LookupEnv.  While it might initially seem possible
// to use a generic function for this purpose, unfortunately that does not seem to be possible.  See the examples below
// for an explanation.
//
// This type assertion doesn't work because viper.Get returns interface{} so value is never a T and ok is never true.
// 		func EnvOrDefault[T any](key string, defaultValue T) T {
// 			if viper.IsSet(key) {
// 				value := viper.Get(key)
// 				result, ok := value.(T)
// 				if ok {
// 					return result
// 				}
// 			}
// 			return defaultValue
// 		}
//
// This doesn't work because it isn't possible to type cast using a type parameter without first passing a type
// assertion, which won't work because viper.Get returns an interface{} (see above).
// 		func EnvOrDefault[T any](key string, defaultValue T) T {
// 			if viper.IsSet(key) {
// 				value := viper.Get(key)
// 				return T(value)
// 			}
// 			return defaultValue
// 		}
//
// There are several problems trying to use reflection.
// 		func EnvOrDefault[T any](key string, defaultValue T) T {
// 			if viper.IsSet(key) {
// 				switch reflect.TypeOf(defaultValue).String() {
// 				 case "int":
// 				 	// This doesn't work because viper.GetInt returns an int, but we need to return a T.
// 				 	return viper.GetInt(key)
// 				 }
//
// 				 case "int":
// 				 	value := viper.GetInt(key)
// 				 	// This doesn't work because value is not an interface type
// 				 	result, ok := value.(T)
// 				 	if ok {
// 				 		return result
// 				 	}
//
// 				 case "int":
// 				 	value := viper.GetInt(key)
// 				 	// This doesn't work because Go doesn't support typecasting using a type parameter.
// 				 	return T(value)
// 				 }
// 			}
// 			return defaultValue
// 		}
//
// Using a type switch also doesn't work because of all the reasons in the previous example plus the fact that Go
// doesn't support type switches using type parameters.
// 		func EnvOrDefault[T any](key string, defaultValue T) T {
// 			if viper.IsSet(key) {
// 				switch defaultValue.(type) {
// 				case int:
// 					return viper.GetInt(key)
// 				}
// 			}
// 			return defaultValue
// 		}

// EnvString returns the value of the specified env var or defaultValue if the env var is not set.
func EnvString(key string, defaultValue string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}
	return defaultValue
}

// EnvBool returns the value of the specified env var or defaultValue if the env var is not set.
func EnvBool(key string, defaultValue bool) bool {
	if viper.IsSet(key) {
		return viper.GetBool(key)
	}
	return defaultValue
}

// EnvInt returns the value of the specified env var or defaultValue if the env var is not set.
func EnvInt(key string, defaultValue int) int {
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}
	return defaultValue
}

// EnvInt32 returns the value of the specified env var or defaultValue if the env var is not set.
func EnvInt32(key string, defaultValue int32) int32 {
	if viper.IsSet(key) {
		return viper.GetInt32(key)
	}
	return defaultValue
}

// EnvInt64 returns the value of the specified env var or defaultValue if the env var is not set.
func EnvInt64(key string, defaultValue int64) int64 {
	if viper.IsSet(key) {
		return viper.GetInt64(key)
	}
	return defaultValue
}

// EnvUint16 returns the value of the specified env var or defaultValue if the env var is not set.
func EnvUint16(key string, defaultValue uint16) uint16 {
	if viper.IsSet(key) {
		return viper.GetUint16(key)
	}
	return defaultValue
}

// EnvUint32 returns the value of the specified env var or defaultValue if the env var is not set.
func EnvUint32(key string, defaultValue uint32) uint32 {
	if viper.IsSet(key) {
		return viper.GetUint32(key)
	}
	return defaultValue
}

// EnvUint64 returns the value of the specified env var or defaultValue if the env var is not set.
func EnvUint64(key string, defaultValue uint64) uint64 {
	if viper.IsSet(key) {
		return viper.GetUint64(key)
	}
	return defaultValue
}

// EnvFloat64 returns the value of the specified env var or defaultValue if the env var is not set.
func EnvFloat64(key string, defaultValue float64) float64 {
	if viper.IsSet(key) {
		return viper.GetFloat64(key)
	}
	return defaultValue
}

// EnvTime returns the value of the specified env var or defaultValue if the env var is not set.
func EnvTime(key string, defaultValue time.Time) time.Time {
	if viper.IsSet(key) {
		return viper.GetTime(key)
	}
	return defaultValue
}

// EnvDuration returns the value of the specified env var or defaultValue if the env var is not set.
func EnvDuration(key string, defaultValue time.Duration) time.Duration {
	if viper.IsSet(key) {
		return viper.GetDuration(key)
	}
	return defaultValue
}

// EnvIntSlice returns the value of the specified env var or defaultValue if the env var is not set.
func EnvIntSlice(key string, defaultValue []int) []int {
	if viper.IsSet(key) {
		return viper.GetIntSlice(key)
	}
	return defaultValue
}

// EnvStringSlice returns the value of the specified env var or defaultValue if the env var is not set.
func EnvStringSlice(key string, defaultValue []string) []string {
	if viper.IsSet(key) {
		return viper.GetStringSlice(key)
	}
	return defaultValue
}

// EnvStringMap returns the value of the specified env var or defaultValue if the env var is not set.
func EnvStringMap(key string, defaultValue map[string]any) map[string]any {
	if viper.IsSet(key) {
		return viper.GetStringMap(key)
	}
	return defaultValue
}

// EnvStringMapString returns the value of the specified env var or defaultValue if the env var is not set.
func EnvStringMapString(key string, defaultValue map[string]string) map[string]string {
	if viper.IsSet(key) {
		return viper.GetStringMapString(key)
	}
	return defaultValue
}
