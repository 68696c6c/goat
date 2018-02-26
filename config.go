package goat

import (
	"github.com/spf13/viper"
)

var (
	readConfig        = true
	configFile        string
	haveConfigFile    bool
	defaultConfigFile = "config.yml"
)

func ReadConfig(b bool) {
	readConfig = b
}

func SetConfigFile(s string) error {
	if haveConfigFile {
		return addAndGetError("config already set")
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
	if readConfig {
		c := Root() + "/" + GetConfigFile()
		viper.SetConfigFile(c)
		if err := viper.ReadInConfig(); err != nil {
			addError("failed to load config: " + err.Error())
		}
	}
}
