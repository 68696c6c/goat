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

func TestGetDBConfig_Env(t *testing.T) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	host := fake.Word()
	port := "1234"
	database := fake.Word()
	username := fake.Word()
	password := fake.Word()
	debug := "1"
	multi := "1"

	err := os.Setenv("CONFIG_TEST_HOST", host)
	require.Nil(t, err, "failed to set env variable 'db_host'")

	err = os.Setenv("CONFIG_TEST_PORT", port)
	require.Nil(t, err, "failed to set env variable 'db_port'")

	err = os.Setenv("CONFIG_TEST_DATABASE", database)
	require.Nil(t, err, "failed to set env variable 'db_database'")

	err = os.Setenv("CONFIG_TEST_USERNAME", username)
	require.Nil(t, err, "failed to set env variable 'db_username'")

	err = os.Setenv("CONFIG_TEST_PASSWORD", password)
	require.Nil(t, err, "failed to set env variable 'db_password'")

	err = os.Setenv("CONFIG_TEST_DEBUG", debug)
	require.Nil(t, err, "failed to set env variable 'db_debug'")

	err = os.Setenv("CONFIG_TEST_MULTI_STATEMENTS", multi)
	require.Nil(t, err, "failed to set env variable 'multi_statements'")

	c := GetDBConfig("config_test")
	assert.Equal(t, host, c.Host, "unexpected config value for 'host'")
	assert.Equal(t, 1234, c.Port, "unexpected config value for 'port'")
	assert.Equal(t, database, c.Database, "unexpected config value for 'database'")
	assert.Equal(t, username, c.Username, "unexpected config value for 'username'")
	assert.Equal(t, password, c.Password, "unexpected config value for 'password'")
	assert.True(t, c.Debug, "unexpected config value for 'debug'")
	assert.True(t, c.MultiStatements, "unexpected config value for 'multi_statements'")
}

func TestConnectionConfig_String(t *testing.T) {
	host := fake.Word()
	port := 1234
	database := fake.Word()
	username := fake.Word()
	password := fake.Word()
	debug := true
	multi := true

	c := ConnectionConfig{
		Host:            host,
		Port:            port,
		Database:        database,
		Username:        username,
		Password:        password,
		Debug:           debug,
		MultiStatements: multi,
	}

	e := fmt.Sprintf("Host: %v, Port: %v, Database: %v, Username: %v, Password: %v, Debug: true, MultiStatements: true", host, port, database, username, password)
	s := c.String()
	assert.Equal(t, e, s, "unexpected value returned")
}
