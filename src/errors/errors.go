package errors

import (
	"errors"
	"strings"
)

func ErrorsToStrings(errs []error) []string {
	var s []string
	for _, err := range errs {
		s = append(s, err.Error())
	}
	return s
}

func ErrorsToString(errs []error) string {
	s := ErrorsToStrings(errs)
	return strings.Join(s, ", ")
}

func PrependErrors(errs []error, err error) []error {
	return append([]error{err}, errs...)
}

func ErrorsToError(errs []error) error {
	var msg []string
	for _, err := range errs {
		msg = append(msg, err.Error())
	}
	return errors.New(strings.Join(msg, ", "))
}
