package database

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

type ConnectionConfig struct {
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	Debug           bool
	MultiStatements bool
}

func (c ConnectionConfig) String() string {
	t := "Host: %v, Port: %v, Database: %v, Username: %v, Password: %v, Debug: %v, MultiStatements: %v"
	return fmt.Sprintf(t, c.Host, c.Port, c.Database, c.Username, c.Password, c.Debug, c.MultiStatements)
}

// Returns a database connection config struct using app config values.
func GetDBConfig(name string) ConnectionConfig {
	return ConnectionConfig{
		Host:            viper.GetString(name + ".host"),
		Port:            viper.GetInt(name + ".port"),
		Database:        viper.GetString(name + ".database"),
		Username:        viper.GetString(name + ".username"),
		Password:        url.QueryEscape(viper.GetString(name + ".password")),
		Debug:           viper.GetBool(name + ".debug"),
		MultiStatements: viper.GetBool(name + ".multi_statements"),
	}
}
