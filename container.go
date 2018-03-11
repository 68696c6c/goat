package goat

import "goat/types"

type Container struct {
	Utils  types.GoatUtilsInterface
	Config types.ConfigInterface
	Path   types.PathInterface
}

func newContainer(u *types.GoatUtils, p types.PathInterface, useConfig bool) *Container {
	c := &Container{
		Utils: u,
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
