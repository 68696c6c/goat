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
	rootPath        string
)

func executableClean() (string, error) {
	p, err := os.Executable()
	return filepath.Clean(p), err
}

// Set the path to the running executable.
func initPath() (bool, error) {
	exeDir = filepath.Dir(exePath)
	haveRoot = exeErr == nil
	rootPath = exeDir
	return exeErr == nil, exeErr
}

// Returns the path to the running executable.
func Root() string {
	mustHaveRoot()
	return rootPath
}

func RootPath(path string) string {
	return Root() + "/" + path
}

// Set the project root path manually, overriding the default, which is the path
// to the running executable.
func SetRoot(p string) {
	rootPath = p
	haveRoot = true
}

// Returns the path to the running executable.
// Will panic if goat has not been initialized.
func ExePath() string {
	return exePath
}

// Returns the path to the dir holding the running executable.
// Will panic if goat has not been initialized.
func ExeDir() string {
	mustBeInitialized()
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
		panic("goat root not set")
	}
}
