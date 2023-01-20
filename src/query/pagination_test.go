package query

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Pagination_NewPaginationFromUrl(t *testing.T) {
	testCases := []testCase[url.Values, Pagination]{
		{
			input: nil,
			expected: Pagination{
				Page:       0,
				PageSize:   0,
				Total:      0,
				TotalPages: 0,
			},
		},
		{
			input: mustMakeQuery("page=10&pageSize=5"),
			expected: Pagination{
				Page:       10,
				PageSize:   5,
				Total:      0,
				TotalPages: 0,
			},
		},
		{
			input: mustMakeQuery("sort=a&sortDir=desc&page=19&pageSize=6"),
			expected: Pagination{
				Page:       19,
				PageSize:   6,
				Total:      0,
				TotalPages: 0,
			},
		},
		{
			input: mustMakeQuery("page=42&pageSize=17&total=257"),
			expected: Pagination{
				Page:       42,
				PageSize:   17,
				Total:      257,
				TotalPages: 16,
			},
		},
	}
	for _, testCase := range testCases {
		result := NewPaginationFromUrl(testCase.input)
		assert.Equal(t, testCase.expected.Page, result.GetPage())
		assert.Equal(t, testCase.expected.PageSize, result.GetPageSize())
		assert.Equal(t, testCase.expected.Total, result.GetTotal())
		assert.Equal(t, testCase.expected.TotalPages, result.GetTotalPages())
	}
}

func Test_Pagination_ApplyToUrl(t *testing.T) {
	testCases := []testCase[*Pagination, Pagination]{
		{
			input: NewPagination(),
			expected: Pagination{
				Page:       0,
				PageSize:   0,
				Total:      0,
				TotalPages: 0,
			},
		},
		{
			input: NewPagination().SetPage(10).SetPageSize(5).SetTotal(100),
			expected: Pagination{
				Page:       10,
				PageSize:   5,
				Total:      100,
				TotalPages: 20,
			},
		},
		{
			input: NewPagination().SetPage(13).SetPageSize(8).SetTotal(491),
			expected: Pagination{
				Page:       13,
				PageSize:   8,
				Total:      491,
				TotalPages: 62,
			},
		},
	}
	for _, testCase := range testCases {
		subject := mustMakeQuery("")
		testCase.input.ApplyToUrl(subject)
		assert.Equal(t, strconv.Itoa(testCase.expected.Page), subject.Get(pageKey))
		assert.Equal(t, strconv.Itoa(testCase.expected.PageSize), subject.Get(pageSizeKey))
		assert.Equal(t, strconv.Itoa(testCase.expected.Total), subject.Get(totalKey))
		assert.Equal(t, strconv.Itoa(testCase.expected.TotalPages), subject.Get(totalPagesKey))
	}
}

func Test_Pagination_getOffset(t *testing.T) {
	subject := NewPagination().SetPage(1).SetPageSize(5).SetTotal(64)
	assert.Equal(t, 0, subject.getOffset())

	testCases := []testCase[*Pagination, int]{
		{
			input:    NewPagination().SetPage(1).SetPageSize(5).SetTotal(64),
			expected: 0,
		},
		{
			input:    NewPagination().SetPage(2).SetPageSize(30).SetTotal(92),
			expected: 30,
		},
	}
	for _, testCase := range testCases {
		result := testCase.input.getOffset()
		assert.Equal(t, testCase.expected, result)
	}
}

func Test_Pagination_setOffset(t *testing.T) {
	testCases := []testCase[*Pagination, int]{
		{
			input:    NewPagination().SetPageSize(5).setOffset(64),
			expected: 13,
		},
		{
			input:    NewPagination().SetPageSize(29).setOffset(46),
			expected: 2,
		},
	}
	for _, testCase := range testCases {
		result := testCase.input.GetPage()
		assert.Equal(t, testCase.expected, result)
	}
}

type offsetTestCase struct {
	page   int
	size   int
	offset int
}

// When page size, page, and offset are all positive values, offset operations should be reversible.
func Test_Pagination_Offset_PositiveValues(t *testing.T) {
	testCases := []offsetTestCase{
		{
			page:   7,
			size:   23,
			offset: 138,
		},
		{
			page:   5,
			size:   15,
			offset: 60,
		},
		{
			page:   1,
			size:   10,
			offset: 0,
		},
	}
	for _, testCase := range testCases {
		resultOffset := NewPagination().SetPage(testCase.page).SetPageSize(testCase.size).getOffset()
		assert.Equal(t, testCase.offset, resultOffset)

		resultPage := NewPagination().SetPageSize(testCase.size).setOffset(testCase.offset).GetPage()
		assert.Equal(t, testCase.page, resultPage)
	}
}

// When page size or offset is 0, the operation is not reversible and the page should always be 1.
func Test_Pagination_Offset_InvalidValues(t *testing.T) {

	// If page size or page is less than 1, offset should be 0.
	assert.Equal(t, 0, NewPagination().SetPage(0).SetPageSize(2).getOffset())
	assert.Equal(t, 0, NewPagination().SetPage(14).SetPageSize(0).getOffset())
	assert.Equal(t, 0, NewPagination().SetPage(52).SetPageSize(-7).getOffset())
	assert.Equal(t, 0, NewPagination().SetPage(-46).SetPageSize(11).getOffset())

	// If page size or offset are less than 1, page should be 1.
	assert.Equal(t, 1, NewPagination().SetPageSize(0).setOffset(67).GetPage())
	assert.Equal(t, 1, NewPagination().SetPageSize(117).setOffset(0).GetPage())
	assert.Equal(t, 1, NewPagination().SetPageSize(0).setOffset(-9).GetPage())
	assert.Equal(t, 1, NewPagination().SetPageSize(-3).setOffset(88).GetPage())
}

func Test_Pagination_First(t *testing.T) {
	testCases := []testCase[*Pagination, Pagination]{
		{
			input: NewPagination().SetPage(1).SetPageSize(9).SetTotal(32),
			expected: Pagination{
				Page:       1,
				PageSize:   9,
				Total:      32,
				TotalPages: 4,
			},
		},
		{
			input: NewPagination().SetPageSize(5).setOffset(64),
			expected: Pagination{
				// Page:       13,
				Page:       1,
				PageSize:   5,
				Total:      0,
				TotalPages: 0,
			},
		},
		{
			input: NewPagination().SetPageSize(29).setOffset(46),
			expected: Pagination{
				// Page:       2,
				Page:       1,
				PageSize:   29,
				Total:      0,
				TotalPages: 0,
			},
		},
		{
			input: NewPagination().SetPage(14).SetPageSize(8).SetTotal(491),
			expected: Pagination{
				// Page:       13,
				Page:       1,
				PageSize:   8,
				Total:      491,
				TotalPages: 62,
			},
		},
	}
	for _, testCase := range testCases {
		result := testCase.input.First()
		assert.Equal(t, testCase.expected.Page, result.GetPage())
		assert.Equal(t, testCase.expected.PageSize, result.GetPageSize())
		assert.Equal(t, testCase.expected.Total, result.GetTotal())
		assert.Equal(t, testCase.expected.TotalPages, result.GetTotalPages())
	}
}

func Test_Pagination_Previous(t *testing.T) {
	testCases := []testCase[*Pagination, Pagination]{
		{
			input: NewPagination().SetPage(1).SetPageSize(9).SetTotal(32),
			expected: Pagination{
				Page:       1,
				PageSize:   9,
				Total:      32,
				TotalPages: 4,
			},
		},
		{
			input: NewPagination().SetPageSize(5).setOffset(64),
			expected: Pagination{
				Page:       12,
				PageSize:   5,
				Total:      0,
				TotalPages: 0,
			},
		},
		{
			input: NewPagination().SetPageSize(29).setOffset(46),
			expected: Pagination{
				Page:       1,
				PageSize:   29,
				Total:      0,
				TotalPages: 0,
			},
		},
		{
			input: NewPagination().SetPage(14).SetPageSize(8).SetTotal(491),
			expected: Pagination{
				Page:       13,
				PageSize:   8,
				Total:      491,
				TotalPages: 62,
			},
		},
	}
	for _, testCase := range testCases {
		result := testCase.input.Previous()
		assert.Equal(t, testCase.expected.Page, result.GetPage())
		assert.Equal(t, testCase.expected.PageSize, result.GetPageSize())
		assert.Equal(t, testCase.expected.Total, result.GetTotal())
		assert.Equal(t, testCase.expected.TotalPages, result.GetTotalPages())
	}
}

func Test_Pagination_Next(t *testing.T) {
	testCases := []testCase[*Pagination, Pagination]{
		{
			input: NewPagination().SetPage(1).SetPageSize(9).SetTotal(32),
			expected: Pagination{
				Page:       2,
				PageSize:   9,
				Total:      32,
				TotalPages: 4,
			},
		},
		{
			input: NewPagination().SetPageSize(5).setOffset(64),
			expected: Pagination{
				Page:       13,
				PageSize:   5,
				Total:      0,
				TotalPages: 0,
			},
		},
		{
			input: NewPagination().SetPageSize(29).setOffset(46),
			expected: Pagination{
				Page:       2,
				PageSize:   29,
				Total:      0,
				TotalPages: 0,
			},
		},
		{
			input: NewPagination().SetPage(14).SetPageSize(8).SetTotal(491),
			expected: Pagination{
				Page:       15,
				PageSize:   8,
				Total:      491,
				TotalPages: 62,
			},
		},
	}
	for _, testCase := range testCases {
		result := testCase.input.Next()
		assert.Equal(t, testCase.expected.Page, result.GetPage())
		assert.Equal(t, testCase.expected.PageSize, result.GetPageSize())
		assert.Equal(t, testCase.expected.Total, result.GetTotal())
		assert.Equal(t, testCase.expected.TotalPages, result.GetTotalPages())
	}
}

func Test_Pagination_Last(t *testing.T) {
	testCases := []testCase[*Pagination, Pagination]{
		{
			input: NewPagination().SetPage(1).SetPageSize(9).SetTotal(32),
			expected: Pagination{
				Page:       4,
				PageSize:   9,
				Total:      32,
				TotalPages: 4,
			},
		},
		{
			// If total is not set, we can't calculate the last page.
			input: NewPagination().SetPageSize(5).setOffset(64),
			expected: Pagination{
				Page:       13,
				PageSize:   5,
				Total:      0,
				TotalPages: 0,
			},
		},
		{
			input: NewPagination().SetPageSize(29).SetTotal(38),
			expected: Pagination{
				Page:       2,
				PageSize:   29,
				Total:      38,
				TotalPages: 2,
			},
		},
		{
			input: NewPagination().SetPage(14).SetPageSize(8).SetTotal(491),
			expected: Pagination{
				Page:       62,
				PageSize:   8,
				Total:      491,
				TotalPages: 62,
			},
		},
	}
	for _, testCase := range testCases {
		result := testCase.input.Last()
		assert.Equal(t, testCase.expected.Page, result.GetPage(), "unexpected page")
		assert.Equal(t, testCase.expected.PageSize, result.GetPageSize(), "unexpected page size")
		assert.Equal(t, testCase.expected.Total, result.GetTotal(), "unexpected total")
		assert.Equal(t, testCase.expected.TotalPages, result.GetTotalPages(), "unexpected total pages")
	}
}

func Test_ceilQuotient_ZeroDivisor(t *testing.T) {
	result := ceilQuotient(1, 0)
	assert.Equal(t, 0, result)
}
