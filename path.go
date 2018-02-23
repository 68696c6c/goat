package goat

import (
	"path/filepath"
	"os"
	"runtime"
)

var (
	exePath, exeErr = executableClean()
	exeDir          string
	haveRoot        bool
)

func executableClean() (string, error) {
	p, err := os.Executable()
	return filepath.Clean(p), err
}

// Set the path to the running executable.
func initPath() (bool, error) {
	exeDir = filepath.Dir(exePath)
	haveRoot = exeErr == nil
	return exeErr == nil, exeErr
}

// Returns the path to the running executable.
func Root() string {
	mustHaveRoot()
	return exeDir
}

// Returns the path to the current working directory.
func CWD() string {
	_, b, _, ok := runtime.Caller(1)
	if !ok {
		addError("failed to get current directory")
	}
	return filepath.Dir(b)
}

func mustHaveRoot() {
	if !haveRoot {
		addError("goat root not set")
	}
}
