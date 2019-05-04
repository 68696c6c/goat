package cli

const rootTemplate = `
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Root = &cobra.Command{
	Use:   "{{.Path}}",
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
		logger := g.GetLogger()

		// Initialize service container.
		app, err := app.GetApp(logger)
		if err != nil {
			goat.ExitWithError(err)
		}

		// Initialize router.
		router := goat.NewRouter(handlers.InitRoutes)

		// Add the service container to the Gin registry.
		c := map[string]interface{}{
			"app": app,
		}
		router.SetRegistry(c)

		// Run the server.
		err = router.Run(port)
		if err != nil {
			goat.ExitWithError(errors.Wrap(err, "error starting server"))
		}
	},
}

`
