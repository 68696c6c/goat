package query

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type orderGenerateResult struct {
	template string
	params   []string
}

func Test_Order(t *testing.T) {
	testCases := []testCase[*Order, orderGenerateResult]{
		{
			input: NewOrder().By("a"),
			expected: orderGenerateResult{
				template: "? ASC",
				params:   []string{"a"},
			},
		},
		{
			input: NewOrder().By("a", Ascending),
			expected: orderGenerateResult{
				template: "? ASC",
				params:   []string{"a"},
			},
		},
		{
			input: NewOrder().By("a", Descending),
			expected: orderGenerateResult{
				template: "? DESC",
				params:   []string{"a"},
			},
		},
		{
			input: NewOrder().By("a").By("b"),
			// expected: "a ASC, b ASC"
			expected: orderGenerateResult{
				template: "? ASC, ? ASC",
				params:   []string{"a", "b"},
			},
		},
		{
			input: NewOrder().By("a").By("b", Ascending),
			// expected: "a ASC, b ASC"
			expected: orderGenerateResult{
				template: "? ASC, ? ASC",
				params:   []string{"a", "b"},
			},
		},
		{
			input: NewOrder().By("a").By("b", Descending),
			// expected: "a ASC, b DESC"
			expected: orderGenerateResult{
				template: "? ASC, ? DESC",
				params:   []string{"a", "b"},
			},
		},
	}
	for _, testCase := range testCases {
		template, params := testCase.input.Generate()
		assert.Equal(t, testCase.expected.template, template)
		assert.Equal(t, testCase.expected.params, params)
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
	testCases := []testCase[url.Values, orderGenerateResult]{
		{
			input: mustMakeQuery(""),
			// expected: "",
			expected: orderGenerateResult{
				template: "",
				params:   []string{},
			},
		},
		{
			input: mustMakeQuery("sort=a"),
			// expected: "a ASC",
			expected: orderGenerateResult{
				template: "? ASC",
				params:   []string{"a"},
			},
		},
		{
			input: mustMakeQuery("sort=a&dir=asc"),
			// expected: "a ASC",
			expected: orderGenerateResult{
				template: "? ASC",
				params:   []string{"a"},
			},
		},
		{
			input: mustMakeQuery("sort=a&dir=desc"),
			// expected: "a DESC",
			expected: orderGenerateResult{
				template: "? DESC",
				params:   []string{"a"},
			},
		},
	}
	for _, testCase := range testCases {
		template, params := NewOrderFromUrl(testCase.input).Generate()
		assert.Equal(t, testCase.expected.template, template)
		assert.Equal(t, testCase.expected.params, params)
	}
}

func Test_Order_ApplyToUrl(t *testing.T) {
	input := mustMakeQuery("")
	NewOrder().By("a", Descending).ApplyToUrl(input)
	assert.Equal(t, "a", input.Get(sortKey))
	assert.Equal(t, "DESC", input.Get(sortDirKey))
}
