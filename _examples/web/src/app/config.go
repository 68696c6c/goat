package app

import (
	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/68696c6c/web/app/lib/auth"
)

type Config struct {
	Version string
	Auth    auth.Config
}

func GetConfig() (Config, error) {
	result := Config{
		Version: viper.GetString("version"),
		Auth: auth.Config{
			SignatureKey: viper.GetString("auth_signature_key"),
			Clients: []auth.Client{
				{
					ID:     viper.GetString("auth_client_id"),
					Secret: viper.GetString("auth_client_secret"),
					Public: viper.GetBool("auth_client_public"),
				},
			},
		},
	}
	err := validateConfig(result)
	if err != nil {
		return Config{}, err
	}
	return result, nil
}

func validateConfig(c Config) error {
	var errs []error
	if c.Version == "" {
		errs = append(errs, errors.New("invalid version"))
	}
	// TODO: validate the rest of the config
	return goat.ErrorsToError(errs)
}
