package goat

import (
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func ExitError(err error) {
	l := log.New(os.Stderr, "", 0)
	l.Println(err)
	os.Exit(1)
}

func ExitSuccess() {
	os.Exit(0)
}

// RecordNotFound returns true if the provided error is the Gorm "record not found" error.
func RecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// ErrorBesidesRecordNotFound true if the provided error is NOT a Gorm "record not found" error.
func ErrorBesidesRecordNotFound(err error) bool {
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func ErrorsToError(errs []error) error {
	if len(errs) < 1 {
		return nil
	}
	var msg []string
	for _, err := range errs {
		msg = append(msg, err.Error())
	}
	return errors.New(strings.Join(msg, "; "))
}

func MakeValidationError(field, rule string) error {
	return errors.Errorf("%s failed on validation '%s'", field, rule)
}
