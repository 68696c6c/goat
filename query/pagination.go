package query

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage     = 1
	defaultPageSize = 50
	defaultTotal    = 0
)

func NewPagination(c *gin.Context) Pagination {
	if c == nil {
		return Pagination{
			Page:     uint(defaultPage),
			PageSize: uint(defaultPageSize),
			Total:    defaultTotal,
		}
	}
	page, err := strconv.ParseUint(c.Query("page"), 10, 32)
	if err != nil {
		page = defaultPage
	}
	pageSize, err := strconv.ParseUint(c.Query("page_size"), 10, 32)
	if err != nil {
		pageSize = defaultPageSize
	}
	return Pagination{
		Page:     uint(page),
		PageSize: uint(pageSize),
		Total:    defaultTotal,
	}
}

type Pagination struct {
	Page     uint `json:"page"`
	PageSize uint `json:"page_size"`
	Total    uint `json:"total"`
}

func (p Pagination) String() string {
	return fmt.Sprintf("page: %v\n page_size: %v\n total: %v\n", p.Page, p.PageSize, p.Total)
}

func (p Pagination) Paginate(page, size, total uint) {
	p.Page = page
	p.PageSize = size
	p.Total = total
}
