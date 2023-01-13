package logging

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/68696c6c/goat/sys/utils"
)

const (
	loggerPath  = "logs"
	loggerExt   = "log"
	loggerFile  = "custom.log"
	loggerLevel = logrus.InfoLevel
)

// TODO: replace with zap???

// Default logging is done to standard out.
// Custom loggers can be created using log files if needed, but Goat is
// optimized for using a single logger writing to standard out for ease of use
// with Docker and Amazon ECS.

type Service interface {
	NewLogger() *logrus.Logger
	NewFileLogger(fileName ...string) (*logrus.Logger, error)
	GetLogPath() string
	GetLogExt() string
	GetLogLevel() string
}

type Config struct {
	Path  string
	Ext   string
	Level string
}

type service struct {
	logPath    string
	logFileExt string
	level      logrus.Level
}

func NewService(c Config) Service {
	// If an invalid log level is provided, fallback to the default.
	level, err := logrus.ParseLevel(c.Level)
	if err != nil {
		level = loggerLevel
	}
	return service{
		logPath:    utils.ArgStringD(c.Path, loggerPath),
		logFileExt: utils.ArgStringD(c.Ext, loggerExt),
		level:      level,
	}
}

// GetLogPath returns the configured log file path.
func (s service) GetLogPath() string {
	return s.logPath
}

// GetLogExt returns the configured log file extension.
func (s service) GetLogExt() string {
	return s.logFileExt
}

// GetLogLevel returns the configured default log level.
func (s service) GetLogLevel() string {
	return s.level.String()
}

// NewLogger returns a new logger instance that writes to standard out.
func (s service) NewLogger() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(s.level)
	return l
}

// NewFileLogger returns a new logger that writes to the specified file, relative to
// the configured log path.If a file name is not specified a default name is used.
func (s service) NewFileLogger(fileName ...string) (*logrus.Logger, error) {
	logger := logrus.New()

	name := loggerFile
	if len(fileName) > 0 && fileName[0] != "" {
		name = fileName[0]
	}

	path := fmt.Sprintf("%s/%s", s.logPath, name)

	logPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "error resolving log file path")
	}

	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return nil, errors.Wrap(err, "error opening log file")
	}

	logger.Out = file

	return logger, nil
}
