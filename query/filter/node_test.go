package filter

//
//import (
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func Test_FilterManualConstruction(t *testing.T) {
//	f := node{
//		Logic: LogicAnd,
//		Children: []*node{
//			{
//				Condition: &fieldCondition{
//					Field: "a",
//					Op:    OpEq,
//					Value: "",
//				},
//			},
//			{
//				Logic: LogicAnd,
//				Children: []*node{
//					{
//						Condition: &fieldCondition{
//							Field: "b",
//							Op:    OpEq,
//							Value: "",
//						},
//					},
//					{
//						Condition: &fieldCondition{
//							Field: "c",
//							Op:    OpLt,
//							Value: 10,
//						},
//						Logic: LogicOr,
//					},
//					{
//						Logic: LogicAnd,
//						Children: []*node{
//							{
//								Condition: &fieldCondition{
//									Field: "d",
//									Op:    OpEq,
//									Value: "",
//								},
//							},
//							{
//								Condition: &fieldCondition{
//									Field: "e",
//									Op:    OpLt,
//									Value: 10,
//								},
//								Logic: LogicOr,
//							},
//						},
//					},
//					{
//						Condition: &fieldCondition{
//							Field: "f",
//							Op:    OpLt,
//							Value: 10,
//						},
//						Logic: LogicOr,
//					},
//				},
//			},
//		},
//	}
//
//	st, params, err := f.Apply()
//	assert.Nil(t, err)
//	assert.Equal(t, "(a = ? AND (b = ? OR c < ? AND (d = ? OR e < ? ) OR f < ? ) ) ", st)
//	assert.EqualValues(t, []interface{}{"", "", 10, "", 10, 10}, params)
//}
//
//func Test_FilterAutoConstruct(t *testing.T) {
//	f := &node{Logic: LogicAnd}
//	f.WhereFieldEq("a", "")
//	f.Group(LogicAnd).
//		WhereFieldEq("b", "").
//		OrWhereFieldEq("c", 10).
//		Group(LogicAnd).
//		WhereFieldEq("d", "").
//		OrWhereField("e", OpGt, 10).
//		Parent.OrWhereField("f", OpNeq, 10).
//		Parent.WhereFieldEq("g", 12)
//
//	st, params, err := f.Apply()
//	assert.Nil(t, err)
//	assert.Equal(t, "(a = ? AND (b = ? OR c = ? AND (d = ? OR e > ? ) OR f <> ? ) AND g = ? ) ", st)
//	assert.EqualValues(t, []interface{}{"", "", 10, "", 10, 10, 12}, params)
//}
