package filesystem

import (
	"path/filepath"
)

type path struct {
	callerFunc func(int) (uintptr, string, int, bool)
	haveRoot   bool
	rootPath   string
	exePath    string
	exeDir     string
}

type pathInterface interface {
	Root() string
	RootPath(string) string
	ExePath() string
	ExeDir() string
	CWD() string
}

func newPath(exePath string, exeError error, rootPath string, cf func(int) (uintptr, string, int, bool)) *path {
	exeDir := filepath.Dir(exePath)
	if rootPath == "" {
		rootPath = exeDir
	}
	return &path{
		callerFunc: cf,
		haveRoot:   exeError == nil,
		rootPath:   rootPath,
		exePath:    exePath,
		exeDir:     exeDir,
	}
}

func (p *path) mustHaveRoot() {
	if !p.haveRoot {
		panic("goat root not set")
	}
}

// Returns the root path.
func (p *path) Root() string {
	p.mustHaveRoot()
	return p.rootPath
}

// Returns the provided path appended to the root path.
func (p *path) RootPath(path string) string {
	return p.Root() + "/" + path
}

// Returns the path to the running executable.
func (p *path) ExePath() string {
	return p.exePath
}

// Returns the path to the dir holding the running executable.
// Will panic if goat has not been initialized.
func (p *path) ExeDir() string {
	return p.exeDir
}

// Returns the path to the current working directory.
func (p *path) CWD() string {
	_, b, _, ok := p.callerFunc(2)
	if !ok {
		// @TODO need an errors interface
		//addError("failed to get current directory")
		return ""
	}
	return filepath.Dir(b)
}
