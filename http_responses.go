package goat

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Message string                 `json:"message"`
	Errors  []string               `json:"errors,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
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

func RespondNotFoundError(c *gin.Context, errs []error) {
	c.JSON(http.StatusOK, Response{"Not found.", ErrorsToStrings(errs), nil})
	c.Abort()
}

func RespondBadRequestError(c *gin.Context, errs []error) {
	c.JSON(http.StatusBadRequest, Response{"Bad Request.", ErrorsToStrings(errs), nil})
	c.Abort()
}

func RespondBadRequest(c *gin.Context, data interface{}) {
	c.JSON(http.StatusBadRequest, data)
	c.Abort()
}

func RespondUnauthorizedError(c *gin.Context, errs []error) {
	c.JSON(http.StatusUnauthorized, Response{"Unauthorized.", ErrorsToStrings(errs), nil})
	c.Abort()
}

func RespondAuthenticationError(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{"Authentication error.", []string{}, nil})
	c.Abort()
}

func RespondServerError(c *gin.Context, errs []error) {
	c.JSON(http.StatusInternalServerError, Response{"Internal server error.", ErrorsToStrings(errs), nil})
	c.Abort()
}
