package sys

import (
	"github.com/68696c6c/goat/src/database"
	"github.com/68696c6c/goat/src/logging"
)

type Container struct {
	DatabaseService database.Service
	LoggerService   logging.Service
}

func NewContainer(mode, mainDBName string, l logging.LoggerConfig) Container {
	c := &Container{
		DatabaseService: database.NewServiceGORM(mainDBName),
		LoggerService:   logging.NewServiceLogrus(l),
	}
	return c
}
