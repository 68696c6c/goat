package goat

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

func RespondServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{
		Message: "Internal server error",
		Errors:  nil,
	})
}

func RespondBadRequest(c *gin.Context, errs interface{}) {
	c.JSON(http.StatusOK, Response{
		Message: "Bad request",
		Errors:  ErrorsToStrings(errs),
	})
}

func RespondCreated(c *gin.Context, i interface{}) {
	c.JSON(http.StatusCreated, i)
}

func RespondOK(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Message: message,
		Errors:  nil,
	})
}
