package sys

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/68696c6c/goat/sys/database"
	"github.com/68696c6c/goat/sys/http/router"
	log "github.com/68696c6c/goat/sys/logging"
	"github.com/68696c6c/goat/sys/utils"
)

type Config struct {
	HttpDebug bool
	DB        database.Config
	// HTTP      http.Config
	Log    log.Config
	Router router.Config
}

const (
	httpHost          = ""
	httpPort          = "80"
	httpAuthType      = "basic"
	contextKeyRequest = "goat_request"
	contextKeyQuery   = "goat_query"
)

func mustGetConfig() Config {

	// Support both config files and env configuration using Viper.
	// Goat uses env configuration by default.
	// To use a config file you will need to tell Viper where to look, e.g:
	// viper.SetDefault("cfgFile", "./config.yml")
	// RootCommand.PersistentFlags().StringVar(&configFile, "config", "./config.yml", "config file (default is ./config.yml)")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Enable viper.Get type inferring; needed by the utils.EnvDefault helper.
	viper.SetTypeByDefaultValue(true)

	baseUrl, err := url.Parse(utils.EnvOrDefault[string]("base_url", ""))
	if err != nil {
		panic(errors.Wrap(err, "missing base_url"))
	}

	httpDebug := viper.GetBool("http_debug")

	// Read database, logger, and router config from the env using Viper.
	return Config{
		HttpDebug: httpDebug,
		DB: database.Config{
			MainConnectionConfig: database.GetMainDBConfig(),
		},
		// HTTP: http.Config{
		// 	// BaseUrl:         baseUrl,
		// 	Debug: httpDebug,
		// 	// Host:            utils.EnvOrDefault("http_host", httpHost),
		// 	// Port:            utils.EnvOrDefault("http_port", httpPort),
		// 	// AuthType:        utils.EnvOrDefault("http_auth_type", httpAuthType),
		// 	ExcludedStructs: viper.GetString("http_binding_excluded_structs"),
		// },
		Log: log.Config{
			Path:  viper.GetString("log.path"),
			Ext:   viper.GetString("log.ext"),
			Level: viper.GetString("log.level"),
		},
		Router: router.Config{
			BaseUrl: baseUrl,
			// TODO: rename these
			Debug:    httpDebug,
			Host:     utils.EnvOrDefault("http_host", httpHost),
			Port:     utils.EnvOrDefault("http_port", httpPort),
			AuthType: utils.EnvOrDefault("http_auth_type", httpAuthType),
		},
	}
}
