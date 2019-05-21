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

//var ops = []Operator{OpEq, OpLt, OpGt, OpLte, OpGte, OpNeq, OpIn, OpLike}
//
//func OperatorFromString(s string) (Operator, error) {
//	for i := range ops {
//		if string(ops[i]) == s {
//			return ops[i], nil
//		}
//	}
//	return Operator(""), errors.Errorf("%s not an operator", s)
//}
