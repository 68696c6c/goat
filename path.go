package goat

import (
	"path/filepath"
	"os"
	"goat/types"
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

func initPath(u types.GoatUtilsInterface) (types.PathInterface, error) {
	if !haveRoot {
		rootPath = exeDir
	}
	path := types.NewPath(u, exePath, exeErr, rootPath, runtime.Caller)
	return path, exeErr
}

// Set the project root path manually, overriding the default, which is the path
// to the running executable.
func SetRoot(p string) {
	rootPath = p
	haveRoot = true
}
