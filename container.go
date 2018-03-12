package goat

type Container struct {
	Utils  utilsInterface
	Config configInterface
	Path   pathInterface
}

func newContainer(p pathInterface, useConfig bool) *Container {
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
