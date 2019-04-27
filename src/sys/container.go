package sys

import "github.com/68696c6c/goat/src/logging"

type Container struct {
	Config        configInterface
	Path          pathInterface
	LoggerService logging.LoggerService
}

func NewContainer(p pathInterface, useConfig bool) *Container {
	c := &Container{
		Utils: newUtils(),
		Path:  p,
	}

	// Config
	if useConfig {
		config, err := initConfig(p)
		panicIfError(err)
		c.Config = config
	}

	return c
}
