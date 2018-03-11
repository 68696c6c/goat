package goat

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r, "Init() panicked")
		if r != nil {
			println("recovered")
			println(r)
		}
		assert.Len(t, GetErrors(), 0, "Init() created errors")
	}()
	Init()
	assert.NotEmpty(t, Root(), "failed to set root")
}
