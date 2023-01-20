package query

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Order(t *testing.T) {
	testCases := []testCase[*Order, string]{
		{input: NewOrder().By("a"), expected: "a ASC"},
		{input: NewOrder().By("a", Ascending), expected: "a ASC"},
		{input: NewOrder().By("a", Descending), expected: "a DESC"},
		{input: NewOrder().By("a").By("b"), expected: "a ASC, b ASC"},
		{input: NewOrder().By("a").By("b", Ascending), expected: "a ASC, b ASC"},
		{input: NewOrder().By("a").By("b", Descending), expected: "a ASC, b DESC"},
	}
	for _, testCase := range testCases {
		result := testCase.input.Generate()
		assert.Equal(t, testCase.expected, result)
	}
}

func Test_DirectionFromString(t *testing.T) {
	testCases := []testCase[string, string]{
		{input: "asc", expected: string(Ascending)},
		{input: "ASC", expected: string(Ascending)},
		{input: "desc", expected: string(Descending)},
		{input: "DESC", expected: string(Descending)},
	}
	for _, testCase := range testCases {
		result, err := DirectionFromString(testCase.input)
		require.Nil(t, err)
		assert.Equal(t, testCase.expected, string(result))
	}
	result, err := DirectionFromString("")
	require.NotNil(t, err)
	assert.Equal(t, "", string(result))
}

func Test_Order_NewOrderFromUrl(t *testing.T) {
	testCases := []testCase[url.Values, string]{
		{input: mustMakeQuery(""), expected: ""},
		{input: mustMakeQuery("sort=a"), expected: "a ASC"},
		{input: mustMakeQuery("sort=a&sortDir=asc"), expected: "a ASC"},
		{input: mustMakeQuery("sort=a&sortDir=desc"), expected: "a DESC"},
	}
	for _, testCase := range testCases {
		result := NewOrderFromUrl(testCase.input).Generate()
		assert.Equal(t, testCase.expected, result)
	}
}

func Test_Order_ApplyToUrl(t *testing.T) {
	input := mustMakeQuery("")
	NewOrder().By("a", Descending).ApplyToUrl(input)
	assert.Equal(t, "a", input.Get(sortKey))
	assert.Equal(t, "DESC", input.Get(sortDirKey))
}
