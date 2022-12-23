package query

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultTotalRecords     int64 = 0
	defaultTotalPages       int64 = 0
	defaultRecordsPerPage   int64 = 50
	defaultCurrentPageIndex int64 = 1

	currentPageKey = "page"
	perPageKey     = "page_size"
)

type Pagination struct {
	Page       int64 `json:"page,omitempty"`      // the json name must match currentPageKey, defined above
	PageSize   int64 `json:"page_size,omitempty"` // the json name must match perPageKey, defined above
	Total      int64 `json:"total,omitempty"`
	TotalPages int64 `json:"total_pages,omitempty"`

	FirstPage    string `json:"first_page,omitempty"`
	PreviousPage string `json:"previous_page,omitempty"`
	NextPage     string `json:"next_page,omitempty"`
	LastPage     string `json:"last_page,omitempty"`
}

func NewPaginationFromGin(c *gin.Context) Pagination {
	result := Pagination{
		Page:       defaultCurrentPageIndex,
		PageSize:   defaultRecordsPerPage,
		Total:      defaultTotalRecords,
		TotalPages: defaultTotalPages,
	}

	if c == nil {
		return result
	}

	page, err := strconv.ParseInt(c.Query(currentPageKey), 10, 32)
	if err == nil {
		result.Page = page
	}

	pageSize, err := strconv.ParseInt(c.Query(perPageKey), 10, 32)
	if err == nil {
		result.PageSize = pageSize
	}

	return result
}

func NewPaginationFromValues(page, size, total int64) Pagination {
	return Pagination{
		Page:     page,
		PageSize: size,
		Total:    total,
		// When performing integer division, Go rounds quotients down to the nearest integer.
		// We want to round up in this case, so we need to divide floats.
		TotalPages: int64(math.Ceil(float64(total) / float64(size))),
	}
}

func (p Pagination) String() string {
	return fmt.Sprintf("page: %v\n page_size: %v\n total: %v\n", p.Page, p.PageSize, p.Total)
}
