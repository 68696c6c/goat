package binding

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/68696c6c/goat/src/types"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v8"
)

// Requests are bound to an application-level (i.e. non-Goat) struct provided by
// the caller.
//
//
// Request struct fields can be annotated with tags to mark them as required and
// to provide custom labels for use in validation error messages.
//
// For example:
//	type userRequest struct {
// 		Name string `json:"name" binding:"required" label:"Username"`
// 	}
//
// Note that the `binding:"required"` annotation is not recursive.  For example,
// if an Order request struct has an Items field that is a slice of Item structs
// and the Item struct also has fields that are required, the binding service
// will not throw an error if one of the Items is missing a field.  To change
// this behavior, use 'dive,required' annotation to tell the binding service to
// check the field's own fields for required annotations as well.
//
// For example:
//	type orderRequest struct {
//		Items []*Item `binding:"required,dive,required"`
//	}
//
//
// If a required field is not provided, or a request body is not sent, the
// binding service will automatically send a bad request response with a list of
// errors.  The service will attempt to use the JSON names of the fields in the
// error messages rather than the struct field names so that the user sees the
// field name as they sent it, e.g. "name" instead of "Name".  If a custom label
// annotation exists on the field, that will be used instead.
//
//
// The binding service can be configured to skip certain struct types when
// binding fields in order to avoid needlessly traversing fields.  For example,
// if an application binds requests to models, the ID and timestamp fields can
// be skipped since they will usually not be included in request bodies.
//
//
// Binding can be done directly in a handler using the Bind function or in a
// middleware using the BindMiddleware function.  The middleware will set the
// bound request in the Gin context where it can be accessed by subsequent
// handler functions.
type Service interface {
	Bind(c *gin.Context, r interface{})
	BindMiddleware(r interface{}) gin.HandlerFunc
}

type Config struct {
	Debug           bool
	ExcludedStructs []string
}

type ServiceGin struct {
	debug           bool
	excludedStructs []string
	jsonTagPath     []reflect.StructField
}

func NewServiceGin(c Config) ServiceGin {
	return ServiceGin{
		debug:           c.Debug,
		excludedStructs: c.ExcludedStructs,
	}
}

// @TODO need to return an error here so that the calling handler can return.
func (s ServiceGin) Bind(c *gin.Context, r interface{}) {
	s.bind(c, r)
}

func (s ServiceGin) BindMiddleware(r interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		s.bind(c, r)
	}
}

// Attempts to bind a JSON request body from the Gin Context to the provided
// struct. If any of the struct's required fields are missing from the request
// body, a 400 response is sent.
func (s ServiceGin) bind(c *gin.Context, r interface{}) {
	value := reflect.ValueOf(r)
	if value.Kind() == reflect.Ptr {
		panic("Bind struct can not be a pointer")
	}
	typ := value.Type()
	obj := reflect.New(typ).Interface()
	if err := c.ShouldBindWith(obj, binding.JSON); err != nil {
		s.respondRequestBindingError(c, err, typ)
		return
	}
	return
}

func (s ServiceGin) log(m string) {
	if s.debug {
		println(m)
	}
}

func (s ServiceGin) respondRequestBindingError(c *gin.Context, err error, t reflect.Type) {
	// If no request body was sent at all, Gin will try to return 'EOF'
	// Show the user something more helpful instead.
	if err.Error() == "EOF" {
		e := errors.New("a request body is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, types.Response{"Bad Request.", []error{e}, nil})
		return
	}

	// Cast the error to a ValidationErrors struct so that we can access detailed
	// information about the error.  If the error cannot be cast, a generic
	// Bad Request response will be sent instead of the detailed message.
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		e := errors.New("invalid request")
		c.AbortWithStatusJSON(http.StatusBadRequest, types.Response{"Bad Request.", []error{e}, nil})
		return
	}

	// Create an error message for each missing field.
	msgs := make(map[string]string)
	for _, e := range ve {
		meta, err := s.getStructFieldValidationMeta(t, e)
		if err != nil {
			// If we couldn't find a JSON tag annotation for the field, fallback to
			// the struct field name.
			meta = &fieldMeta{
				Path:  e.Field,
				Label: e.Name,
			}
		}
		msgs[meta.Path] = meta.Label + " is required"
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, types.ValidationResponse{"Invalid Request.", msgs})
}

type fieldMeta struct {
	Path  string
	Label string
}

// Recursively loop over the fields of the provided struct to find the field
func (s ServiceGin) getStructFieldValidationMeta(strct reflect.Type, e *validator.FieldError) (*fieldMeta, error) {

	// The target field is the field from the error.
	targetField := e.Field

	// See if the target field exists on the source struct.
	var field reflect.StructField
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

func (s ServiceGin) structExcluded(ss string) bool {
	for _, a := range s.excludedStructs {
		if a == ss {
			return true
		}
	}
	return false
}

// Parse the provided struct field's annotations to build an error message.
func (s ServiceGin) getFieldMeta(field reflect.StructField) (*fieldMeta, error) {
	jsonTag, err := getJsonTag(field)
	if err != nil {
		return nil, err
	}

	var ss []string
	for _, p := range s.jsonTagPath {
		pTag, err := getJsonTag(p)
		if err == nil {
			ss = append(ss, pTag)
		}
	}
	ss = append(ss, jsonTag)

	sPath := strings.Join(ss, ".")

	label, err := getLabelTag(field)
	if err != nil {
		label = jsonTag
	}

	return &fieldMeta{
		Path:  sPath,
		Label: label,
	}, nil
}

// Recursively search the provided structs fields for the provided field name,
// building a path to the field.
// @TODO if a parent struct has a child struct field and both structs have a field with the same name, this function will always return the path to the field on the parent.
func (s ServiceGin) getFieldPath(strct reflect.Type, targetFieldName string) (reflect.StructField, error) {
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

			if s.structExcluded(sft.String()) {
				// This struct type has been marked as excluded.
				continue
			}

			if field, ok := sft.FieldByName(targetFieldName); ok {
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
			s.jsonTagPath = append(s.jsonTagPath, sf)
			return f, nil

		// The field is an array.
		case reflect.Slice, reflect.Array:

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
