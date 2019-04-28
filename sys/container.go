package sys

import (
	"github.com/68696c6c/goat/database"
	"github.com/68696c6c/goat/http"
	"github.com/68696c6c/goat/logging"
)

type Container struct {
	DatabaseService database.Service
	HTTPService     http.Service
	LoggingService  logging.Service
}

func NewContainer(mainDBName string, h http.Config, l logging.LoggerConfig) Container {
	c := Container{
		DatabaseService: database.NewServiceGORM(mainDBName),
		HTTPService:     http.NewServiceGin(h),
		LoggingService:  logging.NewServiceLogrus(l),
	}
	return c
}
