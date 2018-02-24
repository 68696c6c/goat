package goat

var initialized bool

func Init() {
	initPath()
	initConfig()
	initialized = true
}

func mustBeInitialized() {
	if !initialized {
		panic("goat is not initialized! Call goat.Init() before calling any other goat functions.")
	}
}
