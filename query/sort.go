package query

import "fmt"

type SortDir string

const (
	SortAsc  SortDir = "ASC"
	SortDesc SortDir = "DESC"
)

type Sort struct {
	Field string
	Dir   SortDir
}

func (s Sort) String() string {
	dString := string(s.Dir)
	return fmt.Sprintf("field: %v\n dir: %v\n", s.Field, dString)
}
