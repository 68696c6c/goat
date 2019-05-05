package http

import (
	"fmt"
	"net"
	"strconv"
)

// Determines whether the provided value is a valid port that can be listened on.
func validPort(port string) error {

	// Must be numeric.
	if _, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("%s is not a valid port", port)
	}

	// Try and listen to see if the port is available.
	if ln, err := net.Listen("tcp", ":"+port); err == nil {
		_ = ln.Close()
		return nil
	}

	return fmt.Errorf("port %s is already in use", port)
}
