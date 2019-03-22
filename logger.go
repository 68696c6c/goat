package goat

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// @TODO not using this, so could you fucking not
var loggers = map[string]*logrus.Logger{}

// System Logger constructor. Returns a new logger instance writing to
// <binary path>/logs/sys.log, unless configured otherwise in the app config
// file under "logs.sys".
// Will panic if goat.Init() has not been called, since the app root path is
// needed to create the log file.
func NewLogger() *logrus.Logger {
	mustBeInitialized()
	return NewCustomLogger("sys")
}

// Return a new std out logger.
func NewSTDOutLogger() *logrus.Logger {
	return logrus.New()
}

// Returns an arbitrarily named logger instance. If not configured directly in
// the app config under "logs.<name>", goes to <binary path>/logs/<name>.log.
// Will panic if goat.Init() has not been called, since the app root path is
// needed to create the log file.
func NewCustomLogger(l string) *logrus.Logger {
	mustBeInitialized()
	if _, ok := loggers[l]; !ok {
		logger := logrus.New()

		conf := viper.GetString("logs." + l)

		if conf == "" {
			conf = fmt.Sprintf("%s/logs/%s.log", Root(), l)
		}

		fname, err := filepath.Abs(conf)
		if err != nil {
			// @TODO add a SetLoggerPanicMode(bool) function
			log.Printf("Error resolving path for log file %s: %s", conf, err)
			panic(err)
		}

		file, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
		if err != nil {
			panic(err)
		}

		logger.Out = file
		loggers[l] = logger
	}

	return loggers[l]
}
