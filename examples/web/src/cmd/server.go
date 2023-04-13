package cmd

import (
	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/68696c6c/web/app"
	"github.com/68696c6c/web/http"
)

func init() {
	Root.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Runs the web server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// config, err := app.GetConfig()
			// if err != nil {
			// 	return errors.Wrap(err, "failed to load app config")
			// }

			services, err := app.InitApp(goat.GetMainDB)
			if err != nil {
				return errors.Wrap(err, "failed to initialize app")
			}

			router, err := http.InitRouter(services)
			if err != nil {
				return errors.Wrap(err, "failed to initialize router")
			}

			err = router.Run()
			if err != nil {
				return errors.Wrap(err, "error starting server")
			}

			return nil
		},
	})
}
