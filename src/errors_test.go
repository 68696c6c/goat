package goat

import (
	"testing"

	"github.com/icrowley/fake"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_RecordNotFound_True(t *testing.T) {
	result := RecordNotFound(gorm.ErrRecordNotFound)
	assert.True(t, result, "failed to recognize 'record not found' error")
}

func Test_RecordNotFound_True_Wrapped(t *testing.T) {
	input := errors.Wrap(gorm.ErrRecordNotFound, fake.Word())
	result := RecordNotFound(input)
	assert.True(t, result, "failed to recognize wrapped 'record not found' error")
}

func Test_RecordNotFound_False(t *testing.T) {
	result := RecordNotFound(errors.New(fake.Word()))
	assert.False(t, result, "false positive looking for 'record not found' error")
}

func Test_ErrorBesidesRecordNotFound_True(t *testing.T) {
	input := errors.Wrap(errors.New(fake.Word()), fake.Word())
	result := ErrorBesidesRecordNotFound(input)
	assert.True(t, result, "failed to recognize 'record not found' error")
}

func Test_ErrorBesidesRecordNotFound_False(t *testing.T) {
	input := errors.Wrap(gorm.ErrRecordNotFound, fake.Word())
	result := ErrorBesidesRecordNotFound(input)
	assert.False(t, result, "false positive looking for 'record not found' error")
}
