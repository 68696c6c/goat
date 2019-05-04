package cli

const containerTemplate = `
package app

import (
	"github.com/68696c6c/goat"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var container ServiceContainer

type ServiceContainer struct {
	DB     *gorm.DB
	Logger *logrus.Logger
}

// Initializes the service container if it hasn't been already.
func GetApp(l *logrus.Logger) (ServiceContainer, error) {
	if container != (ServiceContainer{}) {
		return container, nil
	}

	db, err := goat.GetMainDB()
	if err != nil {
		return ServiceContainer{}, err
	}

	container = ServiceContainer{
		DB:     db,
		Logger: l,
	}
	return container, nil
}

`
