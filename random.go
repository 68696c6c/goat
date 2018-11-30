package goat

import (
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

func RandomBool() bool {
	n := RandomInt(0, 1)
	return n == 1
}

func RandomInt(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

func RandomDecimal(min, max int64) decimal.Decimal {
	f := RandomInt(min*1000, max*1000)
	return decimal.New(f, -3)
}

func RandomIndex(length int) int {
	return int(RandomInt(0, int64(length)))
}
