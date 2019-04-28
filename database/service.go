package database

import (
	"fmt"

	"github.com/68696c6c/goat/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	mainDBNameDefault    = "db"
	dbConnectionTemplate = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true"
)

// Goat assumes a primary database connection, but an arbitrary number of
// database connections if needed.  Only MySQL is directly supported, but since
// Goat uses GORM, any database supported by GORM can theoretically be used.
type Service interface {
	GetMainDBName() string
	GetMainDB() (*gorm.DB, error)
	GetCustomDB(key string) (*gorm.DB, error)
}

type ServiceGORM struct {
	mainDBName string
}

func NewServiceGORM(mainDBName string) ServiceGORM {
	n := utils.ArgStringD(mainDBName, mainDBNameDefault)
	return ServiceGORM{
		mainDBName: n,
	}
}

// Returns a database connection using the provided configuration.
func (s ServiceGORM) getConnection(c ConnectionConfig) (*gorm.DB, error) {
	cs := fmt.Sprintf(dbConnectionTemplate, c.Username, c.Password, c.Host, c.Port, c.Database)
	connection, err := gorm.Open("mysql", cs)
	if err != nil {
		return nil, err
	}
	connection.LogMode(c.Debug)

	// Don't set the updated_at column on create.
	connection.Callback().Create().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {
		if !scope.HasError() {
			scope.SetColumn("CreatedAt", gorm.NowFunc())
		}
	})
	return connection, nil
}

// Returns the name of the default connection, e.g. the configuration key that
// the main connection credentials are stored under.
func (s ServiceGORM) GetMainDBName() string {
	return s.mainDBName
}

// Returns a new database connection using the configured defaults.
func (s ServiceGORM) GetMainDB() (*gorm.DB, error) {
	c := GetDBConfig(s.mainDBName)
	connection, err := s.getConnection(c)
	if err != nil {
		t := "failed to connect to default database '%s' using credentials: %s"
		msg := fmt.Sprintf(t, s.mainDBName, c.String())
		return nil, errors.Wrap(err, msg)
	}
	return connection, nil
}

// Returns a new database connection using DB env variables with the provided
// config key.
func (s ServiceGORM) GetCustomDB(key string) (*gorm.DB, error) {
	c := GetDBConfig(key)
	connection, err := s.getConnection(c)
	if err != nil {
		t := "failed to connect to custom database '%s' using credentials: %s"
		msg := fmt.Sprintf(t, key, c.String())
		return nil, errors.Wrap(err, msg)
	}
	return connection, nil
}
