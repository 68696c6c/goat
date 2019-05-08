package main

import (
	"net/http"

	"github.com/68696c6c/goat/src/types"

	"github.com/gin-gonic/gin"
)

func RespondValid(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, types.BoolResponse{true})
}

func RespondInvalid(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, types.BoolResponse{false})
}

func RespondMessage(c *gin.Context, m string) {
	c.AbortWithStatusJSON(http.StatusOK, types.Response{m, []error{}, nil})
}

func RespondData(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, data)
}

func RespondCreated(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusCreated, data)
}

func RespondNotFoundError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusNotFound, types.Response{"Not found.", []error{err}, nil})
}

func RespondNotFoundErrors(c *gin.Context, errs []error) {
	c.AbortWithStatusJSON(http.StatusNotFound, types.Response{"Not found.", errs, nil})
}

func RespondBadRequestErrors(c *gin.Context, errs []error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, types.Response{"Bad Request.", errs, nil})
}

func RespondBadRequestError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, types.Response{"Bad Request.", []error{err}, nil})
}

func RespondValidationError(c *gin.Context, errs map[string]error) {
	msgs := make(map[string]string)
	for k, v := range errs {
		msgs[k] = v.Error()
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, types.ValidationResponse{"Invalid Request.", msgs})
}

func RespondUnauthorizedError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, types.Response{"Unauthorized.", []error{}, nil})
}

func RespondAuthenticationError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, types.Response{"Authentication error.", []error{}, nil})
}

func RespondServerErrors(c *gin.Context, errs []error) {
	// Only show errors to the user if we are in debug mode.
	if !DebugEnabled() {
		errs = []error{}
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, types.Response{"Internal server error.", errs, nil})
}

func RespondServerError(c *gin.Context, err error) {
	// Only show errors to the user if we are in debug mode.
	if !DebugEnabled() {
		err = nil
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, types.Response{"Internal server error.", []error{err}, nil})
}
