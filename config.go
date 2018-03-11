package goat

import (
	"github.com/spf13/viper"
	"goat/types"
)

const (
	configFileDefault = "config.yml"
)

var (
	config         *types.Config
	configFileSet  bool
	configFile     string
	configPath     string
	configPathType = types.ConfigPathTypeDefault
	readConfig     = true
)

func initConfig() (*types.Config, error) {
	switch configPathType {
	case types.ConfigPathTypeDefault:
		configPath = RootPath(configFileDefault)
		break
	case types.ConfigPathTypeRel:
		configPath = RootPath(configFile)
		break
	case types.ConfigPathTypeAbs:
		configPath = configFile
		break
	}
	config = types.NewConfig(configFile, configPath)
	viper.SetConfigFile(config.FilePath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, addAndGetError("failed to load config: " + err.Error())
	}
	return config, nil
}

func SetConfigFilePath(path string) error {
	if configFileSet {
		return addAndGetError("config already set")
	}
	configPathType = types.ConfigPathTypeAbs
	configFile = path
	configFileSet = true
	return nil
}

func SetConfigFile(filename string) error {
	if configFileSet {
		return addAndGetError("config already set")
	}
	configPathType = types.ConfigPathTypeRel
	configFile = filename
	configFileSet = true
	return nil
}

func GetConfig() (*types.Config, error) {
	mustBeInitialized()
	return config, nil
}

func ReadConfig(b bool) {
	if !initialized {
		readConfig = b
		return
	}
	addError("goat.ReadConfig() called after initialization")
}
