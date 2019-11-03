package main

import (
	"github.com/68696c6c/goat/src/app"
	"strings"

	"github.com/68696c6c/goat/src/http"
	"github.com/68696c6c/goat/src/sys"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
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

func GetRouter(setRoutes http.RouterInitializer, getApp app.Initializer) http.Router {
	return g.HTTP.NewRouter(setRoutes, getApp)
}

func GetLogger() *logrus.Logger {
	return g.Log.NewLogger()
}

func GetFileLogger(name string) (*logrus.Logger, error) {
	return g.Log.NewFileLogger(name)
}

func ErrorIfProd() error {
	if g.Env == sys.EnvironmentProd {
		return errors.Errorf("app environment is set to '%s'", sys.EnvironmentProd.String())
	}
	return nil
}

func DebugEnabled() bool {
	return g.HTTP.DebugEnabled()
}

func BindMiddleware(r interface{}) gin.HandlerFunc {
	return g.HTTP.BindMiddleware(r)
}

// Returns a random string that can be used as a Basic Auth token.
func GenerateToken() string {
	u := uuid.New().String()
	return strings.Replace(u, "-", "", -1)
}
