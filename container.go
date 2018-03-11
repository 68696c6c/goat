package goat

import "goat/types"

type Container struct {
	Utils  types.GoatUtilsInterface
	Config types.ConfigInterface
	Path   types.PathInterface
}

func newContainer(p types.PathInterface, useConfig bool) *Container {
	c := &Container{
		Utils: types.NewGoatUtils(),
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
