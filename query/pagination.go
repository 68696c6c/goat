package query

import (
	"math"
	"net/url"
	"strconv"
)

const (
	pageKey       = "page"
	pageSizeKey   = "size"
	totalKey      = "total"
	totalPagesKey = "pages"
)

type Pagination struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"size,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"pages,omitempty"`
}

func NewPagination() *Pagination {
	return &Pagination{
		Page:       0,
		PageSize:   0,
		Total:      0,
		TotalPages: 0,
	}
}

func NewPaginationFromUrl(q url.Values) *Pagination {
	result := NewPagination()

	if q == nil {
		return result
	}

	page, err := strconv.Atoi(q.Get(pageKey))
	if err == nil {
		result.SetPage(page)
	}

	pageSize, err := strconv.Atoi(q.Get(pageSizeKey))
	if err == nil {
		result.SetPageSize(pageSize)
	}

	total, err := strconv.Atoi(q.Get(totalKey))
	if err == nil {
		result.SetTotal(total)
	}

	return result.setTotalPages()
}

func (p *Pagination) SetPage(page int) *Pagination {
	p.Page = page
	return p
}

func (p *Pagination) GetPage() int {
	return p.Page
}

func (p *Pagination) SetPageSize(size int) *Pagination {
	p.PageSize = size
	return p.setTotalPages()
}

func (p *Pagination) GetPageSize() int {
	return p.PageSize
}

func (p *Pagination) SetTotal(total int) *Pagination {
	p.Total = total
	return p.setTotalPages()
}

func (p *Pagination) GetTotal() int {
	return p.Total
}

// ceilQuotient divides dividend by divisor and rounds the result up to the nearest int.  If divisor is zero, zero is returned.
func ceilQuotient(dividend, divisor int) int {
	if divisor == 0 {
		return 0
	}
	// When performing integer division, Go rounds quotients down to the nearest integer.
	// In order to round up, we need to divide floats.
	return int(math.Ceil(float64(dividend) / float64(divisor)))
}

func (p *Pagination) setTotalPages() *Pagination {
	if p.PageSize > 0 {
		p.TotalPages = ceilQuotient(p.Total, p.PageSize)
	}
	return p
}

func (p *Pagination) SetProperties(totalRows, pageSize, page int) *Pagination {
	p.Total = totalRows
	p.PageSize = pageSize
	if p.PageSize > 0 {
		p.TotalPages = ceilQuotient(p.Total, p.PageSize)
	}
	if page > p.TotalPages {
		p.Page = p.TotalPages
	} else {
		p.Page = page
	}
	return p
}

func (p *Pagination) GetTotalPages() int {
	return p.TotalPages
}

func (p *Pagination) setOffset(offset int) *Pagination {
	if offset < 1 || p.PageSize < 1 {
		p.Page = 1
	} else {
		p.Page = (offset / p.PageSize) + 1
	}
	return p
}

func (p *Pagination) getOffset() int {
	// Minimum offset is 0.
	if p.Page < 1 || p.PageSize < 1 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) ApplyToUrl(q url.Values) {
	if p.Page > -1 {
		q.Set(pageKey, strconv.Itoa(p.Page))
	}
	if p.PageSize > -1 {
		q.Set(pageSizeKey, strconv.Itoa(p.PageSize))
	}
	if p.Total > -1 {
		q.Set(totalKey, strconv.Itoa(p.Total))
	}
	if p.TotalPages > -1 {
		q.Set(totalPagesKey, strconv.Itoa(p.TotalPages))
	}
}

func (p *Pagination) First() *Pagination {
	return &Pagination{
		Page:       1,
		PageSize:   p.PageSize,
		Total:      p.Total,
		TotalPages: p.TotalPages,
	}
}

func (p *Pagination) Previous() *Pagination {
	// Make sure we don't return a negative page value.
	page := p.Page
	if p.Page > 1 {
		page -= 1
	}
	return &Pagination{
		Page:       page,
		PageSize:   p.PageSize,
		Total:      p.Total,
		TotalPages: p.TotalPages,
	}
}

func (p *Pagination) Next() *Pagination {
	// Make sure we don't increment past the last page.
	page := p.Page
	if p.TotalPages > p.Page {
		page += 1
	}
	return &Pagination{
		Page:       page,
		PageSize:   p.PageSize,
		Total:      p.Total,
		TotalPages: p.TotalPages,
	}
}

func (p *Pagination) Last() *Pagination {
	// Try and calculate the total pages if it hasn't been set yet.
	if p.TotalPages == 0 {
		p.setTotalPages()
	}
	// Fallback to the current page if total pages is still 0 (This can happen if page size or total are not set).
	page := p.TotalPages
	if page == 0 {
		page = p.Page
	}
	return &Pagination{
		Page:       page,
		PageSize:   p.PageSize,
		Total:      p.Total,
		TotalPages: p.TotalPages,
	}
}
