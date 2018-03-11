package types

import (
	"testing"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

func TestConfig_FileName(t *testing.T) {
	u := NewGoatUtils()
	file := fake.Word()
	c := NewConfig(u, file, "")
	assert.Equal(t, file, c.FileName(), "failed to set file name")
}

func TestConfig_FilePath(t *testing.T) {
	u := NewGoatUtils()
	path := fake.Word()
	c := NewConfig(u, "", path)
	assert.Equal(t, path, c.FilePath(), "failed to set file name")
}
