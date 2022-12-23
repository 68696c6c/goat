package goat

import "github.com/68696c6c/goat/utils"

// The underlying utils functions are in a sub-package so that the Goat internals can use them without circular imports.

func ErrorsToStrings(errs []error) []string {
	return utils.ErrorsToStrings(errs)
}

func ErrorsToString(errs []error) string {
	return utils.ErrorsToString(errs)
}

func PrependErrors(errs []error, err error) []error {
	return utils.PrependErrors(errs, err)
}

func ErrorsToError(errs []error) error {
	return utils.ErrorsToError(errs)
}

func ExitError(err error) {
	utils.ExitError(err)
}

func ExitErrors(errs []error) {
	utils.ExitErrors(errs)
}

func ExitSuccess() {
	utils.ExitSuccess()
}

// RecordNotFound returns true if the provided error is the Gorm "record not found" error.
func RecordNotFound(err error) bool {
	return utils.RecordNotFound(err)
}

// ErrorBesidesRecordNotFound true if the provided error is NOT a Gorm "record not found" error.
func ErrorBesidesRecordNotFound(err error) bool {
	return utils.ErrorBesidesRecordNotFound(err)
}
