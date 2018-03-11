package types

import (
	"path/filepath"
)

type Path struct {
	callerFunc func(int) (uintptr, string, int, bool)
	utils      GoatUtilsInterface
	haveRoot   bool
	rootPath   string
	exePath    string
	exeDir     string
}

type PathInterface interface {
	Root() string
	RootPath(string) string
	ExePath() string
	ExeDir() string
	CWD() string
}

func NewPath(exePath string, exeError error, rootPath string, cf func(int) (uintptr, string, int, bool)) *Path {
	exeDir := filepath.Dir(exePath)
	if rootPath == "" {
		rootPath = exeDir
	}
	return &Path{
		callerFunc: cf,
		haveRoot:   exeError == nil,
		rootPath:   rootPath,
		exePath:    exePath,
		exeDir:     exeDir,
	}
}

func (p *Path) mustHaveRoot() {
	if !p.haveRoot {
		panic("goat root not set")
	}
}

// Returns the root path.
func (p *Path) Root() string {
	p.mustHaveRoot()
	return p.rootPath
}

// Returns the provided path appended to the root path.
func (p *Path) RootPath(path string) string {
	return p.Root() + "/" + path
}

// Returns the path to the running executable.
func (p *Path) ExePath() string {
	return p.exePath
}

// Returns the path to the dir holding the running executable.
// Will panic if goat has not been initialized.
func (p *Path) ExeDir() string {
	return p.exeDir
}

// Returns the path to the current working directory.
func (p *Path) CWD() string {
	_, b, _, ok := p.callerFunc(1)
	if !ok {
		// @TODO need an errors interface
		//addError("failed to get current directory")
		return ""
	}
	return filepath.Dir(b)
}
