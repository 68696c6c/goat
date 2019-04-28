package sys

import (
	"github.com/68696c6c/goat/src/database"
	"github.com/68696c6c/goat/src/http"
	"github.com/68696c6c/goat/src/logging"
)

type Container struct {
	DatabaseService database.Service
	HTTPService     http.Service
	LoggerService   logging.Service
}

func NewContainer(mainDBName string, h http.Config, l logging.LoggerConfig) Container {
	c := Container{
		DatabaseService: database.NewServiceGORM(mainDBName),
		HTTPService:     http.NewServiceGin(h),
		LoggerService:   logging.NewServiceLogrus(l),
	}
	return c
}
