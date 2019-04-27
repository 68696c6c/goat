package testing

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomDecimal(t *testing.T) {
	d := RandomDecimal(0, 1)
	one, err := decimal.NewFromString("1")
	require.Nil(t, err)
	assert.True(t, d.GreaterThanOrEqual(decimal.Zero))
	assert.True(t, d.LessThanOrEqual(one))
}
