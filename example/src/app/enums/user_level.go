package enums

import (
	"database/sql/driver"

	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
)

type UserLevel string

const (
	UserLevelSuper    UserLevel = "super"
	UserLevelAdmin    UserLevel = "admin"
	UserLevelUser     UserLevel = "user"
	UserLevelCustomer UserLevel = "customer"
)

func UserLevelFromString(value string) (UserLevel, error) {
	values := []UserLevel{
		UserLevelSuper,
		UserLevelAdmin,
		UserLevelUser,
	}
	for _, v := range values {
		if string(v) == value {
			return UserLevel(value), nil
		}
	}
	return "", errors.Errorf("'%s' is not a valid user level", value)
}

func (t *UserLevel) Scan(value any) error {
	stringValue, err := goat.ValueToString(value)
	if err != nil {
		return err
	}
	result, err := UserLevelFromString(stringValue)
	if err != nil {
		return err
	}
	*t = result
	return nil
}

func (t *UserLevel) Value() (driver.Value, error) {
	return string(*t), nil
}

func (t *UserLevel) String() string {
	return string(*t)
}
