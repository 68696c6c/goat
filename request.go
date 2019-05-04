package main

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var contextKeyRequest = "goat_request"

// Returns true and a request object from the Gin context created by the
// bindRequestMiddleware.  If the request is not set, false is returned and
// 500 response headers are set.
func GetRequest(c *gin.Context) interface{} {
	r, exists := c.Get(contextKeyRequest)
	if !exists {
		return nil
	}
	return r
}

// Validates an incoming request and binds the request body to the provided
// struct if the validation passes.
//
// Returns a 400 error with validation errors if binding fails.
//
// Sets the bound request as an interface{} in the Gin registry if binding
// succeeds.  You can retrieve in your handlers it like this:
//
// r, ok := goat.GetRequest(c).(*yourRequestStruct)
//
// This middleware allows you to annotate your request struct fields with
// `binding:"required"` to make required fields.
//
// @TODO it seems that if a request struct has a field that is named the same as one of it's child struct's fields that the validation messages don't prefix the field name with child struct's name
func BindRequestMiddleware(req interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := reflect.ValueOf(req)
		if value.Kind() == reflect.Ptr {
			panic("Bind struct can not be a pointer")
		}
		typ := value.Type()
		obj := reflect.New(typ).Interface()
		if err := c.ShouldBindWith(obj, binding.JSON); err != nil {

			// If no request body was sent at all, Gin will try to return 'EOF'
			// Show the user something more helpful instead.
			if err.Error() == "EOF" {
				RespondBadRequestError(c, errors.New("a request body is required"))
				return
			}

			// If binding failed due to missing fields, return a validation error.
			RespondRequestValidationError(c, err, typ)
			return
		}
		c.Set(contextKeyRequest, obj)
		return
	}
}
