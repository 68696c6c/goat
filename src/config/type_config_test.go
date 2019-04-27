package config

import (
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_FileName(t *testing.T) {
	file := fake.Word()
	c := newConfig(file, "")
	assert.Equal(t, file, c.FileName(), "failed to set file name")
}

func TestConfig_FilePath(t *testing.T) {
	path := fake.Word()
	c := newConfig("", path)
	assert.Equal(t, path, c.FilePath(), "failed to set file name")
}
