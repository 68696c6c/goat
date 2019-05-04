package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func ExitIfProduction() {
	env := viper.GetString("env")
	if env == gin.ReleaseMode {
		err := errors.New("this command cannot be run in production")
		ExitWithError(err)
	}
}

// Determines whether the provided value is a valid port that can be listened on.
func ValidPort(port string) (valid bool, message string) {

	// Must be numeric.
	if _, err := strconv.Atoi(port); err != nil {
		return false, "Not a valid port."
	}

	// Try and listen to see if the port is available.
	if ln, err := net.Listen("tcp", ":"+port); err == nil {
		_ = ln.Close()
		return true, ""
	}

	return false, "Port is already in use."
}

func GetCountArg(countArg string) (int, error) {
	count := 1
	c, err := strconv.Atoi(countArg)
	if err != nil {
		return 0, fmt.Errorf("failed to parse count argument: %v", err.Error())
	}
	count = c
	return count, nil
}

func ExitWithError(err error) {
	fmt.Fprintf(os.Stderr, err.Error()+"\n")
	os.Exit(1)
}

func ExitWithErrors(errs []error) {
	for _, e := range errs {
		fmt.Fprintf(os.Stderr, e.Error()+"\n")
	}
	os.Exit(1)
}

func ExitSuccess() {
	os.Exit(0)
}
