package meta

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
	"github.com/pkg/errors"
)

type FieldMeta struct {
	Path  string
	Label string
}

func (t FieldMeta) String() string {
	return fmt.Sprintf("path: %s, label: %s", t.Path, t.Label)
}

type Service interface {
	GetStructFieldMeta(strct reflect.Type, targetField string) (*FieldMeta, error)
}

type service struct {
	debug           bool
	jsonTagPath     []reflect.StructField
	excludedStructs []string
}

func NewService(excludedStructs string) Service {
	es := strings.Split(excludedStructs, ",")
	return &service{
		debug:           false,
		excludedStructs: es,
	}
}

func (s *service) log(m string) {
	if s.debug {
		println(m)
	}
}

func (s *service) GetStructFieldMeta(strct reflect.Type, targetField string) (*FieldMeta, error) {

	// See if the target field exists on the source struct.
	var field reflect.StructField
	s.log("------------")
	s.jsonTagPath = []reflect.StructField{}
	field, ok := strct.FieldByName(targetField)
	if !ok {
		// Loop over all the SourceStruct struct fields until we find it.
		var err error
		field, err = s.getFieldPath(strct, targetField)
		if err != nil {
			return nil, err
		}
	}

	// If it does, we're done.
	return s.getFieldMeta(field)
}

func (s *service) isStructExcluded(ss string) bool {
	s.log(fmt.Sprintf("checking excludedStructs for '%s'", ss))
	for _, a := range s.excludedStructs {
		if a == ss {
			return true
		}
	}
	return false
}

// Parse the provided struct field's annotations to build an error message.
func (s *service) getFieldMeta(field reflect.StructField) (*FieldMeta, error) {
	s.log("getFieldMeta length " + strconv.Itoa(len(s.jsonTagPath)))
	jsonTag, err := getJsonTag(field)
	if err != nil {
		return nil, err
	}

	var ss []string
	for _, p := range s.jsonTagPath {
		s.log("p: " + p.Name)
		pTag, err := getJsonTag(p)
		if err == nil {
			s.log("pTag: " + pTag)
			ss = append(ss, pTag)
		}
	}
	ss = append(ss, jsonTag)

	sPath := strings.Join(ss, ".")

	label, err := getLabelTag(field)
	if err != nil {
		label = jsonTag
	}

	return &FieldMeta{
		Path:  sPath,
		Label: label,
	}, nil
}

// Recursively search the provided structs fields for the provided field name,
// building a path to the field.
// @TODO if a parent struct has a child struct field and both structs have a field with the same name, this function will always return the path to the field on the parent.
func (s *service) getFieldPath(strct reflect.Type, targetFieldName string) (reflect.StructField, error) {
	s.log(fmt.Sprintf("getting path for field '%s' on struct '%v'", targetFieldName, strct))
	var field reflect.StructField

	// Check if the field exists on the struct.
	if field, ok := strct.FieldByName(targetFieldName); ok {
		return field, nil
	}

	// If the field doesn't exist directly on the struct, check any fields that
	// are also structs or struct slices.
	for i := 0; i < strct.NumField(); i++ {

		sf := strct.Field(i)
		sft := sf.Type
		kind := sft.Kind()
		switch kind {

		// The field is a struct; check its fields for the field.
		case reflect.Struct:
			s.log("field is a struct")

			if s.isStructExcluded(sft.String()) {
				// This struct type has been marked as excluded.
				s.log("field is flat, skipping")
				continue
			}

			// Check this field's fields for the targetField.
			if field, ok := sft.FieldByName(targetFieldName); ok {
				s.log(fmt.Sprintf("field found on field '%v'", sft))
				// The field exists on this struct.
				s.jsonTagPath = append(s.jsonTagPath, sf)
				return field, nil
			}

			// Recurse through the field's fields.
			f, e := s.getFieldPath(sft, targetFieldName)
			if e != nil {
				// The field was not found here; continue checking the other fields.
				continue
			}

			// If we didn't get an error, then we found it.
			s.log(fmt.Sprintf("field found on child field of '%v', adding '%v' to path", sft, sf))
			s.jsonTagPath = append(s.jsonTagPath, sf)
			return f, nil

		// The field is an array.
		case reflect.Slice, reflect.Array:
			s.log("field is an array")

			// Get the type of the array's elements.
			sftElem := sft.Elem()
			if sftElem.Kind() == reflect.Ptr {
				sftElem = sftElem.Elem()
			}

			// If the elements are structs, recurse through the element's fields.
			if sftElem.Kind() == reflect.Struct {
				f, e := s.getFieldPath(sftElem, targetFieldName)
				if e != nil {
					// The field was not found here; continue checking the other fields.
					continue
				}

				// If we didn't get an error, then we found it.
				s.jsonTagPath = append(s.jsonTagPath, sf)
				return f, nil
			}

		// The field is not a struct or an array of structs; skip it.
		default:
			continue
		}
	}

	return field, errors.Errorf("failed to find field '%s' on '%s'", targetFieldName, strct.Name())
}

// Returns the value of the specified tag from provided struct field's annotations.
func getTag(f reflect.StructField, t string) (string, error) {
	tag := string(f.Tag)
	tags, err := structtag.Parse(tag)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse field tags")
	}
	jsonTag, err := tags.Get(t)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get '%s' tag", t)
	}
	return jsonTag.Name, nil
}

// Returns the value of the 'json' tag from provided struct field's annotations.
func getJsonTag(f reflect.StructField) (string, error) {
	return getTag(f, "json")
}

// Returns the value of the 'label' tag from provided struct field's annotations.
func getLabelTag(f reflect.StructField) (string, error) {
	return getTag(f, "label")
}
