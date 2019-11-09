package goat

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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
	SetRoot(os.Getenv("APP_BASE"))
	Init()
	assert.NotEmpty(t, Root(), "failed to set root")
}
