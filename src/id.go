package goat

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

// ID is a UUID type that implements a custom Value function for storing UUIDs as binary(16) columns in a database using Gorm.
type ID uuid.UUID

func (id *ID) Scan(src any) error {
	uid := (*uuid.UUID)(id)
	return uid.Scan(src)
}

func (id ID) Value() (driver.Value, error) {
	return id[:], nil
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

// MarshalText implements encoding.TextMarshaler.
func (id ID) MarshalText() ([]byte, error) {
	return uuid.UUID(id).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (id *ID) UnmarshalText(data []byte) error {
	uid := (*uuid.UUID)(id)
	return uid.UnmarshalText(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (id ID) MarshalBinary() ([]byte, error) {
	return id[:], nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (id *ID) UnmarshalBinary(data []byte) error {
	uid := (*uuid.UUID)(id)
	return uid.UnmarshalBinary(data)
}

func (id ID) Valid() bool {
	return uuid.UUID(id) != uuid.Nil
}

func NewID() ID {
	return ID(uuid.New())
}

func NilID() ID {
	return ID(uuid.Nil)
}

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}

func ParseAllIDs(s []string) ([]ID, error) {
	var ids []ID
	for _, i := range s {
		parsed, err := ParseID(i)
		if err != nil {
			return nil, err
		}

		ids = append(ids, parsed)
	}
	return ids, nil
}
