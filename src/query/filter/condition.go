package filter

import "fmt"

type Condition interface {
	String() string
	Apply() (string, []any, error)
}

func NewCondition(f string, op Operator, value any) Condition {
	return &fieldCondition{
		Field: f,
		Op:    op,
		Value: value,
	}
}

type fieldCondition struct {
	Field string
	Op    Operator
	Value any
}

func (c *fieldCondition) String() string {
	if c == nil {
		return ""
	}
	oString := string(c.Op)
	vString, ok := c.Value.(string)
	if !ok {
		b, ok := c.Value.([]byte)
		if !ok {
			b = []byte{}
		}
		vString = string(b)
	}
	return fmt.Sprintf("field: %v\n op: %v\n value: %v\n", c.Field, oString, vString)
}

func (c *fieldCondition) Apply() (string, []any, error) {
	if c.Op == OpIn {
		return fmt.Sprintf("%s %s (?)", c.Field, c.Op), []any{c.Value}, nil
	}
	return fmt.Sprintf("%s %s ?", c.Field, c.Op), []any{c.Value}, nil
}
