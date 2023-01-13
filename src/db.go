package goat

import (
	"gorm.io/gorm"

	"github.com/68696c6c/goat/sys/database"
)

func GetMainDB() (*gorm.DB, error) {
	return g.DB.GetMainDB()
}

func GetMigrationDB() (*gorm.DB, error) {
	return g.DB.GetMigrationDB()
}

func GetDB(key string) (*gorm.DB, error) {
	return g.DB.GetCustomDB(key)
}

type DatabaseConnection database.ConnectionConfig

func GetCustomDB(c DatabaseConnection) (*gorm.DB, error) {
	return g.DB.GetConnection(database.ConnectionConfig(c))
}
