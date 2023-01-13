package filter

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Filter interface {
	String() string
	WhereField(field string, op Operator, value any) Filter
	OrWhereField(field string, op Operator, value any) Filter
	Apply() (string, []any, error)
}

func NewFilter() Filter {
	return &node{
		Logic: logicAnd,
	}
}

type node struct {
	Condition Condition
	Logic     logic
	Children  []*node
	Parent    *node
}

func (n *node) String() string {
	if n == nil {
		return ""
	}
	aString := ""
	if n.Condition != nil {
		aString = n.Condition.String()
	}
	lString := string(n.Logic)
	var children []string
	for _, c := range n.Children {
		children = append(children, c.String())
	}
	cString := strings.Join(children, "\n ")
	pString := n.Parent.String()
	return fmt.Sprintf("filter: %v\n logic: %v\n children: %v\n parent: %v\n", aString, lString, cString, pString)
}

func (n *node) Group(l logic) Filter {
	node := &node{
		Parent: n,
		Logic:  l,
	}
	n.Children = append(n.Children, node)
	return node
}

func (n *node) where(c Condition, l logic) *node {
	n.Children = append(n.Children, &node{
		Condition: c,
		Logic:     l,
	})
	return n
}

func (n *node) WhereField(field string, op Operator, value any) Filter {
	filter := NewCondition(field, op, value)
	return n.where(filter, logicAnd)
}

func (n *node) OrWhereField(field string, op Operator, value any) Filter {
	filter := NewCondition(field, op, value)
	return n.where(filter, logicOr)
}

func (n *node) Apply() (string, []any, error) {

	q := ""
	var params []any
	if n.Condition != nil && len(n.Children) > 0 {
		return "", nil, errors.New("cannot have both a direct filter and children")
	}

	if n.Condition != nil {
		applied, ps, err := n.Condition.Apply()
		if err != nil {
			return "", nil, errors.New("failed to apply filter")
		}

		q += applied + " "
		params = append(params, ps...)
	}

	if len(n.Children) > 0 {
		q += "("
		first := true
		for _, c := range n.Children {
			applied, ps, err := c.Apply()
			if err != nil {
				return "", nil, errors.Wrap(err, "failed to apply filter children")
			}

			if first {
				q += applied
			} else {
				q += fmt.Sprintf("%s %s", c.Logic, applied)
			}

			params = append(params, ps...)
			first = false
		}
		q += ") "
	}

	return q, params, nil
}
