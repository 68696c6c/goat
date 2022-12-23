package utils

import (
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

func TestArgStringD_Provided(t *testing.T) {
	a := fake.Word()
	d := fake.Word()
	v := ArgStringD(a, d)
	assert.Equal(t, a, v, "failed to return provided arg")
}

func TestArgStringD_Default(t *testing.T) {
	d := fake.Word()
	v := ArgStringD("", d)
	assert.Equal(t, d, v, "failed to return default arg")
}
