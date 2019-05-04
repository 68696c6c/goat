package main

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
)

// A generic response.
// swagger:response Response
type Response struct {
	Message string                 `json:"message"`
	Errors  []string               `json:"errors,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// A validation error response.
// swagger:response ValidationResponse
type ValidationResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// A boolean response.
// swagger:response BoolResponse
type BoolResponse struct {
	Valid bool `json:"valid"`
}

func RespondValid(c *gin.Context) {
	c.JSON(http.StatusOK, BoolResponse{true})
}

func RespondInvalid(c *gin.Context) {
	c.JSON(http.StatusBadRequest, BoolResponse{false})
}

func RespondMessage(c *gin.Context, m string) {
	c.JSON(http.StatusOK, Response{m, []string{}, nil})
	c.Abort()
}

func RespondData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
	c.Abort()
}

func RespondCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
	c.Abort()
}

func RespondNotFoundError(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, Response{"Not found.", []string{err.Error()}, nil})
	c.Abort()
}

func RespondNotFoundErrors(c *gin.Context, errs []error) {
	c.JSON(http.StatusNotFound, Response{"Not found.", ErrorsToStrings(errs), nil})
	c.Abort()
}

func RespondBadRequestErrors(c *gin.Context, errs []error) {
	c.JSON(http.StatusBadRequest, Response{"Bad Request.", ErrorsToStrings(errs), nil})
	c.Abort()
}

func RespondBadRequestError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{"Bad Request.", []string{err.Error()}, nil})
	c.Abort()
}

func RespondBadRequest(c *gin.Context, data interface{}) {
	c.JSON(http.StatusBadRequest, data)
	c.Abort()
}

func RespondValidationError(c *gin.Context, errs map[string]error) {
	msgs := make(map[string]string)
	for k, v := range errs {
		msgs[k] = v.Error()
	}
	c.JSON(http.StatusBadRequest, ValidationResponse{"Invalid Request.", msgs})
	c.Abort()
}

func RespondRequestValidationError(c *gin.Context, err error, t reflect.Type) {
	msgs := make(map[string]string)

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusBadRequest, Response{"Invalid Request.", []string{}, nil})
		c.Abort()
		return
	}
	for _, e := range ve {
		jsonName, label, err := GetStructFieldValidationMeta(t, e)
		if err != nil {
			jsonName = e.Field
			label = e.Name
		}

		msgs[jsonName] = label + " is required"
	}
	c.JSON(http.StatusBadRequest, ValidationResponse{"Invalid Request.", msgs})
	c.Abort()
}

func RespondUnauthorizedError(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{"Unauthorized.", []string{}, nil})
	c.Abort()
}

func RespondAuthenticationError(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{"Authentication error.", []string{}, nil})
	c.Abort()
}

func RespondServerErrors(c *gin.Context, errs []error) {
	// Only show errors to the user if we are in debug mode.
	if gin.Mode() != gin.DebugMode {
		errs = []error{}
	}
	c.JSON(http.StatusInternalServerError, Response{"Internal server error.", ErrorsToStrings(errs), nil})
	c.Abort()
}

func RespondServerError(c *gin.Context, err error) {
	// Only show errors to the user if we are in debug mode.
	if gin.Mode() != gin.DebugMode {
		err = nil
	}
	c.JSON(http.StatusInternalServerError, Response{"Internal server error.", []string{err.Error()}, nil})
	c.Abort()
}
