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

func HasErrors() bool {
	return haveErrors
}

func GetErrors() []error {
	return errs
}

func ErrorsToStrings(i interface{}) []string {
	if errs, ok := i.([]error); ok {
		var s []string
		for _, e := range errs {
			s = append(s, e.Error())
		}
		return s
	}
	if err, ok := i.(error); ok {
		return []string{err.Error()}
	}
	addError("failed to cast error or []error to []string")
	return []string{}
}

func ErrorsToString(i interface{}) string {
	s := ErrorsToStrings(i)
	return strings.Join(s, ", ")
}
