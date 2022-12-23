package query

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultTotalRecords     uint = 0
	defaultTotalPages       uint = 0
	defaultRecordsPerPage   uint = 50
	defaultCurrentPageIndex uint = 1

	currentPageKey = "page"
	perPageKey     = "page_size"
)

type Pagination struct {
	Page       uint `json:"page"`      // the json name must match currentPageKey, defined above
	PageSize   uint `json:"page_size"` // the json name must match perPageKey, defined above
	Total      uint `json:"total"`
	TotalPages uint `json:"total_pages"`

	FirstPage    string `json:"first_page"`
	PreviousPage string `json:"previous_page"`
	NextPage     string `json:"next_page"`
	LastPage     string `json:"last_page"`
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

	page, err := strconv.ParseUint(c.Query(currentPageKey), 10, 32)
	if err == nil {
		result.Page = uint(page)
	}

	pageSize, err := strconv.ParseUint(c.Query(perPageKey), 10, 32)
	if err == nil {
		result.PageSize = uint(pageSize)
	}

	return result
}

func NewPaginationFromValues(page, size, total uint) Pagination {
	return Pagination{
		Page:     page,
		PageSize: size,
		Total:    total,
		// When performing integer division, Go rounds quotients down to the nearest integer.
		// We want to round up in this case, so we need to divide floats.
		TotalPages: uint(math.Ceil(float64(total) / float64(size))),
	}
}

func (p Pagination) String() string {
	return fmt.Sprintf("page: %v\n page_size: %v\n total: %v\n", p.Page, p.PageSize, p.Total)
}
