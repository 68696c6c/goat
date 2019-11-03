package generator

import "github.com/pkg/errors"

const packageCMD = "cmd"

const rootTemplate = `
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Root = &cobra.Command{
	Use:   "{{.ModuleName}}",
	Short: "Root command for {{.Name}}",
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("author", "{{.Author.Name}} <{{.Author.Email}}>")
	viper.SetDefault("license", "{{.License}}")
}

`

const serverTemplate = `
package cmd

import (
	"{{.Imports.Packages.App}}"

	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	Root.AddCommand(serverCommand)
}

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Runs the API web server.",
	Long:  "Runs the API web server.",
	Run: func(cmd *cobra.Command, args []string) {

		g := goat.Init()

		// Initialize router.
		router := goat.NewRouter(handlers.InitRoutes, app.GetApp)

		// Run the server.
		err = router.Run(port)
		if err != nil {
			goat.ExitWithError(errors.Wrap(err, "error starting server"))
		}
	},
}

`

func CreateCMD(config *ProjectConfig) error {
	err := CreateDir(config.Paths.CMD)
	if err != nil {
		return errors.Wrapf(err, "failed to create cmd directory '%s'", config.Paths.CMD)
	}

	// Create root command.
	err = GenerateFile(config.Paths.CMD, "root", rootTemplate, config)
	if err != nil {
		return errors.Wrap(err, "failed to create root command")
	}

	// Create server command.
	err = GenerateFile(config.Paths.CMD, "server", serverTemplate, config)
	if err != nil {
		return errors.Wrap(err, "failed to create server command")
	}

	// Create migrate command.

	// Create make:migration command.

	return nil
}
