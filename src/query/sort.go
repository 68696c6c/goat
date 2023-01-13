package query

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Direction string

// TODO: should include the sort params in the links like we do with pagination params
const (
	Ascending        Direction = "ASC"
	Descending       Direction = "DESC"
	defaultSortField           = "created_at" // TODO: something more generic, like id, or maybe make this configurable?
	defaultSortDir             = Descending

	sortKey    = "sort"
	sortDirKey = "sortDir"
)

func DirectionFromString(input string) (Direction, error) {
	s := strings.ToUpper(input)
	if s == string(Ascending) || s == string(Descending) {
		return Direction(s), nil
	}
	return Direction(""), errors.Errorf("%s not a sort direction", input)
}

type Sort interface {
	fmt.Stringer
	Generate() string
	By(field string) Sort
	Asc() Sort
	Desc() Sort
	Dir(dir Direction) Sort
}

func NewSort() Sort {
	return &sort{
		field:     defaultSortField,
		direction: defaultSortDir,
	}
}

// func NewSortFromUrl() Sort {
// 	return &sort{
// 		field:     defaultSortField,
// 		direction: defaultSortDir,
// 	}
// }

type sort struct {
	field     string
	direction Direction
}

func (s *sort) String() string {
	direction := string(s.direction)
	return fmt.Sprintf("field: %v\n direction: %v\n", s.field, direction)
}

func (s *sort) Generate() string {
	return fmt.Sprintf("%s %s", s.field, s.direction)
}

func (s *sort) By(field string) Sort {
	s.field = field
	return s
}

func (s *sort) Asc() Sort {
	s.direction = Ascending
	return s
}

func (s *sort) Desc() Sort {
	s.direction = Descending
	return s
}

func (s *sort) Dir(dir Direction) Sort {
	s.direction = dir
	return s
}
