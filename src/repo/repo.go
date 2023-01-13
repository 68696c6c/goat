package repo

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query2"
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
	// Filter(cx context.Context, q resource.Pagination, scopes ...QueryFilter) ([]Model, resource.Pagination, error)
	Filter(cx context.Context, q query2.Builder, pagination resource.Pagination) ([]Model, resource.Pagination, error)
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

// type queryFactory func() *gorm.DB
//
// func Filter[M any](queryFilter query.Builder, getBaseQuery queryFactory) ([]*M, error) {
// 	var result []*M
// 	dbQuery, err := queryFilter.ApplyToGorm(getBaseQuery())
// 	if err != nil {
// 		return result, errors.Wrap(err, "failed to build filter query")
// 	}
//
// 	err = dbQuery.Find(&result).Error
// 	if err != nil && goat.ErrorBesidesRecordNotFound(err) {
// 		return result, errors.Wrap(err, "failed to execute filter query")
// 	}
//
// 	err = goat.ApplyPaginationToQuery(queryFilter, getBaseQuery())
// 	if err != nil {
// 		return result, errors.Wrap(err, "failed to paginate filter query")
// 	}
//
// 	return result, nil
// }

func Filter[M any](db *gorm.DB, queryFilter query2.Builder, p resource.Pagination) ([]*M, resource.Pagination, error) {
	// p := queryFilter.GetPagination()

	var result []*M
	// dbQuery, err := queryFilter.ApplyToGorm(db)
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

	// err = dbQuery.Find(&result).Error
	// if err != nil && goat.ErrorBesidesRecordNotFound(err) {
	// 	return result, errors.Wrap(err, "failed to execute filter query")
	// }
	//
	// err = goat.ApplyPaginationToQuery(queryFilter, getBaseQuery())
	// if err != nil {
	// 	return result, errors.Wrap(err, "failed to paginate filter query")
	// }

	return result, pagination, nil
}

// func Paginate(p resource.Pagination) func(*gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		limit := p.GetLimit()
// 		if limit > 0 {
// 			return db.Limit(limit).Offset(p.GetOffset())
// 		}
// 		return db
// 	}
// }
//
// type QueryFilter func(db *gorm.DB) *gorm.DB
//
// func Filter[M any](db *gorm.DB, filters ...QueryFilter) ([]*M, error) {
// 	// for _, filter := range filters {
// 	// 	db.Scopes(filter)
// 	// }
// 	// var result []*M
// 	// err := db.Find(&result).Error
// 	// if err != nil && goat.ErrorBesidesRecordNotFound(err) {
// 	// 	return result, errors.Wrap(err, "failed to execute filter query")
// 	// }
// 	// return result, nil
// 	applyFilters(db, filters...)
// 	return filter[M](db)
// }
//
// type Query struct {
// 	*gorm.DB
// 	Pagination resource.Pagination
// 	Filters    []QueryScope
// }

func paginate(db *gorm.DB, p resource.Pagination) (resource.Pagination, error) {
	// func paginate(db *gorm.DB, p query.Pagination) (query.Pagination, error) {
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return p, err
	}
	println("COUNT: " + strconv.FormatInt(count, 10))
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

// func applyFilters(db *gorm.DB, filters ...QueryFilter) {
// 	if len(filters) > 0 {
// 		for _, filter := range filters {
// 			db.Scopes(filter)
// 		}
// 	}
// }
//
// func PaginatedFilter[M any](db *gorm.DB, p resource.Pagination, filters ...QueryFilter) ([]*M, resource.Pagination, error) {
// 	// if len(filters) > 0 {
// 	// 	for _, filter := range filters {
// 	// 		db.Scopes(filter)
// 	// 	}
// 	// }
// 	applyFilters(db, filters...)
//
// 	// p := query.Pagination
// 	// var result []*M
// 	// if p == (resource.Pagination{}) {
// 	// 	result, err := filter[M](db)
// 	// 	return result, p, err
// 	// 	// err := query.Find(&result).Error
// 	// 	// if err != nil {
// 	// 	// 	return result, p, err
// 	// 	// }
// 	// 	// return result, p, nil
// 	// }
//
// 	// var count int64
// 	// err := query.Count(&count).Error
// 	pagination, err := paginate(db, p)
// 	if err != nil {
// 		return []*M{}, pagination, errors.Wrap(err, "failed get pagination total count")
// 	}
//
// 	result, err := filter[M](db.Scopes(Paginate(p)))
// 	// err = query.Scopes(Paginate(p)).Find(&result).Error
// 	if err != nil {
// 		return result, pagination, err
// 	}
//
// 	// return result, resource.NewPaginationFromValues(p.Page, p.PageSize, count), nil
// 	return result, pagination, nil
// }
//
// func PaginatedFilter[M any](db *gorm.DB, p resource.Pagination, filters ...QueryScope) ([]*M, resource.Pagination, error) {
// 	var count int64
// 	err := db.Count(&count).Error
// 	if err != nil {
// 		return []*M{}, p, errors.Wrap(err, "failed get pagination total count")
// 	}
//
// 	// err = db.Scopes(Paginate(p)).Find(&result).Error
// 	// if err != nil && goat.ErrorBesidesRecordNotFound(err) {
// 	// 	return result, p, errors.Wrap(err, "failed to execute filter query")
// 	// }
// 	result, err := Filter[M](db.Scopes(Paginate(p)))
// 	if err != nil {
// 		return result, p, err
// 	}
//
// 	return result, resource.NewPaginationFromValues(p.Page, p.PageSize, count), nil
// }
//
// func FilterTemp[M any](queryFilter query.Builder, db *gorm.DB) ([]*M, error) {
// 	// getDb := func() *gorm.DB {
// 	// 	return db
// 	// }
// 	//
// 	// var result []*M
// 	// dbQuery, err := queryFilter.ApplyToGorm(getDb())
// 	// if err != nil {
// 	// 	return result, errors.Wrap(err, "failed to build filter query")
// 	// }
// 	//
// 	// err = dbQuery.Find(&result).Error
// 	// if err != nil && goat.ErrorBesidesRecordNotFound(err) {
// 	// 	return result, errors.Wrap(err, "failed to execute filter query")
// 	// }
// 	//
// 	// err = goat.ApplyPaginationToQuery(queryFilter, getDb())
// 	// if err != nil {
// 	// 	return result, errors.Wrap(err, "failed to paginate filter query")
// 	// }
// 	//
// 	// return result, nil
//
// 	// var result []*M
// 	// dbQuery, err := queryFilter.ApplyToGorm(getDb())
// 	// if err != nil {
// 	// 	return result, errors.Wrap(err, "failed to build filter query")
// 	// }
// 	//
// 	// err = dbQuery.Find(&result).Error
// 	// if err != nil && goat.ErrorBesidesRecordNotFound(err) {
// 	// 	return result, errors.Wrap(err, "failed to execute filter query")
// 	// }
// 	//
// 	// unpaginatedFilter := queryFilter.GetGormPageQueryTEMP()
// 	// println(unpaginatedFilter.String())
// 	// ndb, err := unpaginatedFilter.ApplyToGorm(db)
// 	// if err != nil {
// 	// 	return result, errors.Wrap(err, "failed to paginate filter query")
// 	// }
// 	// var count int64
// 	// err = ndb.Count(&count).Error
// 	// // err = goat.ApplyPaginationToQuery(queryFilter, dbQuery.Session(&gorm.Session{NewDB: true}))
// 	// if err != nil {
// 	// 	return result, errors.Wrap(err, "failed get pagination total count")
// 	// }
// 	// queryFilter.ApplyPaginationTotals(count)
// 	//
// 	// return result, nil
//
// 	// This works, but the totalPages value only works on the first page
// 	var result []*M
// 	dbQuery, err := queryFilter.ApplyToGorm(db)
// 	if err != nil {
// 		return result, errors.Wrap(err, "failed to build filter query")
// 	}
//
// 	var count int64
// 	err = dbQuery.Count(&count).Error
// 	// err = goat.ApplyPaginationToQuery(queryFilter, dbQuery.Session(&gorm.Session{NewDB: true}))
// 	if err != nil {
// 		return result, errors.Wrap(err, "failed get pagination total count")
// 	}
// 	queryFilter.ApplyPaginationTotals(count)
//
// 	err = dbQuery.Find(&result).Error
// 	if err != nil && goat.ErrorBesidesRecordNotFound(err) {
// 		return result, errors.Wrap(err, "failed to execute filter query")
// 	}
//
// 	// unpaginatedFilter := queryFilter.GetGormPageQueryTEMP()
// 	// println(unpaginatedFilter.String())
// 	// ndb, err := unpaginatedFilter.ApplyToGorm(db)
// 	// if err != nil {
// 	// 	return result, errors.Wrap(err, "failed to paginate filter query")
// 	// }
// 	// var count int64
// 	// err = dbQuery.Count(&count).Error
// 	// // err = goat.ApplyPaginationToQuery(queryFilter, dbQuery.Session(&gorm.Session{NewDB: true}))
// 	// if err != nil {
// 	// 	return result, errors.Wrap(err, "failed get pagination total count")
// 	// }
// 	// queryFilter.ApplyPaginationTotals(count)
//
// 	return result, nil
// }

func First[M any](db *gorm.DB, where ...any) (*M, error) {
	var result M
	err := db.First(&result, where...).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func handleGormErrors(dbQuery *gorm.DB) error {
	err := dbQuery.Error
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
