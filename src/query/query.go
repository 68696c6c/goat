package query

import (
	"net/url"
	"strings"
)

type Builder interface {
	// Join()
	// Where()
	Filter
	Order(field string, dir Direction) Builder
	// GroupBy()
	Limit(int) Builder
	Offset(int) Builder

	Preload(preload string) Builder
	GetPreload() []string
	GetWhere() (string, []any, error)
	GetOrder() string
	GetLimit() int
	GetOffset() int
}

type query struct {
	Filter
	sort    []Sort
	preload []string
	limit   int
	offset  int
}

func NewQuery() Builder {
	return &query{
		Filter:  NewFilter(),
		sort:    []Sort{},
		preload: []string{},
		limit:   -1,
		offset:  -1,
	}
}

func NewQueryFromUrl(q url.Values) Builder {
	result := NewQuery()

	if q == nil {
		return result
	}

	sortField := q.Get(sortKey)
	if sortField == "" {
		sortField = defaultSortField
	}

	sortDir, err := DirectionFromString(q.Get(sortDirKey))
	if err != nil {
		sortDir = Descending
	}

	result.Order(sortField, sortDir)

	return result
}

func (q *query) Order(field string, dir Direction) Builder {
	q.sort = append(q.sort, NewSort().By(field).Dir(dir))
	return q
}

func (q *query) Limit(limit int) Builder {
	q.limit = limit
	return q
}

func (q *query) GetLimit() int {
	return q.limit
}

func (q *query) Offset(offset int) Builder {
	q.offset = offset
	return q
}

func (q *query) GetOffset() int {
	return q.offset
}

// TODO: should probably accept the same params as gorm Preload
func (q *query) Preload(preload string) Builder {
	q.preload = append(q.preload, preload)
	return q
}

func (q *query) GetPreload() []string {
	return q.preload
}

func (q *query) GetOrder() string {
	if len(q.sort) == 0 {
		return ""
	}
	var result []string
	for _, s := range q.sort {
		result = append(result, s.Generate())
	}
	return strings.Join(result, ", ")
}

func (q *query) GetWhere() (string, []any, error) {
	return q.Filter.Generate()
}
