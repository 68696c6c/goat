package goat

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/68696c6c/goat/query"
)

type Repo[M Model] interface {
	Filter(q *query.Query) ([]M, error)
	GetById(id ID, loadRelations ...bool) (M, error)
	Save(m M) error
	Delete(m M) error
}

type queryFactory func() *gorm.DB

func PaginatedFilter[M Model](queryFilter *query.Query, getBaseQuery queryFactory, result []M) ([]M, error) {
	dbQuery, err := queryFilter.ApplyToGorm(getBaseQuery())
	if err != nil {
		return result, errors.Wrap(err, "failed to build filter query")
	}

	err = dbQuery.Find(&result).Error
	if err != nil && ErrorBesidesRecordNotFound(err) {
		return result, errors.Wrap(err, "failed to execute filter query")
	}

	err = ApplyPaginationToQuery(queryFilter, getBaseQuery())
	if err != nil {
		return result, err
	}

	return result, nil
}

func GetFirst[M Model](db *gorm.DB, result M, where ...any) (M, error) {
	err := db.First(result, where...).Error
	if err != nil {
		if RecordNotFound(err) {
			return nil, err
		}
		return result, err
	}
	return result, nil
}

func handleGormErrors(dbQuery *gorm.DB) error {
	err := dbQuery.Error
	if err != nil {
		return err
	}
	return nil
}

func Create[M Model](db *gorm.DB, input M) error {
	return handleGormErrors(db.Create(input))
}

func Update[M Model](db *gorm.DB, input M) error {
	return handleGormErrors(db.Save(input))
}

func SoftDelete[M ModelSoftDelete](db *gorm.DB, input M) error {
	return handleGormErrors(db.Delete(input))
}

func HardDelete[M ModelHardDelete](db *gorm.DB, input M) error {
	return handleGormErrors(db.Delete(input))
}
