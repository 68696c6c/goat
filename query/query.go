package query

import (
	"fmt"
	"net/url"
	"strings"
)

type Builder interface {
	WhereSQL(sql string, values ...any) Builder

	// Filter methods

	Where(field string, op Operator, value any) Builder
	WhereEq(field string, value any) Builder
	WhereNotEq(field string, value any) Builder
	WhereLt(field string, value any) Builder
	WhereLtEq(field string, value any) Builder
	WhereGt(field string, value any) Builder
	WhereGtEq(field string, value any) Builder
	WhereLike(field string, value any) Builder
	WhereNotLike(field string, value any) Builder
	WhereIn(field string, value any) Builder
	WhereNotIn(field string, value any) Builder
	WhereIs(field string, value any) Builder
	WhereNotIs(field string, value any) Builder
	WhereBetween(field string, valueMin, valueMax any) Builder
	WhereNotBetween(field string, valueMin, valueMax any) Builder

	And(field string, op Operator, value any) Builder
	AndEq(field string, value any) Builder
	AndNotEq(field string, value any) Builder
	AndLt(field string, value any) Builder
	AndLtEq(field string, value any) Builder
	AndGt(field string, value any) Builder
	AndGtEq(field string, value any) Builder
	AndLike(field string, value any) Builder
	AndNotLike(field string, value any) Builder
	AndIn(field string, value any) Builder
	AndNotIn(field string, value any) Builder
	AndIs(field string, value any) Builder
	AndNotIs(field string, value any) Builder
	AndBetween(field string, valueMin, valueMax any) Builder
	AndNotBetween(field string, valueMin, valueMax any) Builder

	Or(field string, op Operator, value any) Builder
	OrEq(field string, value any) Builder
	OrNotEq(field string, value any) Builder
	OrLt(field string, value any) Builder
	OrLtEq(field string, value any) Builder
	OrGt(field string, value any) Builder
	OrGtEq(field string, value any) Builder
	OrLike(field string, value any) Builder
	OrNotLike(field string, value any) Builder
	OrIn(field string, value any) Builder
	OrNotIn(field string, value any) Builder
	OrIs(field string, value any) Builder
	OrNotIs(field string, value any) Builder
	OrBetween(field string, valueMin, valueMax any) Builder
	OrNotBetween(field string, valueMin, valueMax any) Builder

	AndGroup(conditions Filter) Builder
	OrGroup(conditions Filter) Builder

	GetWhere() (string, []any, error)

	// Order methods

	Order(field string, dir ...Direction) Builder
	GetOrderBy() string
	GetOrder() *Order

	// Pagination methods

	Limit(int) Builder
	GetLimit() int
	Offset(int) Builder
	GetOffset() int
	Pagination(p *Pagination) Builder
	GetPagination() *Pagination

	// Join methods

	Join(query string, args ...any) Builder
	GetJoins() []Join

	Build() (Template, error)

	// Raw SQL methods

	ToSQL() (string, []any, error)
	Select(fields ...string) Builder
	From(tableName string) Builder
}

type query struct {
	filter     Filter
	order      *Order
	pagination *Pagination
	joins      []Join
	fields     []string
	from       string
}

func NewQuery() Builder {
	return &query{
		filter:     NewFilter(),
		order:      NewOrder(),
		pagination: NewPagination(),
		joins:      []Join{},
		fields:     []string{},
		from:       "",
	}
}

func NewQueryFromUrl(q url.Values) Builder {
	if q == nil {
		return NewQuery()
	}
	return &query{
		filter:     NewFilter(),
		order:      NewOrderFromUrl(q),
		pagination: NewPaginationFromUrl(q),
		joins:      []Join{},
	}
}

func (q *query) WhereSQL(sql string, values ...any) Builder {
	q.filter.WhereSQL(sql, values...)
	return q
}

// Where aliases
// These methods are just sugar that allow for starting a query with a more natural Where() instead of And().

func (q *query) Where(field string, op Operator, value any) Builder {
	q.filter.And(field, op, value)
	return q
}

func (q *query) WhereEq(field string, value any) Builder {
	q.filter.AndEq(field, value)
	return q
}

func (q *query) WhereNotEq(field string, value any) Builder {
	q.filter.AndNotEq(field, value)
	return q
}

func (q *query) WhereLt(field string, value any) Builder {
	q.filter.AndLt(field, value)
	return q
}

func (q *query) WhereLtEq(field string, value any) Builder {
	q.filter.AndLtEq(field, value)
	return q
}

func (q *query) WhereGt(field string, value any) Builder {
	q.filter.AndGt(field, value)
	return q
}

func (q *query) WhereGtEq(field string, value any) Builder {
	q.filter.AndGtEq(field, value)
	return q
}

func (q *query) WhereLike(field string, value any) Builder {
	q.filter.AndLike(field, value)
	return q
}

func (q *query) WhereNotLike(field string, value any) Builder {
	q.filter.AndNotLike(field, value)
	return q
}

func (q *query) WhereIn(field string, value any) Builder {
	q.filter.AndIn(field, value)
	return q
}

func (q *query) WhereNotIn(field string, value any) Builder {
	q.filter.AndNotIn(field, value)
	return q
}

func (q *query) WhereIs(field string, value any) Builder {
	q.filter.AndIs(field, value)
	return q
}

func (q *query) WhereNotIs(field string, value any) Builder {
	q.filter.AndNotIs(field, value)
	return q
}

func (q *query) WhereBetween(field string, valueMin, valueMax any) Builder {
	q.filter.AndBetween(field, valueMin, valueMax)
	return q
}

func (q *query) WhereNotBetween(field string, valueMin, valueMax any) Builder {
	q.filter.AndNotBetween(field, valueMin, valueMax)
	return q
}

// And conditions

func (q *query) And(field string, op Operator, value any) Builder {
	q.filter.And(field, op, value)
	return q
}

func (q *query) AndEq(field string, value any) Builder {
	q.filter.AndEq(field, value)
	return q
}

func (q *query) AndNotEq(field string, value any) Builder {
	q.filter.AndNotEq(field, value)
	return q
}

func (q *query) AndLt(field string, value any) Builder {
	q.filter.AndLt(field, value)
	return q
}

func (q *query) AndLtEq(field string, value any) Builder {
	q.filter.AndLtEq(field, value)
	return q
}

func (q *query) AndGt(field string, value any) Builder {
	q.filter.AndGt(field, value)
	return q
}

func (q *query) AndGtEq(field string, value any) Builder {
	q.filter.AndGtEq(field, value)
	return q
}

func (q *query) AndLike(field string, value any) Builder {
	q.filter.AndLike(field, value)
	return q
}

func (q *query) AndNotLike(field string, value any) Builder {
	q.filter.AndNotLike(field, value)
	return q
}

func (q *query) AndIn(field string, value any) Builder {
	q.filter.AndIn(field, value)
	return q
}

func (q *query) AndNotIn(field string, value any) Builder {
	q.filter.AndNotIn(field, value)
	return q
}

func (q *query) AndIs(field string, value any) Builder {
	q.filter.AndIs(field, value)
	return q
}

func (q *query) AndNotIs(field string, value any) Builder {
	q.filter.AndNotIs(field, value)
	return q
}

func (q *query) AndBetween(field string, valueMin, valueMax any) Builder {
	q.filter.AndBetween(field, valueMin, valueMax)
	return q
}

func (q *query) AndNotBetween(field string, valueMin, valueMax any) Builder {
	q.filter.AndNotBetween(field, valueMin, valueMax)
	return q
}

// Or conditions

func (q *query) Or(field string, op Operator, value any) Builder {
	q.filter.Or(field, op, value)
	return q
}

func (q *query) OrEq(field string, value any) Builder {
	q.filter.OrEq(field, value)
	return q
}

func (q *query) OrNotEq(field string, value any) Builder {
	q.filter.OrNotEq(field, value)
	return q
}

func (q *query) OrLt(field string, value any) Builder {
	q.filter.OrLt(field, value)
	return q
}

func (q *query) OrLtEq(field string, value any) Builder {
	q.filter.OrLtEq(field, value)
	return q
}

func (q *query) OrGt(field string, value any) Builder {
	q.filter.OrGt(field, value)
	return q
}

func (q *query) OrGtEq(field string, value any) Builder {
	q.filter.OrGtEq(field, value)
	return q
}

func (q *query) OrLike(field string, value any) Builder {
	q.filter.OrLike(field, value)
	return q
}

func (q *query) OrNotLike(field string, value any) Builder {
	q.filter.OrNotLike(field, value)
	return q
}

func (q *query) OrIn(field string, value any) Builder {
	q.filter.OrIn(field, value)
	return q
}

func (q *query) OrNotIn(field string, value any) Builder {
	q.filter.OrNotIn(field, value)
	return q
}

func (q *query) OrIs(field string, value any) Builder {
	q.filter.OrIs(field, value)
	return q
}

func (q *query) OrNotIs(field string, value any) Builder {
	q.filter.OrNotIs(field, value)
	return q
}

func (q *query) OrBetween(field string, valueMin, valueMax any) Builder {
	q.filter.OrBetween(field, valueMin, valueMax)
	return q
}

func (q *query) OrNotBetween(field string, valueMin, valueMax any) Builder {
	q.filter.OrNotBetween(field, valueMin, valueMax)
	return q
}

// Condition groups

func (q *query) AndGroup(conditions Filter) Builder {
	q.filter.AndGroup(conditions)
	return q
}

func (q *query) OrGroup(conditions Filter) Builder {
	q.filter.OrGroup(conditions)
	return q
}

func (q *query) GetWhere() (string, []any, error) {
	return q.filter.Generate()
}

// Order methods

func (q *query) Order(field string, dir ...Direction) Builder {
	q.order.By(field, dir...)
	return q
}

func (q *query) GetOrderBy() string {
	return q.order.Generate()
}

func (q *query) GetOrder() *Order {
	return q.order
}

// Pagination methods

func (q *query) Limit(limit int) Builder {
	q.pagination.SetPageSize(limit)
	return q
}

func (q *query) GetLimit() int {
	return q.pagination.GetPageSize()
}

func (q *query) Offset(offset int) Builder {
	q.pagination.setOffset(offset)
	return q
}

func (q *query) GetOffset() int {
	return q.pagination.getOffset()
}

func (q *query) Pagination(p *Pagination) Builder {
	q.pagination = p
	return q
}

func (q *query) GetPagination() *Pagination {
	return q.pagination
}

// Join methods

type Join struct {
	Query string
	Args  []any
}

func (q *query) Join(query string, args ...any) Builder {
	q.joins = append(q.joins, Join{
		Query: query,
		Args:  args,
	})
	return q
}

func (q *query) GetJoins() []Join {
	return q.joins
}

type Template struct {
	Fields  []string
	From    string
	Where   string
	Params  []any
	OrderBy string
	Joins   []Join
	Limit   int
	Offset  int
}

func (q *query) Build() (Template, error) {
	where, params, err := q.GetWhere()
	if err != nil {
		return Template{}, err
	}
	return Template{
		Fields:  q.fields,
		From:    q.from,
		Where:   where,
		Params:  params,
		OrderBy: q.GetOrderBy(),
		Joins:   q.GetJoins(),
		Limit:   q.GetLimit(),
		Offset:  q.GetOffset(),
	}, nil
}

// Raw SQL methods

// ToSQL builds the query into raw SQL and parameters that can be used with GORM's db.Raw()
func (q *query) ToSQL() (string, []any, error) {
	t, err := q.Build()
	if err != nil {
		return "", []any{}, err
	}
	var clauses []string
	var sqlParams []any
	if len(t.Fields) > 0 {
		clauses = append(clauses, fmt.Sprintf("SELECT %s", strings.Join(t.Fields, ", ")))
	}
	var from []string
	if t.From != "" {
		from = append(from, t.From)
	}
	for _, j := range t.Joins {
		from = append(from, j.Query)
		sqlParams = append(sqlParams, j.Args...)
	}
	if len(from) > 0 {
		clauses = append(clauses, fmt.Sprintf("FROM %s", strings.Join(from, " ")))
	}
	if t.Where != "" {
		clauses = append(clauses, fmt.Sprintf("WHERE %s", t.Where))
		sqlParams = append(sqlParams, t.Params...)
	}
	if t.OrderBy != "" {
		clauses = append(clauses, fmt.Sprintf("ORDER BY %s", t.OrderBy))
	}
	if t.Limit > 0 {
		clauses = append(clauses, fmt.Sprintf("LIMIT %d", t.Limit))
	}
	if t.Offset > 0 {
		clauses = append(clauses, fmt.Sprintf("OFFSET %d", t.Offset))
	}
	return strings.Join(clauses, " "), sqlParams, nil
}

func (q *query) Select(fields ...string) Builder {
	q.fields = append(q.fields, fields...)
	return q
}

func (q *query) From(tableName string) Builder {
	q.from = tableName
	return q
}
