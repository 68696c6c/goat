package query

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/68696c6c/goat/query/filter"
)

type Builder interface {
	GetGormPageQueryTEMP() *Query
	GetPagination() Pagination
	String() string
	Order(field string, dir SortDir)
	ApplyToGorm(g *gorm.DB) (*gorm.DB, error)
	GetGormPageQuery(g *gorm.DB) (*gorm.DB, error)
	ApplyPaginationTotals(totalRecordCount int64)

	WhereEq(field string, value any)
	WhereLike(field string, value any)
	WhereIn(field string, value any)
	WhereLt(field string, value any)
	WhereLtEq(field string, value any)
	WhereGt(field string, value any)
	WhereGtEq(field string, value any)
	WhereNotEq(field string, value any)

	OrWhereEq(field string, value any)
	OrWhereLike(field string, value any)
	OrWhereIn(field string, value any)
	OrWhereLt(field string, value any)
	OrWhereLtEq(field string, value any)
	OrWhereGt(field string, value any)
	OrWhereGtEq(field string, value any)
	OrWhereNotEq(field string, value any)
}

type Query struct {
	Pagination Pagination
	Filter     filter.Filter
	Sort       []Sort
	Preload    []string
}

func NewQueryBuilder(cx *gin.Context) Builder {
	return &Query{
		Pagination: NewPaginationFromGin(cx),
		Filter:     filter.NewFilter(),
		Sort:       NewSort(cx),
	}
}

func (q *Query) GetPagination() Pagination {
	return q.Pagination
}

func (q *Query) sortString() string {
	if len(q.Sort) == 0 {
		return ""
	}
	var ss []string
	for _, s := range q.Sort {
		ss = append(ss, fmt.Sprintf("%s %s", s.Field, s.Dir))
	}
	return strings.Join(ss, ", ")
}

func (q *Query) String() string {
	if q == nil {
		return ""
	}
	pString := q.Pagination.String()
	fString := q.Filter.String()
	var sorts []string
	for _, s := range q.Sort {
		sorts = append(sorts, s.String())
	}
	sString := strings.Join(sorts, "\n ")
	prString := strings.Join(q.Preload, "\n ")
	return fmt.Sprintf("paginator: %v\n filter: %v\n sort: %v\n preload: %v\n", pString, fString, sString, prString)
}

func (q *Query) Order(field string, dir SortDir) {
	q.Sort = append(q.Sort, Sort{
		Field: field,
		Dir:   dir,
	})
}

func (q *Query) ApplyToGorm(g *gorm.DB) (*gorm.DB, error) {
	if q.Filter != nil {
		where, params, err := q.Filter.Apply()
		if err != nil {
			return nil, errors.Wrap(err, "failed to apply filter")
		}

		if where != "" {
			g = g.Where(where, params...)
		}
	}

	if len(q.Sort) > 0 {
		g = g.Order(q.sortString())
	}

	// page := q.Pagination.Page
	// size := q.Pagination.PageSize
	//
	// if size > 0 {
	// 	g = g.Limit(int(size)).Offset(int((page - 1) * size))
	// }

	for _, p := range q.Preload {
		g = g.Preload(p)
	}

	return g, nil
}

func (q *Query) GetGormPageQueryTEMP() *Query {
	return &Query{
		Filter:  q.Filter,
		Preload: q.Preload,
	}
}

// GetGormPageQuery copies the query, removing the pagination and applies it to the provided Gorm instance.
// Can be used to get the un-paginated total count of rows for use in pagination.
// DOES NOT MODIFY THE QUERY INSTANCE
func (q *Query) GetGormPageQuery(g *gorm.DB) (*gorm.DB, error) {
	c := &Query{
		Filter:  q.Filter,
		Preload: q.Preload,
	}
	cg, err := c.ApplyToGorm(g)
	if err != nil {
		return nil, err
	}
	return cg, nil
}

// ApplyPaginationTotals updates the Pagination Total and TotalPages values using the provided new totalRecordCount.
func (q *Query) ApplyPaginationTotals(totalRecordCount int64) {
	original := q.Pagination
	q.Pagination = NewPaginationFromValues(original.Page, original.PageSize, totalRecordCount)
}

// AND Filters

func (q *Query) WhereEq(field string, value any) {
	q.Filter.WhereField(field, filter.OpEq, value)
}

func (q *Query) WhereLike(field string, value any) {
	q.Filter.WhereField(field, filter.OpLike, value)
}

func (q *Query) WhereIn(field string, value any) {
	q.Filter.WhereField(field, filter.OpIn, value)
}

func (q *Query) WhereLt(field string, value any) {
	q.Filter.WhereField(field, filter.OpLt, value)
}

func (q *Query) WhereLtEq(field string, value any) {
	q.Filter.WhereField(field, filter.OpLtEq, value)
}

func (q *Query) WhereGt(field string, value any) {
	q.Filter.WhereField(field, filter.OpGt, value)
}

func (q *Query) WhereGtEq(field string, value any) {
	q.Filter.WhereField(field, filter.OpGtEq, value)
}

func (q *Query) WhereNotEq(field string, value any) {
	q.Filter.WhereField(field, filter.OpNotEq, value)
}

// OR Filters

func (q *Query) OrWhereEq(field string, value any) {
	q.Filter.OrWhereField(field, filter.OpEq, value)
}

func (q *Query) OrWhereLike(field string, value any) {
	q.Filter.OrWhereField(field, filter.OpLike, value)
}

func (q *Query) OrWhereIn(field string, value any) {
	q.Filter.OrWhereField(field, filter.OpIn, value)
}

func (q *Query) OrWhereLt(field string, value any) {
	q.Filter.OrWhereField(field, filter.OpLt, value)
}

func (q *Query) OrWhereLtEq(field string, value any) {
	q.Filter.OrWhereField(field, filter.OpLtEq, value)
}

func (q *Query) OrWhereGt(field string, value any) {
	q.Filter.OrWhereField(field, filter.OpGt, value)
}

func (q *Query) OrWhereGtEq(field string, value any) {
	q.Filter.WhereField(field, filter.OpGtEq, value)
}

func (q *Query) OrWhereNotEq(field string, value any) {
	q.Filter.WhereField(field, filter.OpNotEq, value)
}
