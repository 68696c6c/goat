package database

import (
	"fmt"

	"gorm.io/gorm/logger"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/68696c6c/goat/query"
)

const (
	dbMainConnectionKey  = "db"
	dbConnectionTemplate = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true"
)

// Goat assumes a primary database connection, but an arbitrary number of
// database connections if needed.  Only MySQL is directly supported, but since
// Goat uses Gorm, any database supported by Gorm can theoretically be used.
type Service interface {
	GetMainDB() (*gorm.DB, error)
	GetMigrationDB() (*gorm.DB, error)
	GetCustomDB(key string) (*gorm.DB, error)
	GetConnection(c ConnectionConfig) (*gorm.DB, error)
	ApplyPaginationToQuery(q *query.Query, baseGormQuery *gorm.DB) error
}

type Config struct {
	MainConnectionConfig ConnectionConfig
}

type ServiceGORM struct {
	connections map[string]ConnectionConfig
	dialect     string
}

func NewServiceGORM(c Config) ServiceGORM {
	return ServiceGORM{
		connections: map[string]ConnectionConfig{
			dbMainConnectionKey: c.MainConnectionConfig,
		},
		dialect: "mysql",
	}
}

// Returns a new database connection using the configured defaults.
func (s ServiceGORM) GetMainDB() (*gorm.DB, error) {
	c := s.connections[dbMainConnectionKey]
	connection, err := s.GetConnection(c)
	if err != nil {
		t := "failed to connect to default database using credentials: %s"
		return nil, errors.Wrap(err, fmt.Sprintf(t, c.String()))
	}
	return connection, nil
}

// Returns a new database connection using the configured defaults, but supporting multi-statements for running migrations.
func (s ServiceGORM) GetMigrationDB() (*gorm.DB, error) {
	c := s.connections[dbMainConnectionKey]
	c.MultiStatements = true
	c.Debug = true
	connection, err := s.GetConnection(c)
	if err != nil {
		t := "failed to connect to default migration database using credentials: %s"
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
	connection, err := s.GetConnection(c)
	if err != nil {
		t := "failed to connect to custom database '%s' using credentials: %s"
		return nil, errors.Wrap(err, fmt.Sprintf(t, key, c.String()))
	}
	return connection, nil
}

// GetConnection Returns a database connection using the provided configuration.
// Even though this function does not use any instance properties, it is still attached to ServiceGorm because other
// implementations will have different connection logic.
func (s ServiceGORM) GetConnection(c ConnectionConfig) (*gorm.DB, error) {
	cs := fmt.Sprintf(dbConnectionTemplate, c.Username, c.Password, c.Host, c.Port, c.Database)

	// By default, Gorm logs slow queries and errors.
	logLevel := logger.Error
	if c.Debug {
		logLevel = logger.Info
	}

	connection, err := gorm.Open(mysql.Open(cs), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	return connection, nil
}

// ApplyPaginationToQuery returns a Gorm query with the provided filter Query applied and updates the Query pagination to match the new
func (s ServiceGORM) ApplyPaginationToQuery(q *query.Query, baseGormQuery *gorm.DB) error {
	pageQuery, err := q.GetGormPageQuery(baseGormQuery)
	if err != nil {
		return err
	}

	var count int64
	err = pageQuery.Count(&count).Error
	if err != nil {
		return errors.Wrap(err, "failed to execute filter count query")
	}

	q.ApplyPaginationTotals(count)

	return nil
}
