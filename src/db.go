package goat

import (
	"gorm.io/gorm"

	"github.com/68696c6c/goat/sys/database"
)

func GetMainDB() (*gorm.DB, error) {
	return g.DB.GetMainDB()
}

type DatabaseConfig database.Config

func GetDB(c DatabaseConfig) (*gorm.DB, error) {
	return g.DB.GetConnection(database.Config(c))
}
