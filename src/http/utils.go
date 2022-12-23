package http

import (
	"fmt"
	"net"
	"reflect"
	"strconv"

	"github.com/fatih/structtag"
	"github.com/pkg/errors"
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

// Returns the value of the specified tag from provided struct field's annotations.
func getTag(f reflect.StructField, t string) (string, error) {
	tag := string(f.Tag)
	tags, err := structtag.Parse(tag)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse field tags")
	}
	jsonTag, err := tags.Get(t)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get '%s' tag", t)
	}
	return jsonTag.Name, nil
}

// Returns the value of the 'json' tag from provided struct field's annotations.
func getJsonTag(f reflect.StructField) (string, error) {
	return getTag(f, "json")
}

// Returns the value of the 'label' tag from provided struct field's annotations.
func getLabelTag(f reflect.StructField) (string, error) {
	return getTag(f, "label")
}
