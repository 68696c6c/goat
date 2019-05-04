package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

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
	Use:   "gen:project config",
	Short: "Creates a new Goat project.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configFile := args[0]
		parseConfig(configFile)

		println(fmt.Sprintf("creating project %s from config %s", config.DirName, configFile))

		err := generator.CreateDir(config.Module)
		handleError(err)

		err = generator.CreateApp(config)
		handleError(err)

		err = generator.CreateCMD(config)
		handleError(err)

		err = generator.CreateRepos(config)
		handleError(err)

		err = generator.CreateModels(config)
		handleError(err)

		fmtProject()

		os.Exit(0)
	},
}

func parseConfig(configPath string) {
	yamlFile, err := ioutil.ReadFile(configPath)
	handleError(errors.Wrap(err, "failed read yml spec"))

	config = &generator.ProjectConfig{}
	err = yaml.Unmarshal(yamlFile, config)
	handleError(errors.Wrap(err, "failed parse project spec"))

	d := strings.Split(config.Module, "/")
	config.DirName = d[:len(d)-1][0]
	config.SRCPath = config.Module + "/src"
	config.AppPath = config.SRCPath + "/app"
	config.CMDPath = config.SRCPath + "/cmd"
	config.ModelsPath = config.SRCPath + "/models"
	config.ReposPath = config.SRCPath + "/repos"
}

func fmtProject() {
	err := os.Chdir(config.Module)
	handleError(errors.Wrap(err, "failed change into project dir"))

	cmd := exec.Command("go", "fmt", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	handleError(errors.Wrap(err, "failed format project"))
}
