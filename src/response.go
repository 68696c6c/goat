package goat

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/68696c6c/goat/hal"
)

type ApiProblem hal.ApiProblem

type Resource hal.Resource

type Collection[T any] hal.Collection[T]

func debugError(err error) error {
	if g.HTTP.DebugEnabled() {
		return err
	}
	return nil
}

func logHandlerError(cx *gin.Context, err error) {
	GetLogger().Error(fmt.Sprintf("%s | %s", cx.HandlerName(), err))
}

func RespondOk(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusOK, data)
}

func RespondNoContent(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusNoContent)
}

func RespondUsed(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusIMUsed, data)
}

func RespondAccepted(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusAccepted, data)
}

func RespondCreated(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusCreated, data)
}

func RespondNotFound(cx *gin.Context, err error) {
	status := http.StatusNotFound
	logHandlerError(cx, err)
	cx.AbortWithStatusJSON(status, hal.NewApiProblem(status, debugError(err)))
}

func RespondBadRequest(cx *gin.Context, err error) {
	status := http.StatusBadRequest
	logHandlerError(cx, err)
	cx.AbortWithStatusJSON(status, hal.NewApiProblem(status, debugError(err)))
}

func RespondValidationError(cx *gin.Context, err error) {
	status := http.StatusBadRequest
	logHandlerError(cx, err)
	cx.AbortWithStatusJSON(status, hal.NewApiProblem(status, err))
}

func RespondUnauthorized(cx *gin.Context, err error) {
	status := http.StatusUnauthorized
	logHandlerError(cx, err)
	cx.AbortWithStatusJSON(status, hal.NewApiProblem(status, debugError(err)))
}

func RespondForbidden(cx *gin.Context, err error) {
	status := http.StatusForbidden
	logHandlerError(cx, err)
	cx.AbortWithStatusJSON(status, hal.NewApiProblem(status, debugError(err)))
}

func RespondServerError(cx *gin.Context, err error) {
	status := http.StatusInternalServerError
	logHandlerError(cx, err)
	cx.AbortWithStatusJSON(status, hal.NewApiProblem(status, debugError(err)))
}
