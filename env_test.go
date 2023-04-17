package goat

import (
	"strings"
	"testing"
	"time"

	"github.com/icrowley/fake"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_EnvString_Set(t *testing.T) {
	setupEnvTest()
	assertSet[string](t, EnvString, makeEnvKey(), "value", "default")
}

func Test_EnvString_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[string](t, EnvString, makeEnvKey(), "default")
}

func Test_EnvBool_Set(t *testing.T) {
	setupEnvTest()
	assertSet[bool](t, EnvBool, makeEnvKey(), false, true)
}

func Test_EnvBool_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[bool](t, EnvBool, makeEnvKey(), true)
}

func Test_EnvInt_Set(t *testing.T) {
	setupEnvTest()
	assertSet[int](t, EnvInt, makeEnvKey(), 1, 6)
}

func Test_EnvInt_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[int](t, EnvInt, makeEnvKey(), 6)
}

func Test_EnvInt32_Set(t *testing.T) {
	setupEnvTest()
	assertSet[int32](t, EnvInt32, makeEnvKey(), 1, 6)
}

func Test_EnvInt32_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[int32](t, EnvInt32, makeEnvKey(), 6)
}

func Test_EnvInt64_Set(t *testing.T) {
	setupEnvTest()
	assertSet[int64](t, EnvInt64, makeEnvKey(), 1, 6)
}

func Test_EnvInt64_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[int64](t, EnvInt64, makeEnvKey(), 6)
}

func Test_EnvUint16_Set(t *testing.T) {
	setupEnvTest()
	assertSet[uint16](t, EnvUint16, makeEnvKey(), 1, 6)
}

func Test_EnvUint16_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[uint16](t, EnvUint16, makeEnvKey(), 6)
}

func Test_EnvUint32_Set(t *testing.T) {
	setupEnvTest()
	assertSet[uint32](t, EnvUint32, makeEnvKey(), 1, 6)
}

func Test_EnvUint32_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[uint32](t, EnvUint32, makeEnvKey(), 6)
}

func Test_EnvUint64_Set(t *testing.T) {
	setupEnvTest()
	assertSet[uint64](t, EnvUint64, makeEnvKey(), 1, 6)
}

func Test_EnvUint64_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[uint64](t, EnvUint64, makeEnvKey(), 6)
}

func Test_EnvFloat64_Set(t *testing.T) {
	setupEnvTest()
	assertSet[float64](t, EnvFloat64, makeEnvKey(), 2.3, 7.8)
}

func Test_EnvFloat64_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[float64](t, EnvFloat64, makeEnvKey(), 7.8)
}

func Test_EnvTime_Set(t *testing.T) {
	setupEnvTest()
	assertSet[time.Time](t, EnvTime, makeEnvKey(), time.Now().AddDate(1, 0, 0), time.Now())
}

func Test_EnvTime_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[time.Time](t, EnvTime, makeEnvKey(), time.Now())
}

func Test_EnvDuration_Set(t *testing.T) {
	setupEnvTest()
	interval := time.Now().AddDate(1, 0, 0)
	assertSet[time.Duration](t, EnvDuration, makeEnvKey(), time.Since(time.Now().AddDate(-1, 0, 0)), time.Now().Sub(interval))
}

func Test_EnvDuration_NotSet(t *testing.T) {
	setupEnvTest()
	interval := time.Now().AddDate(1, 0, 0)
	assertNotSet[time.Duration](t, EnvDuration, makeEnvKey(), time.Now().Sub(interval))
}

func Test_EnvIntSlice_Set(t *testing.T) {
	setupEnvTest()
	assertSet[[]int](t, EnvIntSlice, makeEnvKey(), []int{1, 2, 3}, []int{5, 6, 7})
}

func Test_EnvIntSlice_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[[]int](t, EnvIntSlice, makeEnvKey(), []int{5, 6, 7})
}

func Test_EnvStringSlice_Set(t *testing.T) {
	setupEnvTest()
	assertSet[[]string](t, EnvStringSlice, makeEnvKey(), []string{"a", "b", "c"}, []string{"x", "y", "z"})
}

func Test_EnvStringSlice_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[[]string](t, EnvStringSlice, makeEnvKey(), []string{"x", "y", "z"})
}

func Test_EnvStringMap_Set(t *testing.T) {
	setupEnvTest()
	assertSet[map[string]any](t, EnvStringMap, makeEnvKey(), map[string]any{"a": 1, "b": "", "c": false}, map[string]any{"x": 0, "y": "qwerty", "z": []string{"asdf"}})
}

func Test_EnvStringMap_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[map[string]any](t, EnvStringMap, makeEnvKey(), map[string]any{"x": 0, "y": "qwerty", "z": []string{"asdf"}})
}

func Test_EnvStringMapString_Set(t *testing.T) {
	setupEnvTest()
	assertSet[map[string]string](t, EnvStringMapString, makeEnvKey(), map[string]string{"a": "1", "b": "2", "c": "3"}, map[string]string{"x": "4", "y": "5", "z": "6"})
}

func Test_EnvStringMapString_NotSet(t *testing.T) {
	setupEnvTest()
	assertNotSet[map[string]string](t, EnvStringMapString, makeEnvKey(), map[string]string{"x": "4", "y": "5", "z": "6"})
}

// Helpers

func setupEnvTest() {
	viper.Reset()
	viper.AutomaticEnv()
}

func makeEnvKey() string {
	return strings.ToUpper(fake.Characters())
}

type envReader[T any] func(key string, defaultValue T) T

func assertSet[T any](t *testing.T, subject envReader[T], key string, value, defaultValue T) {
	t.Helper()
	viper.Set(key, value)
	assert.Equal(t, value, subject(key, defaultValue))
}

func assertNotSet[T any](t *testing.T, subject envReader[T], key string, defaultValue T) {
	t.Helper()
	assert.Equal(t, defaultValue, subject(key, defaultValue))
}
