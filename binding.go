package goat

import (
	"errors"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/fatih/structtag"
	"gopkg.in/go-playground/validator.v8"
)

var (
	jsonTagExcludedStructs []string
	jsonTagPath            []reflect.StructField
)

type RequestValidator struct {
	*validator.Validate
}

func (v *RequestValidator) ValidateStruct(i interface{}) error {
	return v.Struct(i)
}

// Set structs that will be skipped while traversing structs to find field json
// names. Do this to avoid needlessly traversing struct fields.
func SetJSONTagExcludedStructs(s []string) {
	jsonTagExcludedStructs = s
}

func GetStructFieldValidationMeta(strct reflect.Type, e *validator.FieldError) (string, string, error) {
	logger := NewCustomLogger("reflection")
	logger.Level = logrus.DebugLevel

	logger.Debug("------------------")
	// The target field is the field from the error.
	targetField := e.Field
	logger.Debug("Searching for TargetField '" + targetField + "' on SourceStruct '" + strct.String() + "'")

	// See if the target field exists on the source struct.
	var field reflect.StructField
	jsonTagPath = []reflect.StructField{}
	field, ok := strct.FieldByName(targetField)
	if !ok {
		// Loop over all the SourceStruct struct fields until we find it.
		logger.Debug("Checking SourceStruct fields for TargetField")
		var err error
		field, err = traverseStruct(strct, targetField, logger)
		if err != nil {
			return "", "", err
		}
	}

	// If it does, we're done.
	logger.Debug("Target field exists on parent struct.")
	return getFieldValidationMeta(field, logger)
}

func structExcluded(s string) bool {
	for _, a := range jsonTagExcludedStructs {
		if a == s {
			return true
		}
	}
	return false
}

func getFieldValidationMeta(field reflect.StructField, logger *logrus.Logger) (string, string, error) {
	jsonTag, err := getJsonTag(field)
	if err != nil {
		logger.Error("getFieldValidationMeta failed getting tag name: " + err.Error())
		return "", "", err
	}
	logger.Debug("json field name is " + jsonTag)

	logger.Debug("building field path")
	var s []string
	for _, p := range jsonTagPath {
		logger.Debug("getting tag name for " + p.Name)
		pTag, err := getJsonTag(p)
		if err == nil {
			s = append(s, pTag)
		}
	}
	s = append(s, jsonTag)

	sPath := strings.Join(s, ".")

	label, err := getLabelTag(field)
	if err != nil {
		label = field.Name
	}

	return sPath, label, nil
}

func getJsonTag(field reflect.StructField) (string, error) {
	tag := string(field.Tag)
	tags, err := structtag.Parse(tag)
	if err != nil {
		return "", err
	}
	jsonTag, err := tags.Get("json")
	if err != nil {
		return "", err
	}
	return jsonTag.Name, nil
}

func getLabelTag(field reflect.StructField) (string, error) {
	tag := string(field.Tag)
	tags, err := structtag.Parse(tag)
	if err != nil {
		return "", err
	}
	labelTag, err := tags.Get("label")
	if err != nil {
		return "", err
	}
	return labelTag.Name, nil
}

func traverseStruct(strct reflect.Type, targetFieldName string, logger *logrus.Logger) (reflect.StructField, error) {
	logger.Debug("------ traverseStruct ------")
	var field reflect.StructField

	logger.Debug("checking struct " + strct.String() + " for field " + targetFieldName)
	ok := false
	if field, ok = strct.FieldByName(targetFieldName); ok {
		logger.Debug("found field '" + field.Name + "' on struct '" + strct.String() + "'")
		return field, nil
	}
	logger.Debug("field not found, checking struct for struct fields")

	for i := 0; i < strct.NumField(); i++ {
		// We only need to check fields that are structs or arrays of structs.
		sf := strct.Field(i)
		sft := sf.Type
		logger.Debug("field '" + sf.Name + "' is type '" + sft.String() + "'")

		kind := sft.Kind()
		logger.Debug("field '" + sf.Name + "' is kind '" + kind.String() + "'")
		switch kind {

		case reflect.Struct:
			if structExcluded(sft.String()) {
				logger.Debug("struct type " + sft.String() + " is excluded from traversing, skipping")
				continue
			}
			logger.Debug("checking field '" + sf.Name + "' for targetField " + targetFieldName)
			ok := false
			if field, ok = sft.FieldByName(targetFieldName); ok {
				logger.Debug("found field '" + field.Name + "' on struct '" + strct.String() + "'")
				jsonTagPath = append(jsonTagPath, sf)
				return field, nil
			}
			logger.Debug("target field " + targetFieldName + " not found on field " + sf.Name + ", diving")
			f, e := traverseStruct(sft, targetFieldName, logger)
			if e != nil {
				// If we didn't find it, keep looking here.
				continue
			}
			// If we didn't get an error, then we found it.
			jsonTagPath = append(jsonTagPath, sf)
			return f, nil

		case reflect.Slice, reflect.Array:
			sftElem := sft.Elem()
			if sftElem.Kind() == reflect.Ptr {
				logger.Debug("field '" + sf.Name + "' is an array of pointers " + sftElem.String())
				sftElem = sftElem.Elem()
				logger.Debug("new elem type is " + sftElem.String())
			}
			if sftElem.Kind() == reflect.Struct {
				logger.Debug("field '" + sf.Name + "' is an array of " + sftElem.String() + " structs, diving")
				f, e := traverseStruct(sftElem, targetFieldName, logger)
				if e != nil {
					// If we didn't find it, keep looking here.
					continue
				}
				// If we didn't get an error, then we found it.
				jsonTagPath = append(jsonTagPath, sf)
				return f, nil
			}
			logger.Debug("field '" + sf.Name + "' is an array of non-struct type " + sftElem.String() + ", skipping")

		default:
			logger.Debug("field '" + sf.Name + "' is not a struct or array, skipping")
			continue
		}
	}
	msg := "traverseStruct: failed to find field " + targetFieldName + " on " + strct.Name()
	logger.Error(msg)
	return field, errors.New(msg)
}
