package app

import (
	"sync"

	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/68696c6c/web/app/lib/auth"
	"github.com/68696c6c/web/app/repos"
)

var container App
var once sync.Once

type App struct {
	Version string
	DB      *gorm.DB
	Auth    auth.Service
	repos.UsersRepo
	repos.OrganizationsRepo
}

func initApp(db *gorm.DB, config Config) (App, error) {
	var err error
	once.Do(func() {
		usersRepo := repos.NewUsersRepo(db)

		authService, err := auth.NewAuthService(config.Auth, usersRepo.GetByEmail)
		if err != nil {
			return
		}

		container = App{
			Version:           config.Version,
			DB:                db,
			Auth:              authService,
			UsersRepo:         usersRepo,
			OrganizationsRepo: repos.NewOrganizationsRepo(db),
		}
	})
	if err != nil {
		return App{}, nil
	}
	return container, nil
}

type DBInitializer func() (*gorm.DB, error)

func InitApp(initDB DBInitializer) (App, error) {
	err := goat.Init()
	if err != nil {
		return App{}, errors.Wrap(err, "failed to initialize goat")
	}

	config, err := GetConfig()
	if err != nil {
		return App{}, errors.Wrap(err, "failed to load app config")
	}

	db, err := initDB()
	if err != nil {
		return App{}, errors.Wrap(err, "failed to initialize database connection")
	}

	app, err := initApp(db, config)
	if err != nil {
		return App{}, errors.Wrap(err, "failed to initialize service container")
	}

	return app, nil
}
