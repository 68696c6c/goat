package goat

var (
	initialized bool
	container   *Container
)

func Init() []error {
	p, err := initPath()
	panicIfError(err)

	container = newContainer(p, readConfig)
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

/**
 * Alias functions
 * Call underlying type functions.
 * @TODO is this worth it/the best way?
 */

func Root() string {
	return container.Path.Root()
}

func RootPath(path string) string {
	return container.Path.RootPath(path)
}

func ExePath() string {
	return container.Path.ExePath()
}

func ExeDir() string {
	container.Utils.MustBeInitialized()
	return container.Path.ExeDir()
}

func CWD() string {
	return container.Path.CWD()
}

func ConfigFileName() string {
	return container.Config.FileName()
}

func ConfigFilePath() string {
	return container.Config.FilePath()
}
