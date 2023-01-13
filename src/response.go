package goat

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/68696c6c/goat/resource"
)

type ApiProblem resource.ApiProblem

type Resource resource.Resource

type Collection[T any] resource.Collection[T]

type MessageResponse resource.MessageResponse

func debugError(err error) error {
	if DebugEnabled() {
		return err
	}
	return nil
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
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "not found", debugError(err)))
}

func RespondBadRequest(cx *gin.Context, err error) {
	status := http.StatusBadRequest
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "invalid request", debugError(err)))
}

func RespondValidationError(cx *gin.Context, err error) {
	status := http.StatusBadRequest
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "invalid request", err))
}

func RespondUnauthorized(cx *gin.Context, err error) {
	status := http.StatusUnauthorized
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "unauthorized", debugError(err)))
}

func RespondForbidden(cx *gin.Context, err error) {
	status := http.StatusForbidden
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "authentication error", debugError(err)))
}

func RespondServerError(cx *gin.Context, err error) {
	status := http.StatusInternalServerError
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "internal server error", debugError(err)))
}
