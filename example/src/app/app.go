package app

import (
	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/68696c6c/example/app/repos"
)

var container App

type App struct {
	Version string
	DB      *gorm.DB
	// Logger  *zap.SugaredLogger
	// Errors  goat.ErrorHandler
	repos.UsersRepo
	repos.OrganizationsRepo
}

func GetApp(db *gorm.DB, config Config) (App, error) {
	if container != (App{}) {
		return container, nil
	}

	container = App{
		Version: config.Version,
		DB:      db,
		// Logger:            logger,
		// Errors:            goat.NewErrorHandler(logger),
		UsersRepo:         repos.NewUsersRepo(db),
		OrganizationsRepo: repos.NewOrganizationsRepo(db),
	}

	return container, nil
}

func InitApp() (App, error) {
	goat.Init()

	db, err := goat.GetMainDB()
	if err != nil {
		return App{}, errors.Wrap(err, "failed to initialize database connection")
	}

	config, err := GetConfig()
	if err != nil {
		return App{}, errors.Wrap(err, "failed to load app config")
	}

	app, err := GetApp(db, config)
	if err != nil {
		return App{}, errors.Wrap(err, "failed to initialize service container")
	}

	return app, nil
}
