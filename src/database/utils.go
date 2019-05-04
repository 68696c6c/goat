package database

import "github.com/jinzhu/gorm"

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
