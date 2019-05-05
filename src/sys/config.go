package sys

import (
	"strings"

	"github.com/68696c6c/goat/src/cmd"
	db "github.com/68696c6c/goat/src/database"
	"github.com/68696c6c/goat/src/http"
	log "github.com/68696c6c/goat/src/logging"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var config Config

type Config struct {
	Env  Environment
	CMD  cmd.Config
	DB   db.Config
	HTTP http.Config
	Log  log.Config
}

func mustReadConfig() {
	if config != (Config{}) {
		return
	}

	// Support both config files and env configuration using Viper.
	// Goat uses env configuration by default.
	// To use a config file you will need to tell Viper where to look, e.g:
	// viper.SetDefault("cfgFile", "./config.yml")
	// RootCommand.PersistentFlags().StringVar(&configFile, "config", "./config.yml", "config file (default is ./config.yml)")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Determine the environment that the app is running in.
	e := viper.GetString("env")
	env, err := EnvironmentFromString(e)
	if err != nil {
		panic(errors.Wrapf(err, "failed to determine app environment"))
	}

	// Read database, logger, and router config from the env using Viper.
	// @TODO it would be preferable to read the db connection from Viper here rather than in the database package...
	config = Config{
		Env: env,
		CMD: cmd.Config{},
		DB: db.Config{
			EnvKey: "db",
		},
		HTTP: http.Config{
			Mode:     viper.GetString("mode"),
			Port:     viper.GetString("port"),
			AuthType: viper.GetString("auth_type"),
		},
		Log: log.Config{
			Path:  viper.GetString("log.path"),
			Ext:   viper.GetString("log.ext"),
			Level: viper.GetString("log.level"),
		},
	}
}
