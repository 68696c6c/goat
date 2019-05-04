package main

import (
	"fmt"
	"os"

	"github.com/68696c6c/goat/cli"
)

func main() {
	if err := cli.Goat.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
