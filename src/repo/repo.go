package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/resource"
)

type CRUD[Model, Request any] interface {
	Maker[Model]
	Saver[Model]
	Identifier[Model]
	Filterer[Model]
	Creator[Model, Request]
	Updater[Model, Request]
	Deleter[Model]
}

type Maker[Model any] interface {
	Make() Model
}

type Saver[Model any] interface {
	Save(cx context.Context, m Model) error
}

type Identifier[Model any] interface {
	GetById(cx context.Context, id goat.ID, loadRelations ...bool) (Model, error)
}

type Filterer[Model any] interface {
	Filter(cx context.Context, q query.Builder, pagination resource.Pagination) ([]Model, resource.Pagination, error)
}

type Creator[Model any, Request any] interface {
	Create(cx context.Context, r Request) (Model, error)
}

type Updater[Model any, Request any] interface {
	Update(cx context.Context, id goat.ID, r Request) (Model, error)
}

type Deleter[Model any] interface {
	Delete(cx context.Context, m Model) error
}

func Filter[M any](db *gorm.DB, queryFilter query.Builder, p resource.Pagination) ([]*M, resource.Pagination, error) {
	var result []*M

	dbQuery, err := goat.ApplyQueryToGormNoLimitOffset(db, queryFilter)
	if err != nil {
		return result, p, errors.Wrap(err, "failed to build filter query")
	}

	pagination, err := paginate(dbQuery, p)
	if err != nil {
		return []*M{}, pagination, errors.Wrap(err, "failed get pagination total count")
	}

	result, err = filter[M](dbQuery)
	if err != nil {
		return result, pagination, err
	}

	return result, pagination, nil
}

func paginate(db *gorm.DB, p resource.Pagination) (resource.Pagination, error) {
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return p, err
	}
	return resource.NewPaginationFromValues(p.Page, p.PageSize, count), nil
}

func filter[M any](db *gorm.DB) ([]*M, error) {
	var result []*M
	err := db.Find(&result).Error
	if err != nil && goat.ErrorBesidesRecordNotFound(err) {
		return result, errors.Wrap(err, "failed to execute filter query")
	}
	return result, nil
}

func First[M any](db *gorm.DB, where ...any) (*M, error) {
	var result M
	err := db.First(&result, where...).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func handleGormErrors(db *gorm.DB) error {
	err := db.Error
	if err != nil {
		return err
	}
	return nil
}

func Create[M any](db *gorm.DB, input M) error {
	return handleGormErrors(db.Create(input))
}

func Update[M any](db *gorm.DB, input M) error {
	return handleGormErrors(db.Save(input))
}

func Delete[M any](db *gorm.DB, input M) error {
	return handleGormErrors(db.Delete(input))
}
