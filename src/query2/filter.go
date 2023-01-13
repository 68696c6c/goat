package query2

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/68696c6c/girraph"
)

type Filter interface {
	girraph.Tree[Filter]
	// fmt.Stringer
	SetLogic(logic Logic) Filter
	GetLogic() Logic
	SetCondition(condition Condition) Filter
	GetCondition() Condition
	Where(field string, op Operator, value any) Filter
	And(field string, op Operator, value any) Filter
	Or(field string, op Operator, value any) Filter
	Not(field string, op Operator, value any) Filter
	// Group(logic Logic) Filter
	// GroupEnd() Filter
	AndGroup(conditions Filter) Filter
	OrGroup(conditions Filter) Filter
	NotGroup(conditions Filter) Filter
	Generate() (string, []any, error)
}

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

func (w *filter) addChild() Filter {
	child := NewFilter()
	child.SetParent(w)
	w.AddChild(child)
	return child
}

func (w *filter) addCondition(logic Logic, condition Condition) Filter {
	child := w.addChild().SetLogic(logic).SetCondition(condition)
	return child
}

func (w *filter) SetLogic(logic Logic) Filter {
	w.logic = logic
	return w
}

func (w *filter) GetLogic() Logic {
	return w.logic
}

func (w *filter) SetCondition(condition Condition) Filter {
	w.condition = condition
	return w
}

func (w *filter) GetCondition() Condition {
	return w.condition
}

func (w *filter) Where(field string, op Operator, value any) Filter {
	return w.And(field, op, value)
}

func (w *filter) And(field string, op Operator, value any) Filter {
	w.addCondition(And, NewCondition(field, op, value))
	return w
}

func (w *filter) Or(field string, op Operator, value any) Filter {
	w.addCondition(Or, NewCondition(field, op, value))
	return w
}

func (w *filter) Not(field string, op Operator, value any) Filter {
	w.addCondition(Not, NewCondition(field, op, value))
	return w
}

// func (w *filter) Group(logic Logic) Filter {
// 	group := w.addChild().SetLogic(logic)
// 	return group
// }
//
// func (w *filter) GroupEnd() Filter {
// 	return w.GetParent().GetMeta()
// }

func (w *filter) AndGroup(conditions Filter) Filter {
	conditions.SetLogic(And).SetParent(w)
	w.AddChild(conditions)
	return w
}

func (w *filter) OrGroup(conditions Filter) Filter {
	conditions.SetLogic(Or).SetParent(w)
	w.AddChild(conditions)
	return w
}

func (w *filter) NotGroup(conditions Filter) Filter {
	conditions.SetLogic(Not).SetParent(w)
	w.AddChild(conditions)
	return w
}

func (w *filter) Generate() (string, []any, error) {
	children := w.GetChildren()
	hasChildren := len(children) > 0

	result := ""
	var params []any
	if w.condition != nil && hasChildren {
		return "", nil, errors.New("cannot have both a direct filter and children")
	}

	if w.condition != nil {
		conditionBody, conditionParams := w.condition.Generate()
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
