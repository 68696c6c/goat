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

// Returns true if the provided slice of errors contains a GORM "record not found" error.
func RecordNotFound(errs []error) bool {
	return utils.RecordNotFound(errs)
}

// Returns true if there are any errors in the provided array that are NOT a GORM "record not found" error.
func ErrorsBesidesRecordNotFound(errs []error) bool {
	return utils.ErrorsBesidesRecordNotFound(errs)
}

// Returns true if the provided error is the GORM "record not found" error.  Note that this function only returns true
// on exact matches, so wrapped or flattened errors containing the "record not found" error will not be caught.  Given
// that the GORM "record not found" error is rather generic, handling those cases is too likely to result in false
// positives.  For this reason, the check should be performed before wrapping or flattening any GORM errors using the
// other helper functions in this file.
func IsNotFoundError(err error) bool {
	return utils.IsNotFoundError(err)
}
