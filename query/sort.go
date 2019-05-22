package query

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type SortDir string

const (
	SortAsc          SortDir = "asc"
	SortDesc         SortDir = "desc"
	defaultSortField         = "created_at"
	defaultSortDir           = SortDesc
)

func SortDirFromString(s string) (SortDir, error) {
	if s == string(SortAsc) ||
		s == string(SortDesc) {
		return SortDir(s), nil
	}
	return SortDir(""), errors.Errorf("%s not a sort dir", s)
}

type Sort struct {
	Field string
	Dir   SortDir
}

func NewSort(c *gin.Context) []Sort {
	if c == nil {
		return []Sort{{Field: defaultSortField, Dir: defaultSortDir}}
	}
	field := c.Query("sort")
	if field == "" {
		field = defaultSortField
	}
	dir, err := SortDirFromString(c.Query("sort_dir"))
	if err != nil {
		dir = defaultSortDir
	}
	return []Sort{{Field: field, Dir: dir}}
}

func (s Sort) String() string {
	dString := string(s.Dir)
	return fmt.Sprintf("field: %v\n dir: %v\n", s.Field, dString)
}
