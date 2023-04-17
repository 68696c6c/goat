package goat

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func AssertValidID(t *testing.T, id ID) {
	assert.NotEqual(t, NilID(), id, "id incorrect for a new record")
}

func AssertValidInitialDeletedAt(t *testing.T, d *time.Time) {
	assert.Nil(t, d, "deleted_at value incorrect for a new record")
}

func AssertValidInitialCreatedAt(t *testing.T, d time.Time) {
	assert.WithinDuration(t, time.Now(), d, 1*time.Second, "created_at value incorrect for a new record")
}

func AssertValidInitialUpdatedAt(t *testing.T, d *time.Time) {
	assert.Nil(t, d, "updated_at value incorrect for a new record")
}

func AssertValidModifiedUpdatedAt(t *testing.T, d *time.Time) {
	assert.WithinDuration(t, time.Now(), *d, 1*time.Second, "updated_at value incorrect for an updated record")
}

func RequireDecimalEqual(t *testing.T, exp, act any, msgAndArgs ...any) {
	msg := messageFromMsgAndArgs(msgAndArgs...)

	expDec, ok := exp.(decimal.Decimal)
	require.True(t, ok, msg+"\n not a decimal: %s", exp)
	actDec, ok := act.(decimal.Decimal)
	require.True(t, ok, msg+"\n not a decimal: %s", exp)

	require.True(t, expDec.Equal(actDec), msg+"\n expected %s to equal %s", actDec, expDec)
}

func AssertRecordDeleted[M any](t *testing.T, db *gorm.DB, input M, msg string) {
	err := db.First(input).Error
	assert.NotNil(t, err)
	assert.True(t, RecordNotFound(err), msg)
}

func messageFromMsgAndArgs(msgAndArgs ...any) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		fmt.Printf("%+v", msgAndArgs)
		return msgAndArgs[0].(string)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}
