package logging

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/68696c6c/goat/utils"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

const (
	loggerPath  = "logs"
	loggerExt   = "log"
	loggerFile  = "custom"
	loggerLevel = logrus.InfoLevel
)

// Default logging is done to standard out.
// Custom loggers can be created using log files if needed, but Goat is
// optimized for using a single logger writing to standard out for ease of use
// with Docker and Amazon ECS.
type Service interface {
	NewLogger() *logrus.Logger
	NewFileLogger(path string) (*logrus.Logger, error)
	GetLogPath() string
	GetLogExt() string
	GetLogLevel() string
}

type Config struct {
	Path  string
	Ext   string
	Level string
}

type ServiceLogrus struct {
	logPath    string
	logFileExt string
	level      logrus.Level
}

func NewServiceLogrus(c Config) ServiceLogrus {
	// If an invalid log level is provided, fallback to the default.
	level, err := logrus.ParseLevel(c.Level)
	if err != nil {
		level = loggerLevel
	}
	return ServiceLogrus{
		logPath:    utils.ArgStringD(c.Path, loggerPath),
		logFileExt: utils.ArgStringD(c.Ext, loggerExt),
		level:      level,
	}
}

// Returns the configured log file path.
func (s ServiceLogrus) GetLogPath() string {
	return s.logPath
}

// Returns the configured log file extension.
func (s ServiceLogrus) GetLogExt() string {
	return s.logFileExt
}

// Returns the configured default log level.
func (s ServiceLogrus) GetLogLevel() string {
	return s.level.String()
}

// Returns a new standard out logger.
func (s ServiceLogrus) NewLogger() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(s.level)
	return l
}

// Returns a new logger that writes to the specified file, relative to the
// configured log path using the configured log file extension.
// If a file name is not specified a default name is used.
func (s ServiceLogrus) NewFileLogger(name string) (*logrus.Logger, error) {
	logger := logrus.New()

	if name == "" {
		name = loggerFile
	}

	path := fmt.Sprintf("%s/%s.%s", s.logPath, name, s.logFileExt)

	fname, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "error resolving log file path")
	}

	file, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return nil, errors.Wrap(err, "error opening log file")
	}

	logger.Out = file

	return logger, nil
}
