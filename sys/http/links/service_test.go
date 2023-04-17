package links

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestService(t *testing.T) Service {
	baseUrl, err := url.Parse("https://example.com")
	require.Nil(t, err, "failed to parse test base url")
	return NewService(baseUrl)
}

func Test_GetBaseUrl(t *testing.T) {
	subject := newTestService(t)
	result := subject.GetBaseUrl()
	assert.Equal(t, "https://example.com", result.String())
}

func Test_GetUrl_Default(t *testing.T) {
	subject := newTestService(t)
	expected := subject.GetBaseUrl()
	result := subject.GetUrl()
	assert.Equal(t, expected.String(), result.String())
}

func Test_GetUrl_Key(t *testing.T) {
	subject := newTestService(t)

	input, err := url.Parse("https://test.org")
	require.Nil(t, err, "failed to parse test url")
	subject.SetUrl("key", input)

	result := subject.GetUrl("key")
	assert.Equal(t, input.String(), result.String())
}

func Test_SetUrl(t *testing.T) {
	subject := newTestService(t)

	inputA, err := url.Parse("https://a.test")
	require.Nil(t, err, "failed to parse test url 'a'")
	subject.SetUrl("a", inputA)

	inputB, err := url.Parse("https://a.test")
	require.Nil(t, err, "failed to parse test url 'b'")
	subject.SetUrl("b", inputB)

	assert.Equal(t, inputA.String(), subject.GetUrl("a").String())

	assert.Equal(t, inputB.String(), subject.GetUrl("b").String())
}

func Test_AddBaseUrlPath(t *testing.T) {
	subject := newTestService(t)
	subject.AddBaseUrlPath("foo", "/bar")
	result := subject.GetUrl("foo")
	assert.Equal(t, "https://example.com/bar", result.String())
}

func Test_AddBaseUrlPath_NoLeadingSlash(t *testing.T) {
	subject := newTestService(t)
	subject.AddBaseUrlPath("asdf", "a/b/c")
	result := subject.GetUrl("asdf")
	assert.Equal(t, "https://example.com/a/b/c", result.String())
}
