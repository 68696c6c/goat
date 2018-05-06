package goat

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConnectionTemplate    = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true"
	panicOnFailedConnection = true
)

type DBConfig struct {
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	Debug           bool
	MultiStatements bool
}

func (d *DBConfig) String() string {
	return fmt.Sprintf("Host: %v, Port: %v, Username: %v, Password: %v, Database: %v, Debug: %v", d.Host, d.Port, d.Username, d.Password, d.Database, d.Debug)
}

func GetDefaultDBConfig() DBConfig {
	mustBeInitialized()
	return DBConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Database: viper.GetString("db.database"),
		Username: viper.GetString("db.username"),
		Password: url.QueryEscape(viper.GetString("db.password")),
		Debug:    viper.GetBool("db.debug"),
	}
}

// Set whether or not to panic if a database connection fails.  Default is true.
// Will panic if goat has not been initialized.
func SetDBPanicMode(b bool) {
	mustBeInitialized()
	panicOnFailedConnection = b
}

// Database connection constructor.  Will attempt to connect to a database using
// connection info from the app config.  Will panic if goat has not been
// initialized and add an error to the error stack if the connection fails.
func NewDB() (*gorm.DB, error) {
	dbConfig := GetDefaultDBConfig()
	return NewCustomDB(dbConfig)
}

// Returns a new database connection using the provided connection info.
// Will panic if goat has not been initialized and add an error to the error
// stack if the connection fails.
func NewCustomDB(c DBConfig) (*gorm.DB, error) {
	mustBeInitialized()
	cs := fmt.Sprintf(dbConnectionTemplate, c.Username, c.Password, c.Host, c.Port, c.Database)
	connection, err := gorm.Open("mysql", cs)
	if err != nil {
		msg := "failed to connect to database: " + err.Error()
		if panicOnFailedConnection {
			panic(msg)
		} else {
			err := addAndGetError(msg)
			return nil, err
		}
	}
	connection.LogMode(c.Debug)
	return connection, nil
}

func RecordNotFound(errs []error) bool {
	for _, err := range errs {
		if err == gorm.ErrRecordNotFound {
			return true
		}
	}
	return false
}
