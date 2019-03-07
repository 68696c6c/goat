package goat

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	configTestFixtureConfigFile = "fixtures/config.yml"
	configTestFixtureConfigPath = "/go/src/goat/" + configTestFixtureConfigFile
	configTestMsgFileName       = "config.fileName not set correctly"
	configTestMsgFilePath       = "config.filePath not set correctly"
	configTestMsgConfigFileSet  = "configFileSet not set correctly"
	configTestMsgSetConfigPath  = "calling SetConfigPath returned an unexpected value"
)

var (
	configTestContainer *Container
)

func mockPath() *path {
	_, b, _, ok := runtime.Caller(1)
	if !ok {
		panic("failed to set config test root dir")
	}
	rootPath := filepath.Dir(b)
	exePath := rootPath + "/mock"
	p := newPath(exePath, nil, rootPath, runtime.Caller)
	return p
}

// Reset config variables to simulate a fresh initialization.
func configTestReset() {
	configTestContainer = nil
	initialized = false
	errs = []error{}
	configFileSet = false
	configFile = ""
	configPath = ""
	configFilePathType = configPathTypeDefault
	readConfig = true
}

// A config test analog of goat.Init().
func configTestInit() {
	p := mockPath()
	configTestContainer = newContainer(p, readConfig)
	errs := GetErrors()
	if len(errs) == 0 {
		configTestContainer.Utils.SetInitialized(true)
		return
	}
	errString := ErrorsToString(errs)
	panic("failed to initialize config test: " + errString)
}

func TestInitConfig_Default(t *testing.T) {
	configTestReset()
	configTestInit()
	assert.NotEmpty(t, configTestContainer.Config.FileName(), configTestMsgFileName)
	assert.NotEmpty(t, configTestContainer.Config.FilePath(), configTestMsgFilePath)
	assert.True(t, configFileSet, configTestMsgConfigFileSet)
}

func TestSetConfigFilePath_Success(t *testing.T) {
	configTestReset()

	err := SetConfigFilePath(configTestFixtureConfigPath)
	assert.Nil(t, err, configTestMsgSetConfigPath)

	configTestInit()
	assert.Equal(t, configTestFixtureConfigPath, configTestContainer.Config.FileName(), configTestMsgFileName)
	assert.Equal(t, configTestFixtureConfigPath, configTestContainer.Config.FilePath(), configTestMsgFilePath)
	assert.Equal(t, configFilePathType, configPathTypeAbs, configTestMsgFilePath)
	assert.True(t, configFileSet, configTestMsgConfigFileSet)
}

func TestSetConfigFilePath_Error(t *testing.T) {
	configTestReset()
	configTestInit()

	err := SetConfigFilePath(configTestFixtureConfigPath)
	assert.NotNil(t, err, configTestMsgSetConfigPath)

	p := configTestContainer.Path.RootPath(configFileDefault)
	assert.Equal(t, configFileDefault, configTestContainer.Config.FileName(), configTestMsgFileName)
	assert.Equal(t, p, configTestContainer.Config.FilePath(), configTestMsgFilePath)
	assert.True(t, configFileSet, configTestMsgConfigFileSet)
}

func TestSetConfigFile_Success(t *testing.T) {
	configTestReset()

	err := SetConfigFile(configTestFixtureConfigFile)
	assert.Nil(t, err, configTestMsgSetConfigPath)

	configTestInit()
	assert.Equal(t, configTestFixtureConfigFile, configTestContainer.Config.FileName(), configTestMsgFileName)
	assert.Equal(t, configTestFixtureConfigPath, configTestContainer.Config.FilePath(), configTestMsgFilePath)
	assert.Equal(t, configFilePathType, configPathTypeRel, configTestMsgFilePath)
	assert.True(t, configFileSet, configTestMsgConfigFileSet)
}

func TestSetConfigFile_Error(t *testing.T) {
	configTestReset()
	configTestInit()

	err := SetConfigFile(configTestFixtureConfigFile)
	assert.NotNil(t, err, configTestMsgSetConfigPath)

	p := configTestContainer.Path.RootPath(configFileDefault)
	assert.Equal(t, configFileDefault, configTestContainer.Config.FileName(), configTestMsgFileName)
	assert.Equal(t, p, configTestContainer.Config.FilePath(), configTestMsgFilePath)
	assert.True(t, configFileSet, configTestMsgConfigFileSet)
}

func TestReadConfig_True(t *testing.T) {
	configTestReset()

	ReadConfig(true)

	configTestInit()
	assert.NotEmpty(t, configTestContainer.Config.FileName(), configTestMsgFileName)
	assert.NotEmpty(t, configTestContainer.Config.FilePath(), configTestMsgFilePath)
	assert.True(t, configFileSet, configTestMsgConfigFileSet)
}

func TestReadConfig_False(t *testing.T) {
	configTestReset()

	ReadConfig(false)

	configTestInit()
	assert.Nil(t, configTestContainer.Config, "config was read when it shouldn't have been")
	assert.False(t, configFileSet, configTestMsgConfigFileSet)
}
