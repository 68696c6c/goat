package foo

import (
	"github.com/68696c6c/goat"
	"goat/example/bar"
)

func FooRoot() []string {
	s := []string{
		"Root path from foo package: " + goat.Root(),
		bar.BarRoot(),
	}
	return s
}

func FooCWD() []string {
	s := []string{
		"CWD path from foo package: " + goat.CWD(),
		bar.BarCWD(),
	}
	return s
}
