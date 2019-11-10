package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	dbMainConnectionKey  = "db"
	dbConnectionTemplate = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true"
)

// Goat assumes a primary database connection, but an arbitrary number of
// database connections if needed.  Only MySQL is directly supported, but since
// Goat uses GORM, any database supported by GORM can theoretically be used.
type Service interface {
	GetMainDB() (*gorm.DB, error)
	GetCustomDB(key string) (*gorm.DB, error)
}

type Config struct {
	MainConnectionConfig ConnectionConfig
}

type ServiceGORM struct {
	connections map[string]ConnectionConfig
}

func NewServiceGORM(c Config) ServiceGORM {
	return ServiceGORM{
		connections: map[string]ConnectionConfig{
			dbMainConnectionKey: c.MainConnectionConfig,
		},
	}
}

// Returns a new database connection using the configured defaults.
func (s ServiceGORM) GetMainDB() (*gorm.DB, error) {
	c := s.connections[dbMainConnectionKey]
	connection, err := s.getConnection(c)
	if err != nil {
		t := "failed to connect to default database using credentials: %s"
		return nil, errors.Wrap(err, fmt.Sprintf(t, c.String()))
	}
	return connection, nil
}

// Returns a new database connection using DB env variables with the provided
// config key.
func (s ServiceGORM) GetCustomDB(key string) (*gorm.DB, error) {
	c, ok := s.connections[key]
	if !ok {
		c = getDBConfig(key)
		s.connections[key] = c
	}
	connection, err := s.getConnection(c)
	if err != nil {
		t := "failed to connect to custom database '%s' using credentials: %s"
		return nil, errors.Wrap(err, fmt.Sprintf(t, key, c.String()))
	}
	return connection, nil
}

// Returns a database connection using the provided configuration.
// Even though this function does not use any instance properties, it is still attached to ServiceGORM because other
// implementations will have different connection logic.
func (s ServiceGORM) getConnection(c ConnectionConfig) (*gorm.DB, error) {
	cs := fmt.Sprintf(dbConnectionTemplate, c.Username, c.Password, c.Host, c.Port, c.Database)
	connection, err := gorm.Open("mysql", cs)
	if err != nil {
		return nil, err
	}

	// GORM logging defaults to only errors; calling LogMode(true) enables detailed logging; calling LogMode(false)
	// disables logging entirely, which is not desirable.
	if c.Debug {
		connection.LogMode(c.Debug)
	}

	// Don't set the updated_at column on create.
	connection.Callback().Create().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {
		if !scope.HasError() {
			scope.SetColumn("CreatedAt", gorm.NowFunc())
		}
	})
	return connection, nil
}
