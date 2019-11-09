package filter

type Operator string

const (
	OpEq    Operator = "="
	OpLt    Operator = "<"
	OpGt    Operator = ">"
	OpLtEq  Operator = "<="
	OpGtEq  Operator = ">="
	OpNotEq Operator = "<>"
	OpLike  Operator = "LIKE"
	OpIn    Operator = "IN"
)
