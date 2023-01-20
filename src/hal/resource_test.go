package hal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewResource(t *testing.T) {
	result := NewResource("https://test.com/example", user{Id: 1, Name: "Example"})
	assert.Equal(t, user{Id: 1, Name: "Example"}, result.Embeds)
	assert.Equal(t, "https://test.com/example", result.Links["self"].Href)
}
