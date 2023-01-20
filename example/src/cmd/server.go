package cmd

import (
	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/68696c6c/example/app"
	"github.com/68696c6c/example/app/controllers"
	"github.com/68696c6c/example/app/models"
)

func init() {
	Root.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Runs the web server",
		Run: func(cmd *cobra.Command, args []string) {
			services, err := app.InitApp()
			if err != nil {
				goat.ExitError(errors.Wrap(err, "failed to initialize app"))
			}

			router := goat.InitRouter()

			router.GET(goat.HealthHandler("/health"))
			router.GET(goat.VersionHandler("/version", services.Version))

			api := router.Group("/api")

			api.Use(controllers.SetCurrentUser())
			{
				routes := api.Group(models.UserLinkKey, "/users")
				controller := controllers.NewUsersController(services.UsersRepo)
				routes.GET("", controller.List)
				routes.GET("/:id", controller.View)
				routes.POST("", controller.Create)
				routes.PUT("/:id", controller.Update)
				routes.DELETE("/:id", controller.Delete)
			}
			{
				routes := api.Group(models.OrganizationLinkKey, "/organizations")
				controller := controllers.NewOrganizationsController(services.OrganizationsRepo)
				routes.GET("", controller.List)
				routes.GET("/:id", controller.View)
				routes.POST("", controller.Create)
				routes.PUT("/:id", controller.Update)
				routes.DELETE("/:id", controller.Delete)
			}

			err = router.Run()
			if err != nil {
				goat.ExitError(errors.Wrap(err, "error starting server"))
			}
		},
	})
}
