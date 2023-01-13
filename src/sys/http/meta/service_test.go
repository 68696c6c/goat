package meta

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StructField struct {
	FieldOne string `json:"field_one"`
	FieldTwo string `json:"field_two"`
}

type ExcludedStructField struct {
	FlatFieldOne string `json:"flat_field_one"`
	FlatFieldTwo string `json:"flat_field_two"`
}

type testStruct struct {
	Field         StructField         `json:"field"`
	ExcludedField ExcludedStructField `json:"excluded"`
}

func TestBinding_GetStructFieldMeta(t *testing.T) {
	s := NewService("")

	r := testStruct{
		Field: StructField{
			FieldOne: "one",
			FieldTwo: "two",
		},
	}

	value := reflect.ValueOf(r)
	typ := value.Type()

	m, err := s.GetStructFieldMeta(typ, "FieldOne")

	require.Nil(t, err, "unexpected error returned")
	assert.Equal(t, "field.field_one", m.Path, "unexpected field path returned")
	println(m.String())
}

func TestBinding_GetStructFieldMeta_Flat(t *testing.T) {
	s := NewService("meta.ExcludedStructField")

	r := testStruct{
		ExcludedField: ExcludedStructField{
			FlatFieldOne: "one",
			FlatFieldTwo: "two",
		},
	}

	value := reflect.ValueOf(r)
	typ := value.Type()

	_, err := s.GetStructFieldMeta(typ, "FlatFieldOne")

	assert.NotNil(t, err, "unexpected error returned")
}
