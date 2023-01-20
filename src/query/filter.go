package query

import (
	"fmt"
	"strings"

	"github.com/68696c6c/girraph"
)

type Operator string

const (
	Equals             Operator = "="
	LessThan           Operator = "<"
	GreaterThan        Operator = ">"
	LessThanEqualTo    Operator = "<="
	GreaterThanEqualTo Operator = ">="
	NotEqual           Operator = "!="
	NotLike            Operator = "NOT LIKE"
	NotIn              Operator = "NOT IN"
	Like               Operator = "LIKE"
	In                 Operator = "IN"
)

type Logic string

const (
	And Logic = "AND"
	Or  Logic = "OR"
)

type Filter interface {
	girraph.Tree[Filter]

	addChild() Filter
	addCondition(logic Logic, condition *condition) Filter
	setLogic(logic Logic) Filter
	getLogic() Logic
	setCondition(condition *condition) Filter

	Generate() (string, []any)

	And(field string, op Operator, value any) Filter
	AndEq(field string, value any) Filter
	AndLike(field string, value any) Filter
	AndIn(field string, value any) Filter
	AndLt(field string, value any) Filter
	AndLtEq(field string, value any) Filter
	AndGt(field string, value any) Filter
	AndGtEq(field string, value any) Filter
	AndNotEq(field string, value any) Filter
	AndNotLike(field string, value any) Filter
	AndNotIn(field string, value any) Filter

	Or(field string, op Operator, value any) Filter
	OrEq(field string, value any) Filter
	OrLike(field string, value any) Filter
	OrIn(field string, value any) Filter
	OrLt(field string, value any) Filter
	OrLtEq(field string, value any) Filter
	OrGt(field string, value any) Filter
	OrGtEq(field string, value any) Filter
	OrNotEq(field string, value any) Filter
	OrNotLike(field string, value any) Filter
	OrNotIn(field string, value any) Filter

	AndGroup(conditions Filter) Filter
	OrGroup(conditions Filter) Filter
}

func NewFilter() Filter {
	result := &filter{
		TreeNode: girraph.MakeTreeNode[Filter](),
	}
	result.SetMeta(result)
	return result
}

func Where(field string, op Operator, value any) Filter {
	return NewFilter().And(field, op, value)
}

func WhereEq(field string, value any) Filter {
	return NewFilter().AndEq(field, value)
}

func WhereLike(field string, value any) Filter {
	return NewFilter().AndLike(field, value)
}

func WhereIn(field string, value any) Filter {
	return NewFilter().AndIn(field, value)
}

func WhereLt(field string, value any) Filter {
	return NewFilter().AndLt(field, value)
}

func WhereLtEq(field string, value any) Filter {
	return NewFilter().AndLtEq(field, value)
}

func WhereGt(field string, value any) Filter {
	return NewFilter().AndGt(field, value)
}

func WhereGtEq(field string, value any) Filter {
	return NewFilter().AndGtEq(field, value)
}

func WhereNotEq(field string, value any) Filter {
	return NewFilter().AndNotEq(field, value)
}

func WhereNotLike(field string, value any) Filter {
	return NewFilter().AndNotLike(field, value)
}

func WhereNotIn(field string, value any) Filter {
	return NewFilter().AndNotIn(field, value)
}

type filter struct {
	*girraph.TreeNode[Filter]
	condition *condition
	logic     Logic
}

func (f *filter) addChild() Filter {
	child := NewFilter()
	child.SetParent(f)
	f.AddChild(child)
	return child
}

func (f *filter) addCondition(logic Logic, condition *condition) Filter {
	child := f.addChild().setLogic(logic).setCondition(condition)
	return child
}

func (f *filter) setLogic(logic Logic) Filter {
	f.logic = logic
	return f
}

func (f *filter) getLogic() Logic {
	return f.logic
}

func (f *filter) setCondition(condition *condition) Filter {
	f.condition = condition
	return f
}

// And conditions

func (f *filter) And(field string, op Operator, value any) Filter {
	f.addCondition(And, newCondition(field, op, value))
	return f
}

func (f *filter) AndEq(field string, value any) Filter {
	return f.And(field, Equals, value)
}

func (f *filter) AndLike(field string, value any) Filter {
	return f.And(field, Like, value)
}

func (f *filter) AndIn(field string, value any) Filter {
	return f.And(field, In, value)
}

func (f *filter) AndLt(field string, value any) Filter {
	return f.And(field, LessThan, value)
}

func (f *filter) AndLtEq(field string, value any) Filter {
	return f.And(field, LessThanEqualTo, value)
}

func (f *filter) AndGt(field string, value any) Filter {
	return f.And(field, GreaterThan, value)
}

func (f *filter) AndGtEq(field string, value any) Filter {
	return f.And(field, GreaterThanEqualTo, value)
}

func (f *filter) AndNotEq(field string, value any) Filter {
	return f.And(field, NotEqual, value)
}

func (f *filter) AndNotLike(field string, value any) Filter {
	return f.And(field, NotLike, value)
}

func (f *filter) AndNotIn(field string, value any) Filter {
	return f.And(field, NotIn, value)
}

// Or conditions

func (f *filter) Or(field string, op Operator, value any) Filter {
	f.addCondition(Or, newCondition(field, op, value))
	return f
}

func (f *filter) OrEq(field string, value any) Filter {
	return f.Or(field, Equals, value)
}

func (f *filter) OrLike(field string, value any) Filter {
	return f.Or(field, Like, value)
}

func (f *filter) OrIn(field string, value any) Filter {
	return f.Or(field, In, value)
}

func (f *filter) OrLt(field string, value any) Filter {
	return f.Or(field, LessThan, value)
}

func (f *filter) OrLtEq(field string, value any) Filter {
	return f.Or(field, LessThanEqualTo, value)
}

func (f *filter) OrGt(field string, value any) Filter {
	return f.Or(field, GreaterThan, value)
}

func (f *filter) OrGtEq(field string, value any) Filter {
	return f.Or(field, GreaterThanEqualTo, value)
}

func (f *filter) OrNotEq(field string, value any) Filter {
	return f.Or(field, NotEqual, value)
}

func (f *filter) OrNotLike(field string, value any) Filter {
	return f.Or(field, NotLike, value)
}

func (f *filter) OrNotIn(field string, value any) Filter {
	return f.Or(field, NotIn, value)
}

// Condition groups

func (f *filter) AndGroup(conditions Filter) Filter {
	conditions.setLogic(And).SetParent(f)
	f.AddChild(conditions)
	return f
}

func (f *filter) OrGroup(conditions Filter) Filter {
	conditions.setLogic(Or).SetParent(f)
	f.AddChild(conditions)
	return f
}

func (f *filter) Generate() (string, []any) {
	children := f.GetChildren()
	hasChildren := len(children) > 0

	result := ""
	var params []any

	if f.condition != nil {
		conditionBody, conditionParams := f.condition.Generate()
		result += conditionBody
		params = append(params, conditionParams...)
	} else if hasChildren {
		result += "("
		first := true
		for _, child := range children {
			c := child.GetMeta()
			applied, ps := c.Generate()

			if first {
				result += applied
			} else {
				result += fmt.Sprintf(" %s %s", c.getLogic(), applied)
			}

			params = append(params, ps...)
			first = false
		}
		result += ")"
	}

	return strings.TrimSpace(result), params
}

func newCondition(f string, op Operator, value any) *condition {
	return &condition{
		field:    f,
		operator: op,
		value:    value,
	}
}

type condition struct {
	field    string
	operator Operator
	value    any
}

func (c *condition) Generate() (string, []any) {
	params := []any{c.value}
	if c.operator == In || c.operator == NotIn {
		return fmt.Sprintf("%s %s (?)", c.field, c.operator), params
	}
	return fmt.Sprintf("%s %s ?", c.field, c.operator), params
}
