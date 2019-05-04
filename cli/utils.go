package cli

import (
	"fmt"
	"os"
)

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
