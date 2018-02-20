package goat

import (
	"path/filepath"
	"os"
)

var (
	exePath, exeErr = executableClean()
	exeDir          string
	initialized     bool
)

func executableClean() (string, error) {
	p, err := os.Executable()
	return filepath.Clean(p), err
}

func Init() (bool, error) {
	exeDir = filepath.Dir(exePath)
	initialized = exeErr == nil
	return exeErr == nil, exeErr
}

func GetProjectRoot() string {
	panicIfNotInitialized()
	return exeDir
}

func panicIfNotInitialized() {
	if !initialized {
		panic("goat is not initialized")
	}
}
