package utils

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
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

// Returns true if the provided slice of errors
func RecordNotFound(errs []error) bool {
	for _, err := range errs {
		if err == gorm.ErrRecordNotFound {
			return true
		}
	}
	return false
}

// Returns true if there are any errors in the provided array that are NOT a 'record not found' error
func ErrorsBesidesRecordNotFound(errs []error) bool {
	for _, e := range errs {
		if e != gorm.ErrRecordNotFound {
			return true
		}
	}
	return false
}
