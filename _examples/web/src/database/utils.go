package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"

	_ "github.com/68696c6c/web/database/migrations"
)

func GetMigrationDB(db *gorm.DB) (*sql.DB, error) {
	err := goose.SetDialect("mysql")
	if err != nil {
		return nil, errors.Wrap(err, "failed to set sql dialect")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get sql db")
	}

	return sqlDB, nil
}

func ResetDB(db *gorm.DB) error {
	sqlDB, err := GetMigrationDB(db)
	if err != nil {
		return errors.Wrap(err, "failed to get sql db")
	}

	err = goose.Run("down-to", sqlDB, ".", "0")
	if err != nil {
		return errors.Wrap(err, "failed to reset database schema")
	}

	err = goose.Run("up", sqlDB, ".")
	if err != nil {
		return errors.Wrap(err, "failed to migrate database")
	}

	return nil
}
