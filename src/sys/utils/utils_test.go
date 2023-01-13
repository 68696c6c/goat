package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValueOrDefault_Input(t *testing.T) {
	input := "foo"
	result := ValueOrDefault[string](&input, "bar")
	assert.Equal(t, input, result)
}

func Test_ValueOrDefault_Default(t *testing.T) {
	result := ValueOrDefault[string](nil, "bar")
	assert.Equal(t, "bar", result)
}
