package database

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"

	"github.com/68696c6c/goat/sys/log"
)

type SSLMode string

const (
	SSLModeDefault SSLMode = "" // Use the default value based on the specified dialect.

	// Mysql values.
	SSLModeSkipVerify SSLMode = "skip-verify"
	SSLModePreferred  SSLMode = "preferred" // Mysql default.
	SSLModeStrict     SSLMode = "true"
	SSLModeOff        SSLMode = "false"

	// Postgres values.
	SSLModeDisable    SSLMode = "disable"     // I don't care about security, and I don't want to pay the overhead of encryption.
	SSLModeAllow      SSLMode = "allow"       // I don't care about security, but I will pay the overhead of encryption if the server insists on it.
	SSLModePrefer     SSLMode = "prefer"      // (Postgres default) I don't care about encryption, but I wish to pay the overhead of encryption if the server supports it.
	SSLModeRequire    SSLMode = "require"     // I want my data to be encrypted, and I accept the overhead. I trust that the network will make sure I always connect to the server I want.
	SSLModeVerifyCA   SSLMode = "verify-ca"   // I want my data encrypted, and I accept the overhead. I want to be sure that I connect to a server that I trust.
	SSLModeVerifyFull SSLMode = "verify-full" // I want my data encrypted, and I accept the overhead. I want to be sure that I connect to a server I trust, and that it's the one I specify.
)

type Dialect string

const (
	DialectMysql    Dialect = "mysql"
	DialectPostgres Dialect = "postgres"
)

type Config struct {
	Dialect   Dialect
	Debug     bool
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
	SSL       SSLMode
	BatchSize int
}

func (c Config) String() string {
	return fmt.Sprintf("Dialect: %v, Host: %v, Port: %v, Database: %v, Username: %v, Password: %v, Debug: %v, SSL: %v, BatchSize: %v", c.Dialect, c.Host, c.Port, c.Database, c.Username, c.Password, c.Debug, c.SSL, c.BatchSize)
}

func (c Config) ConnectionString() string {
	if c.Dialect == DialectPostgres {
		return c.postgresConnectionString()
	}
	return c.mysqlConnectionString()
}

func (c Config) mysqlConnectionString() string {
	result := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true", c.Username, c.Password, c.Host, c.Port, c.Database)
	if c.SSL == "" {
		result += fmt.Sprintf("&tls=%s", string(SSLModePreferred))
	} else {
		result += fmt.Sprintf("&tls=%s", string(c.SSL))
	}
	return result
}

func (c Config) postgresConnectionString() string {
	// result := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", c.Host, c.Username, c.Password, c.Database, c.Port)
	result := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", c.Username, c.Password, c.Host, c.Port, c.Database)
	if c.SSL == "" {
		result += fmt.Sprintf("?sslmode=%s", string(SSLModePrefer))
	} else {
		result += fmt.Sprintf("?sslmode=%s", string(c.SSL))
	}
	return result
}

func (c Config) Dialector() gorm.Dialector {
	if c.Dialect == DialectPostgres {
		return postgres.Open(c.postgresConnectionString())
	}
	return mysql.Open(c.mysqlConnectionString())
}

type Service interface {
	GetMainDB() (*gorm.DB, error)
	GetConnection(c Config) (*gorm.DB, error)
}

func NewService(c Config, l log.Service) Service {
	return service{
		config: c,
		log:    l,
	}
}

type service struct {
	config Config
	log    log.Service
}

// GetMainDB returns a new database connection using the configured defaults.
func (s service) GetMainDB() (*gorm.DB, error) {
	connection, err := s.GetConnection(s.config)
	if err != nil {
		msg := fmt.Sprintf("failed to connect to default database using credentials: %s", s.config.String())
		return nil, errors.Wrap(err, msg)
	}
	return connection, nil
}

// GetConnection returns a database connection using the provided configuration.
func (s service) GetConnection(c Config) (*gorm.DB, error) {
	logLevel := gormlog.Error
	if c.Debug {
		logLevel = gormlog.Info
	}

	gormConfig := gorm.Config{
		Logger: s.log.GormLogger().LogMode(logLevel),
	}
	if c.BatchSize > 0 {
		gormConfig.CreateBatchSize = c.BatchSize
	}

	connection, err := gorm.Open(c.Dialector(), &gormConfig)
	if err != nil {
		return nil, err
	}

	return connection, nil
}
