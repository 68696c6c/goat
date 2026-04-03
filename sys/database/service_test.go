package database

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/68696c6c/goat/sys/log"
)

func Test_Service_GetMainDB(t *testing.T) {
	subject := setupDbTest(t, Config{
		Debug:    true,
		Host:     "test-db-mysql",
		Port:     3306,
		Username: "root",
		Password: "secret",
	})

	db, err := subject.GetMainDB()
	require.Nil(t, err)
	require.NotNil(t, db)
}

func Test_Service_GetMainDB_Mysql(t *testing.T) {
	subject := setupDbTest(t, Config{
		Dialect:   DialectMysql,
		Debug:     true,
		Host:      "test-db-mysql",
		Port:      3306,
		Username:  "root",
		Password:  "secret",
		Database:  "goat",
		BatchSize: 1000,
	})

	db, err := subject.GetMainDB()
	require.Nil(t, err)
	require.NotNil(t, db)
}

func Test_Service_GetMainDB_Postgres(t *testing.T) {
	subject := setupDbTest(t, Config{
		Dialect:   DialectPostgres,
		Debug:     true,
		Host:      "test-db-postgres",
		Port:      5432,
		Username:  "postgres",
		Password:  "secret",
		Database:  "goat",
		BatchSize: 1000,
	})

	db, err := subject.GetMainDB()
	require.Nil(t, err)
	require.NotNil(t, db)
}

func Test_Service_ConnectionStringSSL(t *testing.T) {
	cases := map[Dialect][]struct {
		name          string
		mode          SSLMode
		expectedValue string
	}{
		DialectMysql: {
			{"mysql: default ssl", SSLMode(""), string(SSLModePreferred)},
			{"mysql: strict", SSLModeStrict, "true"},
			{"mysql: skip verify", SSLModeSkipVerify, "skip-verify"},
			{"mysql: preferred", SSLModePreferred, "preferred"},
			{"mysql: off", SSLModeOff, "false"},
		},
		DialectPostgres: {
			{"postgres: default ssl", SSLMode(""), string(SSLModePrefer)},
			{"postgres: disable", SSLModeDisable, "disable"},
			{"postgres: allow", SSLModeAllow, "allow"},
			{"postgres: require", SSLModeRequire, "require"},
			{"postgres: verify-ca", SSLModeVerifyCA, "verify-ca"},
			{"postgres: verify-full", SSLModeVerifyFull, "verify-full"},
		},
	}

	for dialect, testCases := range cases {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				conf := Config{Dialect: dialect, Host: "h", Port: 3306, Username: "u", Password: "p", Database: "d", SSL: testCase.mode}

				dsn := conf.ConnectionString()
				parsed, err := url.Parse(dsn)
				require.NoError(t, err)

				key := "tls"
				if conf.Dialect == DialectPostgres {
					key = "sslmode"
				}
				require.Equal(t, testCase.expectedValue, parsed.Query().Get(key))
			})
		}
	}
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
	return NewService(c, l)
}
