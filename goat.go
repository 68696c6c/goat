package goat

import "goat/types"

var (
	initialized bool
	container   *Container
)

func Init() []error {
	u := types.NewGoatUtils()
	p, err := initPath(u)
	panicIfError(err)

	container = newContainer(u, p, readConfig)
	errs := GetErrors()
	if len(errs) == 0 {
		container.Utils.SetInitialized(true)
		return errs
	}
	errString := ErrorsToString(errs)
	panic("failed to initialize goat: " + errString)
}

func mustBeInitialized() {
	if !initialized {
		panic("goat is not initialized! Call goat.Init() before calling this function.")
	}
}

func panicIfError(err error) {
	if err != nil {
		panic("failed to initialize container: " + err.Error())
	}
}
