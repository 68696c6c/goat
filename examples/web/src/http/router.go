package http

import (
	"github.com/68696c6c/goat"
	"github.com/pkg/errors"

	"github.com/68696c6c/web/app"
	"github.com/68696c6c/web/app/controllers"
	"github.com/68696c6c/web/app/lib/auth"
	"github.com/68696c6c/web/app/models"
	"github.com/68696c6c/web/http/middlewares"
)

func InitRouter(services app.App) (goat.Router, error) {
	router, err := goat.InitRouter()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize router")
	}

	router.GET(goat.HealthHandler("/health"))
	router.GET(goat.VersionHandler("/version", services.Version))

	api := router.Group("/api")
	{
		routes := api.Group(auth.TokenLinkKey, "/tokens")
		controller := controllers.NewTokensController(services.UsersRepo, services.Auth)
		routes.POST("/authorize", controller.Authorize)
		goat.SetUrl(auth.AuthorizeLinkKey, goat.GetUrl(auth.TokenLinkKey).JoinPath("/authorize"))
		routes.POST("/exchange", controller.Exchange)

		// This route is only for testing the authentication process.  It is probably unnecessary in a real application.
		routes.GET("/", middlewares.ValidateJwt(services.Auth, services.UsersRepo), controller.View)

		// These routes are only for example purposes.  In a real application, they should exist as part of the front-end.
		routes.GET("/new", controller.AuthorizeForm)
		routes.GET("/callback", controller.Callback)
		goat.SetUrl(auth.AuthorizeCallbackLinkKey, goat.GetUrl(auth.TokenLinkKey).JoinPath("/callback"))
	}

	api.Use(middlewares.ValidateJwt(services.Auth, services.UsersRepo))
	{
		routes := api.Group(models.UserLinkKey, "/users")
		controller := controllers.NewUsersController(services.UsersRepo, services.Auth)
		routes.GET("", controller.List)
		routes.GET("/:id", controller.View)
		routes.POST("", controller.Create)
		routes.PUT("/:id", controller.Update)
		routes.DELETE("/:id", controller.Delete)
	}
	{
		routes := api.Group(models.OrganizationLinkKey, "/organizations")
		controller := controllers.NewOrganizationsController(services.OrganizationsRepo, services.Auth)
		routes.GET("", controller.List)
		routes.GET("/:id", controller.View)
		routes.POST("", controller.Create)
		routes.PUT("/:id", controller.Update)
		routes.DELETE("/:id", controller.Delete)
	}

	return router, nil
}
