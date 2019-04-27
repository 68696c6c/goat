package logging

import (
	"fmt"
	"os"
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Path:  fake.Word(),
		Ext:   fake.Word(),
		Level: fake.Word(),
	}
}

func TestLogger_NewLoggerServiceLogrus_Config(t *testing.T) {
	c := getLoggerConfig()
	c.Level = "info"
	s := NewLoggerServiceLogrus(c)

	assert.Equal(t, c.Path, s.GetLogPath(), "unexpected log path")
	assert.Equal(t, c.Ext, s.GetLogExt(), "unexpected log file extension")
	assert.Equal(t, c.Level, s.GetLogLevel(), "unexpected log level")
}

func TestLogger_NewLoggerServiceLogrus_Defaults(t *testing.T) {
	s := NewLoggerServiceLogrus(LoggerConfig{})

	assert.Equal(t, loggerPath, s.GetLogPath(), "unexpected default log path")
	assert.Equal(t, loggerExt, s.GetLogExt(), "unexpected default log file extension")
	assert.Equal(t, loggerLevel.String(), s.GetLogLevel(), "unexpected default log level")
}

func TestLogger_LoggerServiceLogrus_NewLogger(t *testing.T) {
	c := getLoggerConfig()
	s := NewLoggerServiceLogrus(c)
	l := s.NewLogger()

	assert.Equal(t, os.Stderr, l.Out, "unexpected default logger out")
	assert.Equal(t, loggerLevel.String(), l.Level.String(), "unexpected default logger level")
}

func TestLogger_LoggerServiceLogrus_NewLogger_LogLevel(t *testing.T) {
	c := getLoggerConfig()
	c.Level = "error"
	s := NewLoggerServiceLogrus(c)
	l := s.NewLogger()

	assert.Equal(t, c.Level, l.Level.String(), "unexpected logger level")
}

func TestLogger_LoggerServiceLogrus_NewFileLogger_FileCreated(t *testing.T) {
	c := LoggerConfig{
		Path: "test",
		Ext:  "ext",
	}
	s := NewLoggerServiceLogrus(c)
	f := fake.Word()
	l, err := s.NewFileLogger(f)
	require.Nil(t, err, "unexpected error returned")
	require.NotNil(t, l, "nil logger returned")

	// Make sure the log file was created.
	p := fmt.Sprintf("%s/%s.%s", c.Path, f, c.Ext)
	i, err := os.Stat(p)
	assert.Nil(t, err, "failed to create log file")

	// Make sure the file is writable.
	assert.Equal(t, "-rw-r--r--", i.Mode().String(), "unexpected log file mode")

	// Remove the log file.
	err = os.Remove(p)
	require.NotNil(t, l, "failed to remove test log file")
}
