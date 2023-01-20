package hal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewMessage(t *testing.T) {
	result := NewMessage("hello world", "https://test.com/hello", user{Id: 1, Name: "Example"})
	assert.Equal(t, "hello world", result.Message)
	assert.Equal(t, user{Id: 1, Name: "Example"}, result.Embeds)
	assert.Equal(t, "https://test.com/hello", result.Links["self"].Href)
}
