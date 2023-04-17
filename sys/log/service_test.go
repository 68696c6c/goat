package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test_NewService_InvalidLevel(t *testing.T) {
	_, err := NewService(Config{
		Level:      zap.AtomicLevel{},
		Stacktrace: false,
	})
	assert.NotNil(t, err)
}

func Test_Logger(t *testing.T) {
	level := zap.DebugLevel
	subject := setupLogService(t, level)
	result := subject.Logger().Desugar().Check(level, "")
	assert.NotNil(t, result)
}

func Test_StrictLogger(t *testing.T) {
	level := zap.DebugLevel
	subject := setupLogService(t, level)
	result := subject.StrictLogger().Check(level, "")
	assert.NotNil(t, result)
}

func Test_GinLogger(t *testing.T) {
	subject := setupLogService(t, zap.DebugLevel)
	result := subject.GinLogger()
	assert.NotNil(t, result)
}

func Test_GinRecovery(t *testing.T) {
	subject := setupLogService(t, zap.DebugLevel)
	result := subject.GinRecovery()
	assert.NotNil(t, result)
}

func Test_GormLogger(t *testing.T) {
	subject := setupLogService(t, zap.DebugLevel)
	result := subject.GormLogger()
	assert.NotNil(t, result)
}

func setupLogService(t *testing.T, level zapcore.Level) Service {
	result, err := NewService(Config{
		Level:      zap.NewAtomicLevelAt(level),
		Stacktrace: false,
	})
	require.Nil(t, err, "failed to initialize logger service")
	return result
}
