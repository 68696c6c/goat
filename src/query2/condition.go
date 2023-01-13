package query2

import "fmt"

type Condition interface {
	fmt.Stringer
	Generate() (string, []any)
}

func NewCondition(f string, op Operator, value any) Condition {
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

func (c *condition) String() string {
	if c == nil {
		return ""
	}
	operator := string(c.operator)
	value, ok := c.value.(string)
	if !ok {
		b, ok := c.value.([]byte)
		if !ok {
			b = []byte{}
		}
		value = string(b)
	}
	return fmt.Sprintf("field: %v\n operator: %v\n value: %v\n", c.field, operator, value)
}

func (c *condition) Generate() (string, []any) {
	params := []any{c.value}
	if c.operator == In {
		return fmt.Sprintf("%s %s (?)", c.field, c.operator), params
	}
	return fmt.Sprintf("%s %s ?", c.field, c.operator), params
}
