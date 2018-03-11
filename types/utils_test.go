package types

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGoatUtils_Initialization(t *testing.T) {
	u := NewGoatUtils()
	assert.False(t, u.IsInitialized(), "unexpected initial initialized value")
}

func TestGoatUtils_SetInitialized(t *testing.T) {
	u := NewGoatUtils()
	u.SetInitialized(true)
	defer func() {
		r := recover()
		assert.Nil(t, r, "SetInitialized failed to set initialized")
	}()
	u.MustBeInitialized()
}

func TestGoatUtils_MustBeInitialized(t *testing.T) {
	u := NewGoatUtils()
	defer func() {
		r := recover()
		assert.NotNil(t, r, "SetInitialized didn't panic")
	}()
	u.MustBeInitialized()
}
