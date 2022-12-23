package database

import (
	"fmt"

	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/utils"

	"github.com/68696c6c/goose"
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
	GetMigrationDB() (*gorm.DB, error)
	GetCustomDB(key string) (*gorm.DB, error)
	GetSchema(connection *gorm.DB) (goose.SchemaInterface, error)
	GetConnection(c ConnectionConfig) (*gorm.DB, error)
	ApplyPaginationToQuery(q *query.Query, baseGormQuery *gorm.DB) error
}

type Config struct {
	MainConnectionConfig ConnectionConfig
	MigrationPath        string
}

type ServiceGORM struct {
	connections   map[string]ConnectionConfig
	migrationPath string
	dialect       string
}

func NewServiceGORM(c Config) ServiceGORM {
	return ServiceGORM{
		connections: map[string]ConnectionConfig{
			dbMainConnectionKey: c.MainConnectionConfig,
		},
		migrationPath: c.MigrationPath,
		dialect:       "mysql",
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

func (s ServiceGORM) GetSchema(connection *gorm.DB) (goose.SchemaInterface, error) {
	// @TODO handle logging better?
	schema, err := goose.NewSchema(connection, s.migrationPath, nil)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

// Returns a database connection using the provided configuration.
// Even though this function does not use any instance properties, it is still attached to ServiceGORM because other
// implementations will have different connection logic.
func (s ServiceGORM) GetConnection(c ConnectionConfig) (*gorm.DB, error) {
	cs := fmt.Sprintf(dbConnectionTemplate, c.Username, c.Password, c.Host, c.Port, c.Database)
	connection, err := gorm.Open(s.dialect, cs)
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

// Returns a GORM query with the provided filter Query applied and updates the Query pagination to match the new
func (s ServiceGORM) ApplyPaginationToQuery(q *query.Query, baseGormQuery *gorm.DB) error {
	pageQuery, err := q.GetGormPageQuery(baseGormQuery)
	if err != nil {
		return err
	}

	var count uint
	errs := pageQuery.Count(&count).GetErrors()
	if len(errs) > 0 {
		err := utils.ErrorsToError(errs)
		return errors.Wrap(err, "failed to execute filter sites count query")
	}

	q.ApplyPaginationTotals(count)

	return nil
}
