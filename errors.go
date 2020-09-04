package goat

import "github.com/68696c6c/goat/utils"

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

// Returns true if the provided slice of errors
func RecordNotFound(errs []error) bool {
	return utils.RecordNotFound(errs)
}

// Returns true if there are any errors in the provided array that are NOT a 'record not found' error
func ErrorsBesidesRecordNotFound(errs []error) bool {
	return utils.ErrorsBesidesRecordNotFound(errs)
}
