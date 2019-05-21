package query

import "fmt"

const (
	defaultPage     = 1
	defaultPageSize = 50
	defaultTotal    = 0
)

func NewPagination() *Pagination {
	return &Pagination{
		Page:     defaultPage,
		PageSize: defaultPageSize,
		Total:    defaultTotal,
	}
}

type Pagination struct {
	Page     uint `json:"page"`
	PageSize uint `json:"page_size"`
	Total    uint `json:"total"`
}

func (p *Pagination) String() string {
	return fmt.Sprintf("page: %v\n page_size: %v\n total: %v\n", p.Page, p.PageSize, p.Total)
}

func (p *Pagination) Paginate(page, size, total uint) {
	p.Page = page
	p.PageSize = size
	p.Total = total
}
