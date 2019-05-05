package main

import (
	"log"
	"os"
)

func ExitError(err error) {
	l := log.New(os.Stderr, "", 0)
	l.Println(err)
	os.Exit(1)
}

func ExitErrors(errs []error) {
	l := log.New(os.Stderr, "", 0)
	for _, e := range errs {
		l.Println(e)
	}
	os.Exit(1)
}

func ExitSuccess() {
	os.Exit(0)
}
