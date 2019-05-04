package generator

import "github.com/pkg/errors"

type AuthorConfig struct {
	Name  string
	Email string
}

type PathsConfig struct {
	App      string
	CMD      string
	Handlers string
	Models   string
	Repos    string
}

type ProjectConfig struct {
	Name    string
	License string
	Author  AuthorConfig

	Module   string
	RootPath string
	DirName  string
	SRCPath  string

	Models []*Model
	Repos  []*Repo

	Paths   PathsConfig
	Imports PathsConfig
}

func (c *ProjectConfig) SetPaths() {
	srcPathRel := "/src"
	appPathRel := srcPathRel + "/app"
	cmdPathRel := srcPathRel + "/cmd"
	handlersPathRel := srcPathRel + "/handlers"
	modelsPathRel := srcPathRel + "/models"
	reposPathRel := srcPathRel + "/repos"

	c.SRCPath = c.RootPath + srcPathRel

	c.Paths = PathsConfig{
		App:      c.RootPath + appPathRel,
		CMD:      c.RootPath + cmdPathRel,
		Handlers: c.RootPath + handlersPathRel,
		Models:   c.RootPath + modelsPathRel,
		Repos:    c.RootPath + reposPathRel,
	}

	c.Imports = PathsConfig{
		App:      c.Module + appPathRel,
		CMD:      c.Module + cmdPathRel,
		Handlers: c.Module + handlersPathRel,
		Models:   c.Module + modelsPathRel,
		Repos:    c.Module + reposPathRel,
	}
}

func CreateProject(config *ProjectConfig) error {
	err := CreateDir(config.RootPath)
	if err != nil {
		return errors.Wrapf(err, "failed to create project directory '%s'", config.RootPath)
	}

	return nil
}
