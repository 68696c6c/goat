package goat

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/68696c6c/goat/query"
)

type testModel struct{}

func Test_ApplyQueryToGorm(t *testing.T) {
	MustInit()
	db, err := GetMainDB()
	require.Nil(t, err)
	require.NotNil(t, db)

	result := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tdb := tx.Model(&testModel{})
		err = ApplyQueryToGorm(tdb, query.NewQuery().WhereIs("a", nil), false)
		require.Nil(t, err)
		temp := testModel{}
		return tdb.Find(&temp)
	})
	assert.Equal(t, "SELECT * FROM `test_models` WHERE (a IS NULL)", result)
}
