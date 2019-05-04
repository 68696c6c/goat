package database

import (
	"testing"

	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

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
