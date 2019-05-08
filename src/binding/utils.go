package binding

import (
	"reflect"

	"github.com/fatih/structtag"
	"github.com/pkg/errors"
)

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
