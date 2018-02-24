package goat

var initialized bool

func Init() []error {
	initPath()
	initConfig()
	errs := GetErrors()
	if len(errs) == 0 {
		initialized = true
	}
	return errs
}

func mustBeInitialized() {
	if !initialized {
		panic("goat is not initialized! Call goat.Init() before calling any other goat functions.")
	}
}
