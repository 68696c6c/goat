package database

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"

	"github.com/68696c6c/goat/sys/log"
)

type TLSMode string

const (
	TLSModeNone       = ""
	TLSModeSkipVerify = "skip-verify"
	TLSModePreferred  = "preferred"
	TLSModeStrict     = "true"
	TLSModeOff        = "false"
)

type Config struct {
	Debug    bool
	Host     string
	Port     int
	Database string
	Username string
	Password string
	TLS      TLSMode
}

func (c Config) String() string {
	return fmt.Sprintf("Host: %v, Port: %v, Database: %v, Username: %v, Password: %v, Debug: %v, TLS: %v", c.Host, c.Port, c.Database, c.Username, c.Password, c.Debug, c.TLS)
}

func (c Config) ConnectionString() string {
	base := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true", c.Username, c.Password, c.Host, c.Port, c.Database)
	if c.TLS != "" {
		base += fmt.Sprintf("&tls=%s", string(c.TLS))
	}

	return base
}

type Service interface {
	GetMainDB() (*gorm.DB, error)
	GetConnection(c Config) (*gorm.DB, error)
}

func NewService(c Config, l log.Service) Service {
	return service{
		config:  c,
		dialect: "mysql",
		log:     l,
	}
}

type service struct {
	config  Config
	dialect string
	log     log.Service
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

	connection, err := gorm.Open(mysql.Open(c.ConnectionString()), &gorm.Config{
		Logger: s.log.GormLogger().LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	return connection, nil
}
