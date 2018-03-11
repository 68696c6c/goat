package goat

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r, "Init() panicked")
		assert.Len(t, GetErrors(), 0, "Init() created errors")
	}()
	Init()
}
