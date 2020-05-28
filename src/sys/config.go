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

type Config struct {
	Env  Environment
	CMD  cmd.Config
	DB   db.Config
	HTTP http.Config
	Log  log.Config
}

func mustGetConfig() Config {

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
		panic(errors.Wrap(err, "failed to determine app environment"))
	}

	// If a Gin log mode was not specified, assume one based on the environment.
	httpDebug := false
	if hd := viper.GetString("http_debug"); hd == "" {
		httpDebug = DebugFromEnvironment(env)
	} else if hd == "1" {
		httpDebug = true
	}

	// In order to avoid relying on hacky 'base path' assumptions, require the user to provide a path.
	migrationPath := viper.GetString("migration_path")
	if migrationPath == "" {
		panic(errors.Wrap(err, "failed to determine path to migration files"))
	}

	// Read database, logger, and router config from the env using Viper.
	return Config{
		Env: env,
		CMD: cmd.Config{},
		DB: db.Config{
			MainConnectionConfig: db.GetMainDBConfig(),
			MigrationPath:        migrationPath,
		},
		HTTP: http.Config{
			Debug:           httpDebug,
			Host:            viper.GetString("http_host"),
			Port:            viper.GetString("http_port"),
			AuthType:        viper.GetString("auth_type"),
			ExcludedStructs: viper.GetString("binding_excluded_structs"),
		},
		Log: log.Config{
			Path:  viper.GetString("log.path"),
			Ext:   viper.GetString("log.ext"),
			Level: viper.GetString("log.level"),
		},
	}
}
