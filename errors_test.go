package goat

import (
	"testing"

	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// RecordNotFound tests.

func TestRecordNotFound_True(t *testing.T) {
	e := []error{
		gorm.ErrRecordNotFound,
		errors.New(fake.Word()),
		errors.New(fake.Word()),
	}
	v := RecordNotFound(e)
	assert.True(t, v, "failed to find 'record not found' error")
}

func TestRecordNotFound_False(t *testing.T) {
	e := []error{
		errors.New(fake.Word()),
		errors.New(fake.Word()),
		errors.New(fake.Word()),
	}
	v := RecordNotFound(e)
	assert.False(t, v, "false positive looking for 'record not found' error")
}

// ErrorsBesidesRecordNotFound tests.

func TestErrorsBesidesRecordNotFound_True(t *testing.T) {
	e := []error{
		errors.New(fake.Word()),
		errors.New(fake.Word()),
		errors.New(fake.Word()),
	}
	v := ErrorsBesidesRecordNotFound(e)
	assert.True(t, v, "failed to find 'record not found' error")
}

func TestErrorsBesidesRecordNotFound_False(t *testing.T) {
	e := []error{
		gorm.ErrRecordNotFound,
	}
	v := ErrorsBesidesRecordNotFound(e)
	assert.False(t, v, "false positive looking for 'record not found' error")
}

// IsNotFoundError tests.

func TestIsNotFoundError_True(t *testing.T) {
	v := IsNotFoundError(gorm.ErrRecordNotFound)
	assert.True(t, v, "failed to recognize 'record not found' error")
}

func TestIsNotFoundError_False(t *testing.T) {
	v := IsNotFoundError(errors.New(fake.Sentence()))
	assert.False(t, v, "false positive looking for 'record not found' error")
}

func TestIsNotFoundError_Flattened_False(t *testing.T) {
	// Given that the GORM "record not found" is rather generic, IsNotFoundError intentionally does not support flattened
	// or wrapped "record not found" errors to avoid the potential of false positives.
	e := errors.Wrap(ErrorsToError([]error{
		gorm.ErrRecordNotFound,
		errors.New(fake.Sentence()),
	}), fake.Sentence())
	v := IsNotFoundError(e)
	assert.False(t, v, "false positive looking for 'record not found' error")
}
