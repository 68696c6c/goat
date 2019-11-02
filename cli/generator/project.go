package generator

import (
	"github.com/pkg/errors"
	"go/build"
	"os"
	"path/filepath"
)

const packageSRC = "src"

type AuthorConfig struct {
	Name  string
	Email string
}

type Path struct {
	Abs string
	Rel string
}

type Packages struct {
	App    string
	CMD    string
	HTTP   string
	Models string
	Repos  string
}

type Paths struct {
	Root string
	Packages
}

type Imports struct {
	Packages
}

type ProjectConfig struct {
	Name       string
	Module     string
	ModuleName string
	License    string
	Author     AuthorConfig

	Paths   Paths
	Imports Imports

	Models []*Model
	Repos  []*Repo
}

func newPackages(base string) Packages {
	srcPath := joinPath(base, packageSRC)
	return Packages{
		App:    joinPath(srcPath, packageApp),
		CMD:    joinPath(srcPath, packageCMD),
		HTTP:   joinPath(srcPath, packageHandlers),
		Models: joinPath(srcPath, packageModels),
		Repos:  joinPath(srcPath, packageRepos),
	}
}

func (c *ProjectConfig) SetPaths() error {

	c.ModuleName = filepath.Base(c.Module)

	// Get the absolute path to the project (e.g. within $GOPATH)
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = build.Default.GOPATH
	}
	path := joinPath(goPath, "src")
	rootPath, err := filepath.Abs(path)
	if err != nil {
		return errors.Wrap(err, "failed determine absolute project path")
	}
	projectPath := joinPath(rootPath, c.Module)
	println("project path: ", projectPath)

	c.Paths = Paths{
		Root:     projectPath,
		Packages: newPackages(projectPath),
	}

	c.Imports = Imports{
		Packages: newPackages(c.Module),
	}

	return nil
}

func CreateProject(config *ProjectConfig) error {
	err := CreateDir(config.Paths.Root)
	if err != nil {
		return errors.Wrapf(err, "failed to create project directory '%s'", config.Paths.Root)
	}

	return nil
}
