package query

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

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
	Like               Operator = "LIKE"
	In                 Operator = "IN"
)

type Logic string

const (
	And Logic = "AND"
	Or  Logic = "OR"
	Not Logic = "NOT"
)

type Filter interface {
	girraph.Tree[Filter]
	fmt.Stringer

	SetLogic(logic Logic) Filter
	GetLogic() Logic
	SetCondition(condition Condition) Filter
	GetCondition() Condition
	Where(field string, op Operator, value any) Filter
	And(field string, op Operator, value any) Filter
	Or(field string, op Operator, value any) Filter
	AndGroup(conditions Filter) Filter
	OrGroup(conditions Filter) Filter
	NotGroup(conditions Filter) Filter
	Generate() (string, []any, error)

	WhereEq(field string, value any) Filter
	WhereLike(field string, value any) Filter
	WhereIn(field string, value any) Filter
	WhereLt(field string, value any) Filter
	WhereLtEq(field string, value any) Filter
	WhereGt(field string, value any) Filter
	WhereGtEq(field string, value any) Filter
	WhereNotEq(field string, value any) Filter

	OrWhereEq(field string, value any) Filter
	OrWhereLike(field string, value any) Filter
	OrWhereIn(field string, value any) Filter
	OrWhereLt(field string, value any) Filter
	OrWhereLtEq(field string, value any) Filter
	OrWhereGt(field string, value any) Filter
	OrWhereGtEq(field string, value any) Filter
	OrWhereNotEq(field string, value any) Filter
}

// TODO: add this to girraph!
func makeTreeTemp[T any]() *girraph.TreeNode[T] {
	return &girraph.TreeNode[T]{
		ID:       uuid.New().String(),
		Children: []girraph.Tree[T]{},
	}
}

func NewFilter() Filter {
	result := &filter{
		TreeNode: makeTreeTemp[Filter](),
	}
	result.SetMeta(result)
	return result
}

func Where(field string, op Operator, value any) Filter {
	return NewFilter().And(field, op, value)
}

type filter struct {
	*girraph.TreeNode[Filter]
	condition Condition
	logic     Logic
}

func (f *filter) String() string {
	if f == nil {
		return ""
	}

	conditionString := ""
	if f.condition != nil {
		conditionString = f.condition.String()
	}

	logicString := string(f.logic)

	var children []string
	for _, child := range f.GetChildren() {
		c := child.GetMeta()
		children = append(children, c.String())
	}
	childrenString := strings.Join(children, "\n ")

	parentString := ""
	parent := f.GetParent()
	if parent != nil {
		parentString = parent.GetMeta().String()
	}

	return fmt.Sprintf("condition: %v\n logic: %v\n children: %v\n parent: %v\n", conditionString, logicString, childrenString, parentString)
}

func (f *filter) addChild() Filter {
	child := NewFilter()
	child.SetParent(f)
	f.AddChild(child)
	return child
}

func (f *filter) addCondition(logic Logic, condition Condition) Filter {
	child := f.addChild().SetLogic(logic).SetCondition(condition)
	return child
}

func (f *filter) SetLogic(logic Logic) Filter {
	f.logic = logic
	return f
}

func (f *filter) GetLogic() Logic {
	return f.logic
}

func (f *filter) SetCondition(condition Condition) Filter {
	f.condition = condition
	return f
}

func (f *filter) GetCondition() Condition {
	return f.condition
}

// And conditions

func (f *filter) And(field string, op Operator, value any) Filter {
	f.addCondition(And, NewCondition(field, op, value))
	return f
}

func (f *filter) Where(field string, op Operator, value any) Filter {
	return f.And(field, op, value)
}

func (f *filter) WhereEq(field string, value any) Filter {
	return f.And(field, Equals, value)
}

func (f *filter) WhereLike(field string, value any) Filter {
	return f.And(field, Like, value)
}

func (f *filter) WhereIn(field string, value any) Filter {
	return f.And(field, In, value)
}

func (f *filter) WhereLt(field string, value any) Filter {
	return f.And(field, LessThan, value)
}

func (f *filter) WhereLtEq(field string, value any) Filter {
	return f.And(field, LessThanEqualTo, value)
}

func (f *filter) WhereGt(field string, value any) Filter {
	return f.And(field, GreaterThan, value)
}

func (f *filter) WhereGtEq(field string, value any) Filter {
	return f.And(field, GreaterThanEqualTo, value)
}

func (f *filter) WhereNotEq(field string, value any) Filter {
	return f.And(field, NotEqual, value)
}

// Or conditions

func (f *filter) Or(field string, op Operator, value any) Filter {
	f.addCondition(Or, NewCondition(field, op, value))
	return f
}

func (f *filter) OrWhereEq(field string, value any) Filter {
	return f.Or(field, Equals, value)
}

func (f *filter) OrWhereLike(field string, value any) Filter {
	return f.Or(field, Like, value)
}

func (f *filter) OrWhereIn(field string, value any) Filter {
	return f.Or(field, In, value)
}

func (f *filter) OrWhereLt(field string, value any) Filter {
	return f.Or(field, LessThan, value)
}

func (f *filter) OrWhereLtEq(field string, value any) Filter {
	return f.Or(field, LessThanEqualTo, value)
}

func (f *filter) OrWhereGt(field string, value any) Filter {
	return f.Or(field, GreaterThan, value)
}

func (f *filter) OrWhereGtEq(field string, value any) Filter {
	return f.Or(field, GreaterThanEqualTo, value)
}

func (f *filter) OrWhereNotEq(field string, value any) Filter {
	return f.Or(field, NotEqual, value)
}

// Condition groups

func (f *filter) AndGroup(conditions Filter) Filter {
	conditions.SetLogic(And).SetParent(f)
	f.AddChild(conditions)
	return f
}

func (f *filter) OrGroup(conditions Filter) Filter {
	conditions.SetLogic(Or).SetParent(f)
	f.AddChild(conditions)
	return f
}

func (f *filter) NotGroup(conditions Filter) Filter {
	conditions.SetLogic(Not).SetParent(f)
	f.AddChild(conditions)
	return f
}

func (f *filter) Generate() (string, []any, error) {
	children := f.GetChildren()
	hasChildren := len(children) > 0

	result := ""
	var params []any
	if f.condition != nil && hasChildren {
		return "", nil, errors.New("cannot have both a direct filter and children")
	}

	if f.condition != nil {
		conditionBody, conditionParams := f.condition.Generate()
		result += conditionBody
		params = append(params, conditionParams...)
	}

	if hasChildren {
		result += "("
		first := true
		for _, child := range children {
			c := child.GetMeta()
			applied, ps, err := c.Generate()
			if err != nil {
				return "", nil, errors.Wrap(err, "failed to apply filter children")
			}

			if first {
				result += applied
			} else {
				result += fmt.Sprintf(" %s %s", c.GetLogic(), applied)
			}

			params = append(params, ps...)
			first = false
		}
		result += ")"
	}

	return strings.TrimSpace(result), params, nil
}
