package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func Test_Filter_ManualConstruction(t *testing.T) {
// 	subject := Where("a", Equals, "").
// 		Group(And).And("b", Equals, "").Or("c", LessThan, 10).
// 		Group(And).And("d", Equals, "").Or("e", LessThan, 10).GroupEnd().
// 		Or("f", LessThan, 10).
// 		GroupEnd()
//
// 	st, params, err := subject.Generate()
// 	assert.Nil(t, err)
// 	assert.Equal(t, "(a = ? AND (b = ? OR c < ? AND (d = ? OR e < ?) OR f < ?))", st)
// 	assert.EqualValues(t, []any{"", "", 10, "", 10, 10}, params)
// }

func Test_Filter(t *testing.T) {
	subject := Where("a", Equals, "").
		AndGroup(
			Where("b", Equals, "").Or("c", LessThan, 100).AndGroup(
				Where("d", Equals, "").Or("e", LessThan, 20),
			).Or("f", LessThan, 3),
		)
	st, params, err := subject.Generate()
	assert.Nil(t, err)
	assert.Equal(t, "(a = ? AND (b = ? OR c < ? AND (d = ? OR e < ?) OR f < ?))", st)
	assert.EqualValues(t, []any{"", "", 100, "", 20, 3}, params)
}

// func Test_Filter_AutoConstruct(t *testing.T) {
// 	subject := Where("a", Equals, "").
// 		Group(And).
// 		And("b", Equals, "").
// 		Or("c", Equals, 1).
// 		Group(And).
// 		And("d", Equals, "").
// 		Or("e", GreaterThan, 22).
// 		GroupEnd().Or("f", NotEqual, 300).
// 		GroupEnd().And("g", Equals, 12)
//
// 	st, params, err := subject.Generate()
// 	assert.Nil(t, err)
// 	assert.Equal(t, "(a = ? AND (b = ? OR c = ? AND (d = ? OR e > ?) OR f != ?) AND g = ?)", st)
// 	assert.EqualValues(t, []any{"", "", 1, "", 22, 300, 12}, params)
// }

func Test_Filter_AutoConstruct_Group(t *testing.T) {
	subject := Where("a", Equals, "").
		AndGroup(
			Where("b", Equals, "").Or("c", Equals, 1).AndGroup(
				Where("d", Equals, "").Or("e", GreaterThan, 22),
			).Or("f", NotEqual, 300),
		).
		And("g", Equals, 12)

	st, params, err := subject.Generate()
	assert.Nil(t, err)
	assert.Equal(t, "(a = ? AND (b = ? OR c = ? AND (d = ? OR e > ?) OR f != ?) AND g = ?)", st)
	assert.EqualValues(t, []any{"", "", 1, "", 22, 300, 12}, params)
}
