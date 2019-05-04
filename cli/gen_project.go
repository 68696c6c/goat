package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/68696c6c/goat/cli/generator"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var config *generator.ProjectConfig

func init() {
	Goat.AddCommand(genProject)
}

var genProject = &cobra.Command{
	Use:   "gen:project name config",
	Short: "Creates a new Goat project.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		project := args[0]
		configFile := args[1]
		parseConfig(project, configFile)

		println(fmt.Sprintf("creating project %s from config %s", project, configFile))

		err := os.MkdirAll(project, os.ModePerm)
		handleError(err)

		err = generator.CreateApp(config)
		handleError(err)

		err = generator.CreateCMD(config)
		handleError(err)

		err = generator.CreateModels(config)
		handleError(err)

		fmtProject()

		os.Exit(0)
	},
}

func parseConfig(projectPath, configPath string) {
	yamlFile, err := ioutil.ReadFile(configPath)
	handleError(errors.Wrap(err, "failed read yml spec"))

	config = &generator.ProjectConfig{}
	err = yaml.Unmarshal(yamlFile, config)
	handleError(errors.Wrap(err, "failed parse project spec"))

	config.Path = projectPath
	config.SRCPath = projectPath + "/src"
}

func fmtProject() {
	err := os.Chdir(config.Path)
	handleError(errors.Wrap(err, "failed change into project dir"))

	cmd := exec.Command("go", "fmt", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	handleError(errors.Wrap(err, "failed format project"))
}
