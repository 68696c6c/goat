package main

import (
	"strings"

	"github.com/68696c6c/goat/src/http"
	"github.com/68696c6c/goat/src/sys"

	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var g sys.Goat

// Goat has three primary concerns: database connections and schema management,
// request handling, and logging.  These concerns are encapsulated inside of
// services that are bootstrapped when goat.Init() is called.
func Init() {
	if g != (sys.Goat{}) {
		return
	}
	g = sys.Init()
}

// Global functions for calling encapsulated services.

func GetMainDB() (*gorm.DB, error) {
	return g.DB.GetMainDB()
}

func GetDB(key string) (*gorm.DB, error) {
	return g.DB.GetCustomDB(key)
}

func GetRouter(setRoutes func(http.Router)) http.Router {
	return g.HTTP.NewRouter(setRoutes)
}

func GetLogger() *logrus.Logger {
	return g.Log.NewLogger()
}

func GetFileLogger(name string) (*logrus.Logger, error) {
	return g.Log.NewFileLogger(name)
}

func GetEnv() sys.Environment {
	return g.Env
}

func ErrorIfProd() error {
	if g.Env == sys.EnvironmentProd {
		return errors.Errorf("app environment is set to '%s'", sys.EnvironmentProd.String())
	}
	return nil
}

// Returns a random string that can be used as a Basic Auth token.
func GenerateToken() string {
	u := uuid.New().String()
	return strings.Replace(u, "-", "", -1)
}
