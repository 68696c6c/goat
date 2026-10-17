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
	sortDirKey = "dir"
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
		sort: []Sort{},
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
	sort []Sort
}

func (o *Order) By(field string, dir ...Direction) *Order {
	d := defaultSortDir
	if len(dir) > 0 {
		d = dir[0]
	}
	o.sort = append(o.sort, Sort{
		Field:     field,
		Direction: d,
	})
	return o
}

func (o *Order) Generate() (string, []string) {
	if len(o.sort) == 0 {
		return "", []string{}
	}
	var result []string
	var params []string
	for _, s := range o.sort {
		template, field := s.Generate()
		result = append(result, template)
		params = append(params, field)
	}
	return strings.Join(result, ", "), params
}

func (o *Order) GetSorts() []Sort {
	return o.sort
}

func (o *Order) ApplyToUrl(q url.Values) {
	// TODO: support sorting by multiple fields
	if len(o.sort) > 0 {
		s := o.sort[0]
		q.Set(sortKey, s.Field)
		q.Set(sortDirKey, string(s.Direction))
	}
}

type Sort struct {
	Field     string
	Direction Direction
}

func (s *Sort) Generate() (string, string) {
	return fmt.Sprintf("? %s", strings.ToUpper(string(s.Direction))), s.Field
}
