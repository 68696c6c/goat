package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
)

type CRUD[Model, Request any] interface {
	Saver[Model]
	Identifier[Model]
	Filterer[Model]
	Creator[Model, Request]
	Updater[Model, Request]
	Deleter[Model]
}

type Saver[Model any] interface {
	Save(cx context.Context, m Model) error
}

type Identifier[Model any] interface {
	GetByID(cx context.Context, id goat.ID, loadRelations ...bool) (Model, error)
}

type Filterer[Model any] interface {
	Filter(cx context.Context, q query.Builder) ([]Model, query.Builder, error)
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

func Filter[M any](db *gorm.DB, q query.Builder) ([]*M, query.Builder, error) {
	var result []*M
	err := goat.ApplyQueryToGorm(db, q, false)
	if err != nil {
		return []*M{}, q, errors.Wrap(err, "failed get build query")
	}

	pagination, err := paginate(db, q.GetPagination())
	if err != nil {
		return []*M{}, q, errors.Wrap(err, "failed get pagination total count")
	}

	result, err = filter[M](db)
	if err != nil {
		return result, q, err
	}

	return result, q.Pagination(pagination), nil
}

func paginate(db *gorm.DB, p *query.Pagination) (*query.Pagination, error) {
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return p, err
	}
	return query.NewPagination().SetPage(p.GetPage()).SetPageSize(p.GetPageSize()).SetTotal(int(count)), nil
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

func Create[M any](db *gorm.DB, input M) error {
	return handleGormErrors(db.Create(input))
}

func Update[M any](db *gorm.DB, input M) error {
	return handleGormErrors(db.Save(input))
}

func Delete[M any](db *gorm.DB, input M) error {
	return handleGormErrors(db.Delete(input))
}

func handleGormErrors(db *gorm.DB) error {
	err := db.Error
	if err != nil {
		return err
	}
	return nil
}
