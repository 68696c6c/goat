package goat

import (
	"github.com/spf13/viper"
)

var (
	configFile     string
	haveConfigFile bool
	defaultConfigFile = "config.yml"
)

func SetConfigFile(s string) error {
	if haveConfigFile {
		addAndReturnError("config already set")
	}
	configFile = s
	haveConfigFile = true
	return nil
}

func GetConfigFile() string {
	if configFile == "" {
		configFile = defaultConfigFile
		haveConfigFile = true
	}
	return configFile
}

func initConfig() {
	c := Root() + "/" + GetConfigFile()
	viper.SetConfigFile(c)
	if err := viper.ReadInConfig(); err != nil {
		addError("failed to load config: " + err.Error())
	}
}
