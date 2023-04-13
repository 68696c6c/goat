package query

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Query_NewQueryFromUrl(t *testing.T) {
	testCases := []testCase[url.Values, Template]{
		{
			input: nil,
			expected: Template{
				Where:   "",
				Params:  nil,
				OrderBy: "",
				Joins:   []Join{},
				Limit:   0,
				Offset:  0,
			},
		},
		{
			input: mustMakeQuery("sort=a&dir=desc"),
			expected: Template{
				Where:   "",
				Params:  nil,
				OrderBy: "a DESC",
				Joins:   []Join{},
				Limit:   0,
				Offset:  0,
			},
		},
		{
			input: mustMakeQuery("sort=a&dir=desc&page=10&size=5"),
			expected: Template{
				Where:   "",
				Params:  nil,
				OrderBy: "a DESC",
				Joins:   []Join{},
				Limit:   5,
				Offset:  45,
			},
		},
	}
	for _, testCase := range testCases {
		result := NewQueryFromUrl(testCase.input).Build()
		assert.Equal(t, testCase.expected.Where, result.Where, "unexpected where")
		assert.Equal(t, testCase.expected.Params, result.Params, "unexpected params")
		assert.Equal(t, testCase.expected.OrderBy, result.OrderBy, "unexpected order")
		assert.Equal(t, testCase.expected.Joins, result.Joins, "unexpected preload")
		assert.Equal(t, testCase.expected.Limit, result.Limit, "unexpected limit")
		assert.Equal(t, testCase.expected.Offset, result.Offset, "unexpected offset")
	}
}

func Test_Query_Build(t *testing.T) {
	result := NewQuery().Order("a", "desc").Join("RelationA", "args").AndEq("a", "example").Limit(3).Offset(6).Build()
	assert.Equal(t, "(a = ?)", result.Where)
	assert.Equal(t, []any{"example"}, result.Params)
	assert.Equal(t, "a DESC", result.OrderBy)
	assert.Equal(t, "RelationA", result.Joins[0].Query)
	assert.Equal(t, []any{"args"}, result.Joins[0].Args)
	assert.Equal(t, 3, result.Limit)
	assert.Equal(t, 6, result.Offset)
}

func Test_Query_FilterExample(t *testing.T) {
	subject := NewQuery().Where("a", Equals, "").
		AndGroup(
			Where("b", Equals, "").Or("c", LessThan, 100).AndGroup(
				Where("d", Equals, "").Or("e", LessThan, 20),
			).Or("f", LessThan, 3),
		)
	st, params := subject.GetWhere()
	assert.Equal(t, "(a = ? AND (b = ? OR c < ? AND (d = ? OR e < ?) OR f < ?))", st)
	assert.EqualValues(t, []any{"", "", 100, "", 20, 3}, params)
}

func Test_Query_FilterExample2(t *testing.T) {
	subject := NewQuery().Where("a", Equals, "").
		AndGroup(
			Where("b", Equals, "").Or("c", Equals, 1).AndGroup(
				Where("d", Equals, "").Or("e", GreaterThan, 22),
			).Or("f", NotEqual, 300),
		).
		And("g", Equals, 12)

	st, params := subject.GetWhere()
	assert.Equal(t, "(a = ? AND (b = ? OR c = ? AND (d = ? OR e > ?) OR f != ?) AND g = ?)", st)
	assert.EqualValues(t, []any{"", "", 1, "", 22, 300, 12}, params)
}

func Test_Query_Where(t *testing.T) {
	testCases := []testCase[Builder, filterResult]{
		{input: NewQuery().Where("a", Equals, 6), expected: filterResult{sql: "(a = ?)", params: []any{6}}},
		{input: NewQuery().WhereEq("a", 55), expected: filterResult{sql: "(a = ?)", params: []any{55}}},
		{input: NewQuery().WhereLike("a", "example"), expected: filterResult{sql: "(a LIKE ?)", params: []any{"example"}}},
		{input: NewQuery().WhereIn("a", []any{1, 2, 3}), expected: filterResult{sql: "(a IN (?))", params: []any{[]any{1, 2, 3}}}},
		{input: NewQuery().WhereLt("a", 9), expected: filterResult{sql: "(a < ?)", params: []any{9}}},
		{input: NewQuery().WhereLtEq("a", 37), expected: filterResult{sql: "(a <= ?)", params: []any{37}}},
		{input: NewQuery().WhereGt("a", 12), expected: filterResult{sql: "(a > ?)", params: []any{12}}},
		{input: NewQuery().WhereGtEq("a", 49), expected: filterResult{sql: "(a >= ?)", params: []any{49}}},
		{input: NewQuery().WhereNotEq("a", "asdf"), expected: filterResult{sql: "(a != ?)", params: []any{"asdf"}}},
		{input: NewQuery().WhereNotLike("a", "example"), expected: filterResult{sql: "(a NOT LIKE ?)", params: []any{"example"}}},
		{input: NewQuery().WhereNotIn("a", []any{1, 2, 3}), expected: filterResult{sql: "(a NOT IN (?))", params: []any{[]any{1, 2, 3}}}},
	}
	for _, testCase := range testCases {
		resultSql, resultParams := testCase.input.GetWhere()
		assert.Equal(t, testCase.expected.sql, resultSql)
		assert.EqualValues(t, testCase.expected.params, resultParams)
	}
}

func Test_Query_And(t *testing.T) {
	testCases := []testCase[Builder, filterResult]{
		{input: NewQuery().WhereEq("a", 1).And("b", Equals, 6), expected: filterResult{sql: "(a = ? AND b = ?)", params: []any{1, 6}}},
		{input: NewQuery().WhereEq("a", 1).AndEq("b", 55), expected: filterResult{sql: "(a = ? AND b = ?)", params: []any{1, 55}}},
		{input: NewQuery().WhereEq("a", 1).AndLike("b", "example"), expected: filterResult{sql: "(a = ? AND b LIKE ?)", params: []any{1, "example"}}},
		{input: NewQuery().WhereEq("a", 1).AndIn("b", []any{1, 2, 3}), expected: filterResult{sql: "(a = ? AND b IN (?))", params: []any{1, []any{1, 2, 3}}}},
		{input: NewQuery().WhereEq("a", 1).AndLt("b", 9), expected: filterResult{sql: "(a = ? AND b < ?)", params: []any{1, 9}}},
		{input: NewQuery().WhereEq("a", 1).AndLtEq("b", 37), expected: filterResult{sql: "(a = ? AND b <= ?)", params: []any{1, 37}}},
		{input: NewQuery().WhereEq("a", 1).AndGt("b", 12), expected: filterResult{sql: "(a = ? AND b > ?)", params: []any{1, 12}}},
		{input: NewQuery().WhereEq("a", 1).AndGtEq("b", 49), expected: filterResult{sql: "(a = ? AND b >= ?)", params: []any{1, 49}}},
		{input: NewQuery().WhereEq("a", 1).AndNotEq("b", "asdf"), expected: filterResult{sql: "(a = ? AND b != ?)", params: []any{1, "asdf"}}},
		{input: NewQuery().WhereEq("a", 1).AndNotLike("b", "example"), expected: filterResult{sql: "(a = ? AND b NOT LIKE ?)", params: []any{1, "example"}}},
		{input: NewQuery().WhereEq("a", 1).AndNotIn("b", []any{1, 2, 3}), expected: filterResult{sql: "(a = ? AND b NOT IN (?))", params: []any{1, []any{1, 2, 3}}}},
	}
	for _, testCase := range testCases {
		resultSql, resultParams := testCase.input.GetWhere()
		assert.Equal(t, testCase.expected.sql, resultSql)
		assert.EqualValues(t, testCase.expected.params, resultParams)
	}
}

func Test_Query_Or(t *testing.T) {
	testCases := []testCase[Builder, filterResult]{
		{input: NewQuery().WhereEq("a", 1).Or("b", Equals, 6), expected: filterResult{sql: "(a = ? OR b = ?)", params: []any{1, 6}}},
		{input: NewQuery().WhereEq("a", 1).OrEq("b", 55), expected: filterResult{sql: "(a = ? OR b = ?)", params: []any{1, 55}}},
		{input: NewQuery().WhereEq("a", 1).OrLike("b", "example"), expected: filterResult{sql: "(a = ? OR b LIKE ?)", params: []any{1, "example"}}},
		{input: NewQuery().WhereEq("a", 1).OrIn("b", []any{1, 2, 3}), expected: filterResult{sql: "(a = ? OR b IN (?))", params: []any{1, []any{1, 2, 3}}}},
		{input: NewQuery().WhereEq("a", 1).OrLt("b", 9), expected: filterResult{sql: "(a = ? OR b < ?)", params: []any{1, 9}}},
		{input: NewQuery().WhereEq("a", 1).OrLtEq("b", 37), expected: filterResult{sql: "(a = ? OR b <= ?)", params: []any{1, 37}}},
		{input: NewQuery().WhereEq("a", 1).OrGt("b", 12), expected: filterResult{sql: "(a = ? OR b > ?)", params: []any{1, 12}}},
		{input: NewQuery().WhereEq("a", 1).OrGtEq("b", 49), expected: filterResult{sql: "(a = ? OR b >= ?)", params: []any{1, 49}}},
		{input: NewQuery().WhereEq("a", 1).OrNotEq("b", "asdf"), expected: filterResult{sql: "(a = ? OR b != ?)", params: []any{1, "asdf"}}},
		{input: NewQuery().WhereEq("a", 1).OrNotLike("b", "example"), expected: filterResult{sql: "(a = ? OR b NOT LIKE ?)", params: []any{1, "example"}}},
		{input: NewQuery().WhereEq("a", 1).OrNotIn("b", []any{1, 2, 3}), expected: filterResult{sql: "(a = ? OR b NOT IN (?))", params: []any{1, []any{1, 2, 3}}}},
	}
	for _, testCase := range testCases {
		resultSql, resultParams := testCase.input.GetWhere()
		assert.Equal(t, testCase.expected.sql, resultSql)
		assert.EqualValues(t, testCase.expected.params, resultParams)
	}
}

func Test_Query_AndGroup(t *testing.T) {
	testCases := []testCase[Builder, filterResult]{
		{input: NewQuery().WhereEq("a", 1).AndGroup(WhereEq("b", 2).Or("c", Equals, 6)), expected: filterResult{sql: "(a = ? AND (b = ? OR c = ?))", params: []any{1, 2, 6}}},
		{input: NewQuery().WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrEq("c", 55)), expected: filterResult{sql: "(a = ? AND (b = ? OR c = ?))", params: []any{1, 2, 55}}},
		{input: NewQuery().WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrLike("c", "example")), expected: filterResult{sql: "(a = ? AND (b = ? OR c LIKE ?))", params: []any{1, 2, "example"}}},
		{input: NewQuery().WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrIn("c", []any{1, 2, 3})), expected: filterResult{sql: "(a = ? AND (b = ? OR c IN (?)))", params: []any{1, 2, []any{1, 2, 3}}}},
		{input: NewQuery().WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrLt("c", 9)), expected: filterResult{sql: "(a = ? AND (b = ? OR c < ?))", params: []any{1, 2, 9}}},
		{input: NewQuery().WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrLtEq("c", 37)), expected: filterResult{sql: "(a = ? AND (b = ? OR c <= ?))", params: []any{1, 2, 37}}},
		{input: NewQuery().WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrGt("c", 12)), expected: filterResult{sql: "(a = ? AND (b = ? OR c > ?))", params: []any{1, 2, 12}}},
		{input: NewQuery().WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrGtEq("c", 49)), expected: filterResult{sql: "(a = ? AND (b = ? OR c >= ?))", params: []any{1, 2, 49}}},
		{input: NewQuery().WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrNotEq("c", "asdf")), expected: filterResult{sql: "(a = ? AND (b = ? OR c != ?))", params: []any{1, 2, "asdf"}}},
		{input: NewQuery().WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrNotLike("c", "example")), expected: filterResult{sql: "(a = ? AND (b = ? OR c NOT LIKE ?))", params: []any{1, 2, "example"}}},
		{input: NewQuery().WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrNotIn("c", []any{1, 2, 3})), expected: filterResult{sql: "(a = ? AND (b = ? OR c NOT IN (?)))", params: []any{1, 2, []any{1, 2, 3}}}},
	}
	for _, testCase := range testCases {
		resultSql, resultParams := testCase.input.GetWhere()
		assert.Equal(t, testCase.expected.sql, resultSql)
		assert.EqualValues(t, testCase.expected.params, resultParams)
	}
}

func Test_Query_OrGroup(t *testing.T) {
	testCases := []testCase[Builder, filterResult]{
		{input: NewQuery().WhereEq("a", 1).OrGroup(WhereEq("b", 2).Or("c", Equals, 6)), expected: filterResult{sql: "(a = ? OR (b = ? OR c = ?))", params: []any{1, 2, 6}}},
		{input: NewQuery().WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrEq("c", 55)), expected: filterResult{sql: "(a = ? OR (b = ? OR c = ?))", params: []any{1, 2, 55}}},
		{input: NewQuery().WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrLike("c", "example")), expected: filterResult{sql: "(a = ? OR (b = ? OR c LIKE ?))", params: []any{1, 2, "example"}}},
		{input: NewQuery().WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrIn("c", []any{1, 2, 3})), expected: filterResult{sql: "(a = ? OR (b = ? OR c IN (?)))", params: []any{1, 2, []any{1, 2, 3}}}},
		{input: NewQuery().WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrLt("c", 9)), expected: filterResult{sql: "(a = ? OR (b = ? OR c < ?))", params: []any{1, 2, 9}}},
		{input: NewQuery().WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrLtEq("c", 37)), expected: filterResult{sql: "(a = ? OR (b = ? OR c <= ?))", params: []any{1, 2, 37}}},
		{input: NewQuery().WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrGt("c", 12)), expected: filterResult{sql: "(a = ? OR (b = ? OR c > ?))", params: []any{1, 2, 12}}},
		{input: NewQuery().WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrGtEq("c", 49)), expected: filterResult{sql: "(a = ? OR (b = ? OR c >= ?))", params: []any{1, 2, 49}}},
		{input: NewQuery().WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrNotEq("c", "asdf")), expected: filterResult{sql: "(a = ? OR (b = ? OR c != ?))", params: []any{1, 2, "asdf"}}},
		{input: NewQuery().WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrNotLike("c", "example")), expected: filterResult{sql: "(a = ? OR (b = ? OR c NOT LIKE ?))", params: []any{1, 2, "example"}}},
		{input: NewQuery().WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrNotIn("c", []any{1, 2, 3})), expected: filterResult{sql: "(a = ? OR (b = ? OR c NOT IN (?)))", params: []any{1, 2, []any{1, 2, 3}}}},
	}
	for _, testCase := range testCases {
		resultSql, resultParams := testCase.input.GetWhere()
		assert.Equal(t, testCase.expected.sql, resultSql)
		assert.EqualValues(t, testCase.expected.params, resultParams)
	}
}

func Test_Query_Pagination(t *testing.T) {
	input := Pagination{
		Page:       1,
		PageSize:   1,
		Total:      1,
		TotalPages: 1,
	}
	result := NewQueryFromUrl(mustMakeQuery("page=42&size=17&total=257")).Pagination(&input).GetPagination()
	assert.Equal(t, input.Page, result.GetPage())
	assert.Equal(t, input.PageSize, result.GetPageSize())
	assert.Equal(t, input.Total, result.GetTotal())
	assert.Equal(t, input.TotalPages, result.GetTotalPages())
}

func Test_Query_GetPagination(t *testing.T) {
	input := mustMakeQuery("page=42&size=17&total=257")
	expected := Pagination{
		Page:       42,
		PageSize:   17,
		Total:      257,
		TotalPages: 16,
	}
	result := NewQueryFromUrl(input).GetPagination()
	assert.Equal(t, expected.Page, result.GetPage())
	assert.Equal(t, expected.PageSize, result.GetPageSize())
	assert.Equal(t, expected.Total, result.GetTotal())
	assert.Equal(t, expected.TotalPages, result.GetTotalPages())
}

func Test_Query_GetOrder(t *testing.T) {
	input := mustMakeQuery("sort=a&dir=desc")
	expected := sort{field: "a", direction: "DESC"}
	result := NewQueryFromUrl(input).GetOrder()
	require.Len(t, result.sort, 1)
	assert.Equal(t, expected.field, result.sort[0].field)
	assert.Equal(t, expected.direction, result.sort[0].direction)
}
