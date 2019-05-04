package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

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

		createApp()
		createCMD()
		createModels()

		fmtProject()

		os.Exit(0)
	},
}

var config *projectConfig

type authorConfig struct {
	Name  string
	Email string
}

type projectConfig struct {
	Path    string
	SRCPath string
	Name    string
	License string
	Author  authorConfig
	Models  []*model
}

func parseConfig(projectPath, configPath string) {
	yamlFile, err := ioutil.ReadFile(configPath)
	handleError(errors.Wrap(err, "failed read yml file"))

	config = &projectConfig{}
	err = yaml.Unmarshal(yamlFile, config)
	handleError(errors.Wrap(err, "failed read yml file"))

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

func createApp() {
	appPath := config.SRCPath + "/app"
	err := os.MkdirAll(appPath, os.ModePerm)
	handleError(err)

	// Create a service container.
	generateFile(appPath, "container", containerTemplate)
}

func createCMD() {
	cmdPath := config.SRCPath + "/cmd"
	err := os.MkdirAll(cmdPath, os.ModePerm)
	handleError(err)

	// Create root command.
	generateFile(cmdPath, "root", rootTemplate)

	// Create server command.
	generateFile(cmdPath, "server", serverTemplate)

	// Create migrate command.

	// Create make:migration command.

}

func generateFile(basePath, fileName, fileTemplate string) {
	t := template.Must(template.New(fileName).Parse(fileTemplate))

	filePath := fmt.Sprintf("%s/%s.go", basePath, fileName)
	f, err := os.Create(filePath)
	handleError(errors.Wrapf(err, "failed create file '%s'", filePath))

	err = t.Execute(f, config)
	handleError(errors.Wrapf(err, "failed write file '%s'", filePath))

	err = f.Close()
	handleError(errors.Wrapf(err, "failed to close file '%s'", filePath))
}

func generateModel(basePath, fileName, fileTemplate string, modelConfig model) {
	t := template.Must(template.New(fileName).Parse(fileTemplate))

	filePath := fmt.Sprintf("%s/%s.go", basePath, fileName)
	f, err := os.Create(filePath)
	handleError(errors.Wrapf(err, "failed create file '%s'", filePath))

	err = t.Execute(f, modelConfig)
	handleError(errors.Wrapf(err, "failed write file '%s'", filePath))

	err = f.Close()
	handleError(errors.Wrapf(err, "failed to close file '%s'", filePath))
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}

func snakeToCamel(input string) string {
	isToUpper := false
	var output string
	for k, v := range input {
		if k == 0 {
			output = strings.ToUpper(string(input[0]))
		} else {
			if isToUpper {
				output += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					output += string(v)
				}
			}
		}
	}
	return output
}
