package main

import (
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

func RandomBool() bool {
	n := RandomInt(0, 1)
	return n == 1
}

func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func RandomInt32(min, max int32) int32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int31n(max-min) + min
}

func RandomInt64(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

func RandomDecimal(min, max int64) decimal.Decimal {
	exp := RandomInt32(-10, 10)
	f := RandomInt64(min*1000, max*1000)
	return decimal.New(f, exp)
}

func RandomDecimalExp(min, max int64, exp int32) decimal.Decimal {
	f := RandomInt64(min*1000, max*1000)
	return decimal.New(f, exp)
}

func RandomIndex(length int) int {
	return RandomInt(0, length)
}
