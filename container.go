package goat

type Container struct {
	Utils utilsInterface
	Path  pathInterface
}

func newContainer(p pathInterface, useConfig bool) *Container {
	c := &Container{
		Utils: newUtils(),
		Path:  p,
	}
	return c
}
