package goat

import (
	"github.com/spf13/viper"
)

const (
	configFileDefault = "config.yml"
)

var (
	configFileSet      bool
	configFile         string
	configPath         string
	configFilePathType = configPathTypeDefault
	readConfig         = true
)

func initConfig(p pathInterface) error {
	if !readConfig {
		return nil
	}
	switch configFilePathType {
	case configPathTypeDefault:
		configFile = configFileDefault
		configPath = p.RootPath(configFileDefault)
		break
	case configPathTypeRel:
		configPath = p.RootPath(configFile)
		break
	case configPathTypeAbs:
		configPath = configFile
		break
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return addAndGetError("failed to load config: " + err.Error())
	}
	configFileSet = true
	return nil
}

func SetConfigFilePath(path string) error {
	if configFileSet {
		return addAndGetError("config already set")
	}
	configFilePathType = configPathTypeAbs
	configFile = path
	configFileSet = true
	return nil
}

func SetConfigFile(filename string) error {
	if configFileSet {
		return addAndGetError("config already set")
	}
	configFilePathType = configPathTypeRel
	configFile = filename
	configFileSet = true
	return nil
}

func ReadConfig(b bool) {
	if !initialized {
		readConfig = b
		return
	}
	addError("goat.ReadConfig() called after initialization")
}

func ConfigFileName() string {
	return configFile
}

func ConfigFilePath() string {
	return configPath
}
