package query

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type Direction string

const (
	Ascending      Direction = "ASC"
	Descending     Direction = "DESC"
	defaultSortDir           = Ascending

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

func NewOrder() *Order {
	return &Order{
		sort: []*sort{},
	}
}

func NewOrderFromUrl(q url.Values) *Order {
	result := NewOrder()

	sortField := q.Get(sortKey)
	if sortField == "" {
		return result
	}

	sortDir, err := DirectionFromString(q.Get(sortDirKey))
	if err != nil {
		sortDir = defaultSortDir
	}

	return result.By(sortField, sortDir)
}

type Order struct {
	sort []*sort
}

func (o *Order) By(field string, dir ...Direction) *Order {
	d := defaultSortDir
	if len(dir) > 0 {
		d = dir[0]
	}
	o.sort = append(o.sort, newSort().By(field).Dir(d))
	return o
}

func (o *Order) Generate() string {
	if len(o.sort) == 0 {
		return ""
	}
	var result []string
	for _, s := range o.sort {
		result = append(result, s.Generate())
	}
	return strings.Join(result, ", ")
}

func (o *Order) ApplyToUrl(q url.Values) {
	// TODO: support sorting by multiple fields
	if len(o.sort) > 0 {
		s := o.sort[0]
		q.Set(sortKey, s.field)
		q.Set(sortDirKey, string(s.direction))
	}
}

func newSort() *sort {
	return &sort{
		field:     "",
		direction: "",
	}
}

type sort struct {
	field     string
	direction Direction
}

func (s *sort) Generate() string {
	return fmt.Sprintf("%s %s", s.field, strings.ToUpper(string(s.direction)))
}

func (s *sort) By(field string) *sort {
	s.field = field
	return s
}

func (s *sort) Dir(dir Direction) *sort {
	s.direction = dir
	return s
}
