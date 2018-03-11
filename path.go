package goat

import (
	"path/filepath"
	"os"
	"goat/types"
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

func initPath(u types.GoatUtilsInterface) (types.PathInterface, error) {
	if !haveRoot {
		rootPath = exeDir
	}
	path := types.NewPath(u, exePath, exeErr, rootPath)
	return path, exeErr
}

// Returns the path to the running executable.
// @TODO remove
func Root() string {
	mustHaveRoot()
	return rootPath
}

// Set the project root path manually, overriding the default, which is the path
// to the running executable.
func SetRoot(p string) {
	rootPath = p
	haveRoot = true
}

// Panics if the root path has not been set.
// @TODO remove
func mustHaveRoot() {
	if !haveRoot {
		panic("goat root not set")
	}
}
