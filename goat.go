package goat

import (
	"strings"

	"github.com/68696c6c/goat/http"
	"github.com/68696c6c/goat/logging"
	"github.com/68696c6c/goat/sys"

	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var container sys.Container

// Goat has three primary concerns: database connections and schema management,
// request handling, and logging.  These concerns are encapsulated inside of
// services that are bootstrapped when goat.Init() is called.
func Init() {
	if container != (sys.Container{}) {
		return
	}

	// Support both config files and env configuration using Viper.
	// Goat uses env configuration by default.
	// To use a config file you will need to tell Viper where to look, e.g:
	// viper.SetDefault("cfgFile", "./config.yml")
	// RootCommand.PersistentFlags().StringVar(&configFile, "config", "./config.yml", "config file (default is ./config.yml)")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read database, logger, and router config from the env using Viper.
	// @TODO it would be preferable to read the db connection from Viper here rather than in the database package...
	dbName := "db"
	// @TODO how to allow configuration of boolean http settings?
	httpConfig := http.Config{
		Mode:     viper.GetString("mode"),
		Port:     viper.GetString("listen_port"),
		AuthType: viper.GetString("auth_type"),
	}
	loggerConfig := logging.LoggerConfig{
		Path:  viper.GetString("log.path"),
		Ext:   viper.GetString("log.ext"),
		Level: viper.GetString("log.level"),
	}

	container = sys.NewContainer(dbName, httpConfig, loggerConfig)
}

// Global functions for calling encapsulated services.

func GetMainDB() (*gorm.DB, error) {
	return container.DatabaseService.GetMainDB()
}

func GetCustomDB(key string) (*gorm.DB, error) {
	return container.DatabaseService.GetCustomDB(key)
}

func GetRouter(setRoutes func(http.Router)) http.Router {
	return container.HTTPService.NewRouter(setRoutes)
}

func GetLogger() *logrus.Logger {
	return container.LoggingService.NewLogger()
}

func GetFileLogger(name string) (*logrus.Logger, error) {
	return container.LoggingService.NewFileLogger(name)
}

// Returns a random string that can be used as a Basic Auth token.
func GenerateToken() string {
	u := uuid.New().String()
	return strings.Replace(u, "-", "", -1)
}
