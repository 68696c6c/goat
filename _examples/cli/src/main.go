package main

import (
	"os"

	"github.com/68696c6c/cli/cmd"
)

func main() {
	if err := cmd.Root.Execute(); err != nil {
		println(err)
		os.Exit(-1)
	}
}
