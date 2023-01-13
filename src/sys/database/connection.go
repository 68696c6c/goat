package database

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

type ConnectionConfig struct {
	Debug           bool
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	MultiStatements bool
}

func (c ConnectionConfig) String() string {
	t := "Host: %v, Port: %v, Database: %v, Username: %v, Password: %v, Debug: %v, MultiStatements: %v"
	return fmt.Sprintf(t, c.Host, c.Port, c.Database, c.Username, c.Password, c.Debug, c.MultiStatements)
}

// GetMainDBConfig returns the default database connection.
func GetMainDBConfig() ConnectionConfig {
	return getDBConfig(dbMainConnectionKey)
}

// getDBConfig returns a database connection config struct using app config values.
func getDBConfig(name string) ConnectionConfig {
	return ConnectionConfig{
		Debug:           viper.GetBool(name + ".debug"),
		Host:            viper.GetString(name + ".host"),
		Port:            viper.GetInt(name + ".port"),
		Database:        viper.GetString(name + ".database"),
		Username:        viper.GetString(name + ".username"),
		Password:        url.QueryEscape(viper.GetString(name + ".password")),
		MultiStatements: viper.GetBool(name + ".multi_statements"),
	}
}
