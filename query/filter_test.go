package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Filter_Example(t *testing.T) {
	subject := Where("a", Equals, "").
		AndGroup(
			Where("b", Equals, "").Or("c", LessThan, 100).AndGroup(
				Where("d", Equals, "").Or("e", LessThan, 20),
			).Or("f", LessThan, 3),
		)
	st, params, err := subject.Generate()
	require.Nil(t, err)
	assert.Equal(t, "(a = ? AND (b = ? OR c < ? AND (d = ? OR e < ?) OR f < ?))", st)
	assert.EqualValues(t, []any{"", "", 100, "", 20, 3}, params)
}

func Test_Filter_Example2(t *testing.T) {
	subject := Where("a", Equals, "").
		AndGroup(
			Where("b", Equals, "").Or("c", Equals, 1).AndGroup(
				Where("d", Equals, "").Or("e", GreaterThan, 22),
			).Or("f", NotEqual, 300),
		).
		And("g", Equals, 12)

	st, params, err := subject.Generate()
	require.Nil(t, err)
	assert.Equal(t, "(a = ? AND (b = ? OR c = ? AND (d = ? OR e > ?) OR f != ?) AND g = ?)", st)
	assert.EqualValues(t, []any{"", "", 1, "", 22, 300, 12}, params)
}

func Test_Filter_Where(t *testing.T) {
	testCases := []testCase[Filter, filterResult]{
		{input: Where("a", Equals, 6), expected: filterResult{sql: "(a = ?)", params: []any{6}}},
		{input: WhereEq("a", 55), expected: filterResult{sql: "(a = ?)", params: []any{55}}},
		{input: WhereNotEq("a", "asdf"), expected: filterResult{sql: "(a != ?)", params: []any{"asdf"}}},
		{input: WhereLt("a", 9), expected: filterResult{sql: "(a < ?)", params: []any{9}}},
		{input: WhereLtEq("a", 37), expected: filterResult{sql: "(a <= ?)", params: []any{37}}},
		{input: WhereGt("a", 12), expected: filterResult{sql: "(a > ?)", params: []any{12}}},
		{input: WhereGtEq("a", 49), expected: filterResult{sql: "(a >= ?)", params: []any{49}}},
		{input: WhereLike("a", "example"), expected: filterResult{sql: "(a LIKE ?)", params: []any{"example"}}},
		{input: WhereNotLike("a", "example"), expected: filterResult{sql: "(a NOT LIKE ?)", params: []any{"example"}}},
		{input: WhereIn("a", []any{1, 2, 3}), expected: filterResult{sql: "(a IN (?))", params: []any{[]any{1, 2, 3}}}},
		{input: WhereNotIn("a", []any{1, 2, 3}), expected: filterResult{sql: "(a NOT IN (?))", params: []any{[]any{1, 2, 3}}}},
		{input: WhereIs("a", nil), expected: filterResult{sql: "(a IS ?)", params: []any{nil}}},
		{input: WhereNotIs("a", nil), expected: filterResult{sql: "(a IS NOT ?)", params: []any{nil}}},
		{input: WhereBetween("a", 1, 3), expected: filterResult{sql: "(a BETWEEN ? AND ?)", params: []any{1, 3}}},
		{input: WhereNotBetween("a", 1, 3), expected: filterResult{sql: "(a NOT BETWEEN ? AND ?)", params: []any{1, 3}}},
	}
	for _, testCase := range testCases {
		resultSql, resultParams, err := testCase.input.Generate()
		require.Nil(t, err)
		assert.Equal(t, testCase.expected.sql, resultSql)
		assert.EqualValues(t, testCase.expected.params, resultParams)
	}
}

func Test_Filter_And(t *testing.T) {
	testCases := []testCase[Filter, filterResult]{
		{input: WhereEq("a", 1).And("b", Equals, 6), expected: filterResult{sql: "(a = ? AND b = ?)", params: []any{1, 6}}},
		{input: WhereEq("a", 1).AndEq("b", 55), expected: filterResult{sql: "(a = ? AND b = ?)", params: []any{1, 55}}},
		{input: WhereEq("a", 1).AndNotEq("b", "asdf"), expected: filterResult{sql: "(a = ? AND b != ?)", params: []any{1, "asdf"}}},
		{input: WhereEq("a", 1).AndLt("b", 9), expected: filterResult{sql: "(a = ? AND b < ?)", params: []any{1, 9}}},
		{input: WhereEq("a", 1).AndLtEq("b", 37), expected: filterResult{sql: "(a = ? AND b <= ?)", params: []any{1, 37}}},
		{input: WhereEq("a", 1).AndGt("b", 12), expected: filterResult{sql: "(a = ? AND b > ?)", params: []any{1, 12}}},
		{input: WhereEq("a", 1).AndGtEq("b", 49), expected: filterResult{sql: "(a = ? AND b >= ?)", params: []any{1, 49}}},
		{input: WhereEq("a", 1).AndLike("b", "example"), expected: filterResult{sql: "(a = ? AND b LIKE ?)", params: []any{1, "example"}}},
		{input: WhereEq("a", 1).AndNotLike("b", "example"), expected: filterResult{sql: "(a = ? AND b NOT LIKE ?)", params: []any{1, "example"}}},
		{input: WhereEq("a", 1).AndIn("b", []any{1, 2, 3}), expected: filterResult{sql: "(a = ? AND b IN (?))", params: []any{1, []any{1, 2, 3}}}},
		{input: WhereEq("a", 1).AndNotIn("b", []any{1, 2, 3}), expected: filterResult{sql: "(a = ? AND b NOT IN (?))", params: []any{1, []any{1, 2, 3}}}},
		{input: WhereEq("a", 1).AndIs("b", true), expected: filterResult{sql: "(a = ? AND b IS ?)", params: []any{1, true}}},
		{input: WhereEq("a", 1).AndNotIs("b", true), expected: filterResult{sql: "(a = ? AND b IS NOT ?)", params: []any{1, true}}},
		{input: WhereEq("a", 1).AndBetween("b", 6, 9), expected: filterResult{sql: "(a = ? AND b BETWEEN ? AND ?)", params: []any{1, 6, 9}}},
		{input: WhereEq("a", 1).AndNotBetween("b", 6, 9), expected: filterResult{sql: "(a = ? AND b NOT BETWEEN ? AND ?)", params: []any{1, 6, 9}}},
	}
	for _, testCase := range testCases {
		resultSql, resultParams, err := testCase.input.Generate()
		require.Nil(t, err)
		assert.Equal(t, testCase.expected.sql, resultSql)
		assert.EqualValues(t, testCase.expected.params, resultParams)
	}
}

func Test_Filter_Or(t *testing.T) {
	testCases := []testCase[Filter, filterResult]{
		{input: WhereEq("a", 1).Or("b", Equals, 6), expected: filterResult{sql: "(a = ? OR b = ?)", params: []any{1, 6}}},
		{input: WhereEq("a", 1).OrEq("b", 55), expected: filterResult{sql: "(a = ? OR b = ?)", params: []any{1, 55}}},
		{input: WhereEq("a", 1).OrNotEq("b", "asdf"), expected: filterResult{sql: "(a = ? OR b != ?)", params: []any{1, "asdf"}}},
		{input: WhereEq("a", 1).OrLt("b", 9), expected: filterResult{sql: "(a = ? OR b < ?)", params: []any{1, 9}}},
		{input: WhereEq("a", 1).OrLtEq("b", 37), expected: filterResult{sql: "(a = ? OR b <= ?)", params: []any{1, 37}}},
		{input: WhereEq("a", 1).OrGt("b", 12), expected: filterResult{sql: "(a = ? OR b > ?)", params: []any{1, 12}}},
		{input: WhereEq("a", 1).OrGtEq("b", 49), expected: filterResult{sql: "(a = ? OR b >= ?)", params: []any{1, 49}}},
		{input: WhereEq("a", 1).OrLike("b", "example"), expected: filterResult{sql: "(a = ? OR b LIKE ?)", params: []any{1, "example"}}},
		{input: WhereEq("a", 1).OrNotLike("b", "example"), expected: filterResult{sql: "(a = ? OR b NOT LIKE ?)", params: []any{1, "example"}}},
		{input: WhereEq("a", 1).OrIn("b", []any{1, 2, 3}), expected: filterResult{sql: "(a = ? OR b IN (?))", params: []any{1, []any{1, 2, 3}}}},
		{input: WhereEq("a", 1).OrNotIn("b", []any{1, 2, 3}), expected: filterResult{sql: "(a = ? OR b NOT IN (?))", params: []any{1, []any{1, 2, 3}}}},
		{input: WhereEq("a", 1).OrIs("b", false), expected: filterResult{sql: "(a = ? OR b IS ?)", params: []any{1, false}}},
		{input: WhereEq("a", 1).OrNotIs("b", false), expected: filterResult{sql: "(a = ? OR b IS NOT ?)", params: []any{1, false}}},
		{input: WhereEq("a", 1).OrBetween("b", 6, 9), expected: filterResult{sql: "(a = ? OR b BETWEEN ? AND ?)", params: []any{1, 6, 9}}},
		{input: WhereEq("a", 1).OrNotBetween("b", 6, 9), expected: filterResult{sql: "(a = ? OR b NOT BETWEEN ? AND ?)", params: []any{1, 6, 9}}},
	}
	for _, testCase := range testCases {
		resultSql, resultParams, err := testCase.input.Generate()
		require.Nil(t, err)
		assert.Equal(t, testCase.expected.sql, resultSql)
		assert.EqualValues(t, testCase.expected.params, resultParams)
	}
}

func Test_Filter_AndGroup(t *testing.T) {
	testCases := []testCase[Filter, filterResult]{
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).Or("c", Equals, 6)), expected: filterResult{sql: "(a = ? AND (b = ? OR c = ?))", params: []any{1, 2, 6}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrEq("c", 55)), expected: filterResult{sql: "(a = ? AND (b = ? OR c = ?))", params: []any{1, 2, 55}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrNotEq("c", "asdf")), expected: filterResult{sql: "(a = ? AND (b = ? OR c != ?))", params: []any{1, 2, "asdf"}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrLt("c", 9)), expected: filterResult{sql: "(a = ? AND (b = ? OR c < ?))", params: []any{1, 2, 9}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrLtEq("c", 37)), expected: filterResult{sql: "(a = ? AND (b = ? OR c <= ?))", params: []any{1, 2, 37}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrGt("c", 12)), expected: filterResult{sql: "(a = ? AND (b = ? OR c > ?))", params: []any{1, 2, 12}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrGtEq("c", 49)), expected: filterResult{sql: "(a = ? AND (b = ? OR c >= ?))", params: []any{1, 2, 49}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrLike("c", "example")), expected: filterResult{sql: "(a = ? AND (b = ? OR c LIKE ?))", params: []any{1, 2, "example"}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrNotLike("c", "example")), expected: filterResult{sql: "(a = ? AND (b = ? OR c NOT LIKE ?))", params: []any{1, 2, "example"}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrIn("c", []any{1, 2, 3})), expected: filterResult{sql: "(a = ? AND (b = ? OR c IN (?)))", params: []any{1, 2, []any{1, 2, 3}}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrNotIn("c", []any{1, 2, 3})), expected: filterResult{sql: "(a = ? AND (b = ? OR c NOT IN (?)))", params: []any{1, 2, []any{1, 2, 3}}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrIs("c", nil)), expected: filterResult{sql: "(a = ? AND (b = ? OR c IS ?))", params: []any{1, 2, nil}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrNotIs("c", nil)), expected: filterResult{sql: "(a = ? AND (b = ? OR c IS NOT ?))", params: []any{1, 2, nil}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrBetween("c", 6, 9)), expected: filterResult{sql: "(a = ? AND (b = ? OR c BETWEEN ? AND ?))", params: []any{1, 2, 6, 9}}},
		{input: WhereEq("a", 1).AndGroup(WhereEq("b", 2).OrNotBetween("c", 6, 9)), expected: filterResult{sql: "(a = ? AND (b = ? OR c NOT BETWEEN ? AND ?))", params: []any{1, 2, 6, 9}}},
	}
	for _, testCase := range testCases {
		resultSql, resultParams, err := testCase.input.Generate()
		require.Nil(t, err)
		assert.Equal(t, testCase.expected.sql, resultSql)
		assert.EqualValues(t, testCase.expected.params, resultParams)
	}
}

func Test_Filter_OrGroup(t *testing.T) {
	testCases := []testCase[Filter, filterResult]{
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).Or("c", Equals, 6)), expected: filterResult{sql: "(a = ? OR (b = ? OR c = ?))", params: []any{1, 2, 6}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrEq("c", 55)), expected: filterResult{sql: "(a = ? OR (b = ? OR c = ?))", params: []any{1, 2, 55}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrNotEq("c", "asdf")), expected: filterResult{sql: "(a = ? OR (b = ? OR c != ?))", params: []any{1, 2, "asdf"}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrLt("c", 9)), expected: filterResult{sql: "(a = ? OR (b = ? OR c < ?))", params: []any{1, 2, 9}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrLtEq("c", 37)), expected: filterResult{sql: "(a = ? OR (b = ? OR c <= ?))", params: []any{1, 2, 37}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrGt("c", 12)), expected: filterResult{sql: "(a = ? OR (b = ? OR c > ?))", params: []any{1, 2, 12}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrGtEq("c", 49)), expected: filterResult{sql: "(a = ? OR (b = ? OR c >= ?))", params: []any{1, 2, 49}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrLike("c", "example")), expected: filterResult{sql: "(a = ? OR (b = ? OR c LIKE ?))", params: []any{1, 2, "example"}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrNotLike("c", "example")), expected: filterResult{sql: "(a = ? OR (b = ? OR c NOT LIKE ?))", params: []any{1, 2, "example"}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrIn("c", []any{1, 2, 3})), expected: filterResult{sql: "(a = ? OR (b = ? OR c IN (?)))", params: []any{1, 2, []any{1, 2, 3}}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrNotIn("c", []any{1, 2, 3})), expected: filterResult{sql: "(a = ? OR (b = ? OR c NOT IN (?)))", params: []any{1, 2, []any{1, 2, 3}}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrIs("c", nil)), expected: filterResult{sql: "(a = ? OR (b = ? OR c IS ?))", params: []any{1, 2, nil}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrNotIs("c", nil)), expected: filterResult{sql: "(a = ? OR (b = ? OR c IS NOT ?))", params: []any{1, 2, nil}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrBetween("c", 6, 9)), expected: filterResult{sql: "(a = ? OR (b = ? OR c BETWEEN ? AND ?))", params: []any{1, 2, 6, 9}}},
		{input: WhereEq("a", 1).OrGroup(WhereEq("b", 2).OrNotBetween("c", 6, 9)), expected: filterResult{sql: "(a = ? OR (b = ? OR c NOT BETWEEN ? AND ?))", params: []any{1, 2, 6, 9}}},
	}
	for _, testCase := range testCases {
		resultSql, resultParams, err := testCase.input.Generate()
		require.Nil(t, err)
		assert.Equal(t, testCase.expected.sql, resultSql)
		assert.EqualValues(t, testCase.expected.params, resultParams)
	}
}
