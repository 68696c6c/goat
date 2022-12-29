package sys

import (
	"strings"

	"github.com/68696c6c/goat/src/cmd"
	db "github.com/68696c6c/goat/src/database"
	"github.com/68696c6c/goat/src/http"
	log "github.com/68696c6c/goat/src/logging"

	"github.com/spf13/viper"
)

type Config struct {
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

	// Read database, logger, and router config from the env using Viper.
	return Config{
		CMD: cmd.Config{},
		DB: db.Config{
			MainConnectionConfig: db.GetMainDBConfig(),
		},
		HTTP: http.Config{
			Debug:           viper.GetBool("http_debug"),
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
