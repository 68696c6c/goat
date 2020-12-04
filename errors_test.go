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
	// IsNotFoundError intentionally does not support flattened or wrapped "record not found" errors.
	e := errors.Wrap(ErrorsToError([]error{
		gorm.ErrRecordNotFound,
		errors.New(fake.Sentence()),
	}), fake.Sentence())
	v := IsNotFoundError(e)
	assert.False(t, v, "false positive looking for 'record not found' error")
}

// IsOrContainsNotFoundError tests.

func TestIsOrContainsNotFoundError_Single_True(t *testing.T) {
	v := IsOrContainsNotFoundError(gorm.ErrRecordNotFound)
	assert.True(t, v, "failed to recognize 'record not found' error")
}

func TestIsOrContainsNotFoundError_Flattened_True(t *testing.T) {
	e := ErrorsToError([]error{
		gorm.ErrRecordNotFound,
		errors.New(fake.Word()),
	})
	v := IsOrContainsNotFoundError(e)
	assert.True(t, v, "failed to recognize 'record not found' error in a flattened slice of errors")
}

func TestIsOrContainsNotFoundError_Wrapped_True(t *testing.T) {
	e := errors.Wrap(gorm.ErrRecordNotFound, fake.Sentence())
	v := IsOrContainsNotFoundError(e)
	assert.True(t, v, "failed to recognize 'record not found' error in a wrapped error")
}

func TestIsOrContainsNotFoundError_FlattenedWrapped_True(t *testing.T) {
	e := errors.Wrap(ErrorsToError([]error{
		gorm.ErrRecordNotFound,
		errors.New(fake.Sentence()),
	}), fake.Sentence())
	v := IsOrContainsNotFoundError(e)
	assert.True(t, v, "failed to recognize 'record not found' error in a wrapped and flattened error")
}

func TestIsOrContainsNotFoundError_Single_False(t *testing.T) {
	e := errors.New(fake.Sentence())
	v := IsOrContainsNotFoundError(e)
	assert.False(t, v, "false positive looking for 'record not found' error")
}

func TestIsOrContainsNotFoundError_Flattened_False(t *testing.T) {
	e := ErrorsToError([]error{
		errors.New(fake.Sentence()),
		errors.New(fake.Sentence()),
	})
	v := IsOrContainsNotFoundError(e)
	assert.False(t, v, "false positive looking for 'record not found' error in a flattened slice of errors")
}

func TestIsOrContainsNotFoundError_Wrapped_False(t *testing.T) {
	e := errors.Wrap(errors.New(fake.Sentence()), fake.Sentence())
	v := IsOrContainsNotFoundError(e)
	assert.False(t, v, "false positive looking for 'record not found' error in a wrapped error")
}

func TestIsOrContainsNotFoundError_FlattenedWrapped_False(t *testing.T) {
	e := errors.Wrap(ErrorsToError([]error{
		errors.New(fake.Sentence()),
		errors.New(fake.Sentence()),
	}), fake.Sentence())
	v := IsOrContainsNotFoundError(e)
	assert.False(t, v, "false positive looking for 'record not found' error in a wrapped and flattened error")
}
