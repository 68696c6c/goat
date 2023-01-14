package database

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/icrowley/fake"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDBConfig(t *testing.T, key string) ConnectionConfig {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	c := ConnectionConfig{
		Debug:    true,
		Host:     fake.Word(),
		Port:     1234,
		Username: fake.Word(),
		Password: fake.Word(),
		// MultiStatements: true,
	}

	err := os.Setenv(key+"_DEBUG", "1")
	require.Nil(t, err, fmt.Sprintf("failed to set env variable '%s_debug'", key))

	err = os.Setenv(key+"_HOST", c.Host)
	require.Nil(t, err, fmt.Sprintf("failed to set env variable '%s_host'", key))

	err = os.Setenv(key+"_PORT", "1234")
	require.Nil(t, err, fmt.Sprintf("failed to set env variable '%s_port'", key))

	err = os.Setenv(key+"_DATABASE", c.Database)
	require.Nil(t, err, fmt.Sprintf("failed to set env variable '%s_database'", key))

	err = os.Setenv(key+"_USERNAME", c.Username)
	require.Nil(t, err, fmt.Sprintf("failed to set env variable '%s_username'", key))

	err = os.Setenv(key+"_PASSWORD", c.Password)
	require.Nil(t, err, fmt.Sprintf("failed to set env variable '%s_password'", key))

	err = os.Setenv(key+"_MULTI_STATEMENTS", "1")
	require.Nil(t, err, fmt.Sprintf("failed to set env variable '%s_multi_statements'", key))

	return c
}

func assertConnectionEqual(t *testing.T, config, result ConnectionConfig) {
	assert.True(t, config.Debug, "unexpected config value for 'debug'")
	assert.Equal(t, config.Host, result.Host, "unexpected config value for 'host'")
	assert.Equal(t, config.Port, result.Port, "unexpected config value for 'port'")
	assert.Equal(t, config.Database, result.Database, "unexpected config value for 'database'")
	assert.Equal(t, config.Username, result.Username, "unexpected config value for 'username'")
	assert.Equal(t, config.Password, result.Password, "unexpected config value for 'password'")
	// assert.True(t, config.MultiStatements, "unexpected config value for 'multi_statements'")
}

func TestGetDBConfig_Default(t *testing.T) {
	config := setupDBConfig(t, "DB")
	result := GetMainDBConfig()
	assertConnectionEqual(t, config, result)
}

func TestGetDBConfig_Custom(t *testing.T) {
	config := setupDBConfig(t, "CONFIG_TEST")
	result := getDBConfig("config_test")
	assertConnectionEqual(t, config, result)
}

func TestConnectionConfig_String(t *testing.T) {
	host := fake.Word()
	port := 1234
	database := fake.Word()
	username := fake.Word()
	password := fake.Word()
	debug := true
	// multi := true

	c := ConnectionConfig{
		Debug:    debug,
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
		// MultiStatements: multi,
	}

	// e := fmt.Sprintf("Host: %v, Port: %v, Database: %v, Username: %v, Password: %v, Debug: true, MultiStatements: true", host, port, database, username, password)
	e := fmt.Sprintf("Host: %v, Port: %v, Database: %v, Username: %v, Password: %v, Debug: true", host, port, database, username, password)
	s := c.String()
	assert.Equal(t, e, s, "unexpected value returned")
}

func TestNewServiceGORM_MainDB(t *testing.T) {
	config := setupDBConfig(t, dbMainConnectionKey)

	s := NewService(Config{
		MainConnectionConfig: config,
	})

	connections := s.getConnections()

	require.Len(t, connections, 1, "unexpected number of connections")

	result, ok := connections[dbMainConnectionKey]
	require.True(t, ok, "failed to get main connection")

	assertConnectionEqual(t, config, result)
}

func TestNewServiceGORM_GetCustomDB(t *testing.T) {
	s := NewService(Config{
		MainConnectionConfig: setupDBConfig(t, dbMainConnectionKey),
	})

	config := setupDBConfig(t, "CUSTOM_CONNECTION")

	// This will definitely return an error but that isn't the point of this test.
	_, _ = s.GetCustomDB("custom_connection")

	connections := s.getConnections()
	require.Len(t, connections, 2, "unexpected number of connections")

	result, ok := connections["custom_connection"]
	require.True(t, ok, "failed to get custom connection")

	assertConnectionEqual(t, config, result)
}
