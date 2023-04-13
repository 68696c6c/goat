package sys

import (
	"github.com/pkg/errors"

	"github.com/68696c6c/goat/sys/database"
	"github.com/68696c6c/goat/sys/http"
	"github.com/68696c6c/goat/sys/log"
)

type Goat struct {
	DB   database.Service
	HTTP http.Service
	Log  log.Service
}

func Init(config Config) (Goat, error) {
	logService, err := log.NewService(config.Log)
	if err != nil {
		return Goat{}, errors.Wrapf(err, "failed to initialize log service")
	}
	return Goat{
		DB:   database.NewService(config.DB, logService),
		HTTP: http.NewService(config.HTTP, logService),
		Log:  logService,
	}, nil
}

type Config struct {
	DB   database.Config
	HTTP http.Config
	Log  log.Config
}
