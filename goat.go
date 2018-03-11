package goat

var initialized bool

func Init() []error {
	initPath()
	if readConfig {
		initConfig()
	}
	errs := GetErrors()
	if len(errs) == 0 {
		initialized = true
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
