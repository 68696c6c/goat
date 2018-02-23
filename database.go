package goat

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Debug    bool
}

func MustGetDBConnection(c DBConfig) *gorm.DB {
	template := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true"
	cs := fmt.Sprintf(template, c.Username, c.Password, c.Host, c.Port, c.Database)

	connection, err := gorm.Open("mysql", cs)
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	connection.LogMode(c.Debug)

	return connection
}

func MustInitDB() *gorm.DB {
	password := viper.GetString("db.password")
	return MustGetDBConnection(DBConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Database: viper.GetString("db.database"),
		Username: viper.GetString("db.username"),
		Password: url.QueryEscape(password),
		Debug:    viper.GetBool("db.debug"),
	})
}

func RecordNotFound(errs []error) bool {
	for _, err := range errs {
		if err == gorm.ErrRecordNotFound {
			return true
		}
	}
	return false
}