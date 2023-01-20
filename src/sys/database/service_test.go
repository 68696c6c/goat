package database

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/68696c6c/goat/sys/log"
)

func Test_Service_GetMainDB(t *testing.T) {
	subject := setupDbTest(t, Config{
		Debug:    true,
		Host:     "test-db",
		Port:     3306,
		Username: "root",
		Password: "secret",
	})

	db, err := subject.GetMainDB()
	require.Nil(t, err)
	require.NotNil(t, db)
}

func Test_Service_GetMainDB_Invalid(t *testing.T) {
	subject := setupDbTest(t, Config{
		Debug:    true,
		Host:     "invalid",
		Port:     3306,
		Username: "root",
		Password: "secret",
	})

	db, err := subject.GetMainDB()
	require.NotNil(t, err)
	require.Nil(t, db)
}

func setupDbTest(t *testing.T, c Config) Service {
	l, err := log.NewService(log.Config{
		Level:      zap.NewAtomicLevelAt(zap.DebugLevel),
		Stacktrace: false,
	})
	require.Nil(t, err)

	// TODO: this isn't needed here, use in sys.app test
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// viper.AutomaticEnv()
	//
	// err = os.Setenv("DB_DEBUG", "1")
	// require.Nil(t, err, fmt.Sprintf("failed to set env variable 'db_debug'"))
	//
	// err = os.Setenv("DB_HOST", c.Host)
	// require.Nil(t, err, fmt.Sprintf("failed to set env variable 'db_host'"))
	//
	// err = os.Setenv("DB_PORT", "3306")
	// require.Nil(t, err, fmt.Sprintf("failed to set env variable 'db_port'"))
	//
	// err = os.Setenv("DB_DATABASE", c.Database)
	// require.Nil(t, err, fmt.Sprintf("failed to set env variable 'db_database'"))
	//
	// err = os.Setenv("DB_USERNAME", c.Username)
	// require.Nil(t, err, fmt.Sprintf("failed to set env variable 'db_username'"))
	//
	// err = os.Setenv("DB_PASSWORD", c.Password)
	// require.Nil(t, err, fmt.Sprintf("failed to set env variable 'db_password'"))

	return NewService(c, l)
}
