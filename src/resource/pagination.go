package resource

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
)

const (
	pageDefault       int64 = 1
	pageSizeDefault   int64 = 50
	totalDefault      int64 = 0
	totalPagesDefault int64 = 0

	pageKey       = "page"
	pageSizeKey   = "pageSize"
	totalKey      = "total"
	totalPagesKey = "totalPages"
)

type Pagination struct {
	Page       int64 `json:"page,omitempty"`       // the json name must match pageKey, defined above
	PageSize   int64 `json:"pageSize,omitempty"`   // the json name must match pageSizeKey, defined above
	Total      int64 `json:"total,omitempty"`      // the json name must match totalKey, defined above
	TotalPages int64 `json:"totalPages,omitempty"` // the json name must match totalPagesKey, defined above
}

func NewPagination() Pagination {
	return Pagination{
		Page:       pageDefault,
		PageSize:   pageSizeDefault,
		Total:      totalDefault,
		TotalPages: totalPagesDefault,
	}
}

func NewPaginationFromUrl(q url.Values) Pagination {
	result := NewPagination()

	if q == nil {
		return result
	}

	page, err := strconv.ParseInt(q.Get(pageKey), 10, 32)
	if err == nil {
		result.Page = page
	}

	pageSize, err := strconv.ParseInt(q.Get(pageSizeKey), 10, 32)
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
	return fmt.Sprintf("page: %v\n pageSize: %v\n total: %v\n totalPages: %v\n", p.Page, p.PageSize, p.Total, p.TotalPages)
}

func (p Pagination) GetLimit() int {
	return int(p.PageSize)
}

func (p Pagination) GetOffset() int {
	return (int(p.Page) - 1) * int(p.PageSize)
}

func (p Pagination) AddToQuery(result url.Values) {
	if p.Page != 0 {
		result.Set(pageKey, strconv.FormatInt(p.Page, 10))
	}
	if p.PageSize != 0 {
		result.Set(pageSizeKey, strconv.FormatInt(p.PageSize, 10))
	}
	if p.Total != 0 {
		result.Set(totalKey, strconv.FormatInt(p.Total, 10))
	}
	if p.TotalPages != 0 {
		result.Set(totalPagesKey, strconv.FormatInt(p.TotalPages, 10))
	}
}

func (p Pagination) First() Pagination {
	return Pagination{
		Page:       1,
		PageSize:   p.PageSize,
		Total:      p.Total,
		TotalPages: p.TotalPages,
	}
}

func (p Pagination) Previous() Pagination {
	page := p.Page
	if p.Page > 1 {
		page -= 1
	}
	return Pagination{
		Page:       page,
		PageSize:   p.PageSize,
		Total:      p.Total,
		TotalPages: p.TotalPages,
	}
}

func (p Pagination) Next() Pagination {
	page := p.Page
	if p.TotalPages > p.Page {
		page += 1
	}
	return Pagination{
		Page:       page,
		PageSize:   p.PageSize,
		Total:      p.Total,
		TotalPages: p.TotalPages,
	}
}

func (p Pagination) Last() Pagination {
	return Pagination{
		Page:       p.TotalPages,
		PageSize:   p.PageSize,
		Total:      p.Total,
		TotalPages: p.TotalPages,
	}
}
