package query

import "net/url"

type Builder interface {

	// Filter methods

	Where(field string, op Operator, value any) Builder
	WhereEq(field string, value any) Builder
	WhereLike(field string, value any) Builder
	WhereIn(field string, value any) Builder
	WhereLt(field string, value any) Builder
	WhereLtEq(field string, value any) Builder
	WhereGt(field string, value any) Builder
	WhereGtEq(field string, value any) Builder
	WhereNotEq(field string, value any) Builder
	WhereNotLike(field string, value any) Builder
	WhereNotIn(field string, value any) Builder

	And(field string, op Operator, value any) Builder
	AndEq(field string, value any) Builder
	AndLike(field string, value any) Builder
	AndIn(field string, value any) Builder
	AndLt(field string, value any) Builder
	AndLtEq(field string, value any) Builder
	AndGt(field string, value any) Builder
	AndGtEq(field string, value any) Builder
	AndNotEq(field string, value any) Builder
	AndNotLike(field string, value any) Builder
	AndNotIn(field string, value any) Builder

	Or(field string, op Operator, value any) Builder
	OrEq(field string, value any) Builder
	OrLike(field string, value any) Builder
	OrIn(field string, value any) Builder
	OrLt(field string, value any) Builder
	OrLtEq(field string, value any) Builder
	OrGt(field string, value any) Builder
	OrGtEq(field string, value any) Builder
	OrNotEq(field string, value any) Builder
	OrNotLike(field string, value any) Builder
	OrNotIn(field string, value any) Builder

	AndGroup(conditions Filter) Builder
	OrGroup(conditions Filter) Builder

	GetWhere() (string, []any)

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

	Build() Template
}

type query struct {
	filter     Filter
	order      *Order
	pagination *Pagination
	joins      []Join
}

func NewQuery() Builder {
	return &query{
		filter:     NewFilter(),
		order:      NewOrder(),
		pagination: NewPagination(),
		joins:      []Join{},
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

func (q *query) WhereLike(field string, value any) Builder {
	q.filter.AndLike(field, value)
	return q
}

func (q *query) WhereIn(field string, value any) Builder {
	q.filter.AndIn(field, value)
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

func (q *query) WhereNotEq(field string, value any) Builder {
	q.filter.AndNotEq(field, value)
	return q
}

func (q *query) WhereNotLike(field string, value any) Builder {
	q.filter.AndNotLike(field, value)
	return q
}

func (q *query) WhereNotIn(field string, value any) Builder {
	q.filter.AndNotIn(field, value)
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

func (q *query) AndLike(field string, value any) Builder {
	q.filter.AndLike(field, value)
	return q
}

func (q *query) AndIn(field string, value any) Builder {
	q.filter.AndIn(field, value)
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

func (q *query) AndNotEq(field string, value any) Builder {
	q.filter.AndNotEq(field, value)
	return q
}

func (q *query) AndNotLike(field string, value any) Builder {
	q.filter.AndNotLike(field, value)
	return q
}

func (q *query) AndNotIn(field string, value any) Builder {
	q.filter.AndNotIn(field, value)
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

func (q *query) OrLike(field string, value any) Builder {
	q.filter.OrLike(field, value)
	return q
}

func (q *query) OrIn(field string, value any) Builder {
	q.filter.OrIn(field, value)
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

func (q *query) OrNotEq(field string, value any) Builder {
	q.filter.OrNotEq(field, value)
	return q
}

func (q *query) OrNotLike(field string, value any) Builder {
	q.filter.OrNotLike(field, value)
	return q
}

func (q *query) OrNotIn(field string, value any) Builder {
	q.filter.OrNotIn(field, value)
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

// Order methods

func (q *query) Order(field string, dir ...Direction) Builder {
	q.order.By(field, dir...)
	return q
}

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

func (q *query) GetOrderBy() string {
	return q.order.Generate()
}

func (q *query) GetOrder() *Order {
	return q.order
}

func (q *query) GetWhere() (string, []any) {
	return q.filter.Generate()
}

type Template struct {
	Where   string
	Params  []any
	OrderBy string
	Joins   []Join
	Limit   int
	Offset  int
}

func (q *query) Build() Template {
	where, params := q.GetWhere()
	return Template{
		Where:   where,
		Params:  params,
		OrderBy: q.GetOrderBy(),
		Joins:   q.GetJoins(),
		Limit:   q.GetLimit(),
		Offset:  q.GetOffset(),
	}
}
