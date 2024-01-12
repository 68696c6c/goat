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
		Host:     "test-db",
		Port:     3306,
		Username: "root",
		Password: "secret",
	})

	db, err := subject.GetMainDB()
	require.Nil(t, err)
	require.NotNil(t, db)
}

func Test_Service_ConnectionStringTLS(t *testing.T) {
	cases := []struct {
		name          string
		mode          TLSMode
		expectedValue string
	}{
		{"default tls", TLSMode(""), ""},
		{"strict", TLSModeStrict, "true"},
		{"skip verify", TLSModeSkipVerify, "skip-verify"},
		{"preferred", TLSModePreferred, "preferred"},
		{"off", TLSModeOff, "false"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			conf := Config{Host: "h", Port: 3306, Username: "u", Password: "p", Database: "d"}
			if testCase.mode != "" {
				conf.TLS = testCase.mode
			}

			dsn := conf.ConnectionString()
			parsed, err := url.Parse(dsn)

			require.NoError(t, err)
			require.Equal(t, testCase.expectedValue, parsed.Query().Get("tls"))
		})
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
