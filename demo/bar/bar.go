package bar

import (
	"github.com/68696c6c/goat"
)

func BarRoot() string {
	return "Root path from bar package: " + goat.Root()
}

func BarCWD() string {
	return "CWD path from bar package: " + goat.CWD()
}