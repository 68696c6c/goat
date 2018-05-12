package goat

import (
	"strings"
	"errors"
)

var (
	errs       []error
	haveErrors bool
)

func addError(s string) {
	haveErrors = true
	errs = append(errs, errors.New(s))
}

func addAndGetError(s string) error {
	err := errors.New(s)
	haveErrors = true
	errs = append(errs, err)
	return err
}

func HasErrors() bool {
	return haveErrors
}

func GetErrors() []error {
	return errs
}

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
