package goat

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/68696c6c/goat/resource"
)

func debugError(err error) error {
	if DebugEnabled() {
		return err
	}
	return nil
}

// func RespondValid(cx *gin.Context) {
// 	// g.HTTP.RespondValid(cx)
// 	cx.AbortWithStatusJSON(http.StatusOK, true)
// }

func RespondOk(cx *gin.Context, data any) {
	// g.HTTP.RespondOk(cx, data)
	cx.AbortWithStatusJSON(http.StatusOK, data)
}

func RespondNoContent(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusNoContent)
}

// func RespondUsed(cx *gin.Context, data any) {
// 	// g.HTTP.RespondUsed(cx, data)
// 	cx.AbortWithStatusJSON(http.StatusIMUsed, data)
// }

func RespondAccepted(cx *gin.Context, data any) {
	// g.HTTP.RespondAccepted(cx, data)
	cx.AbortWithStatusJSON(http.StatusAccepted, data)
}

func RespondCreated(cx *gin.Context, data any) {
	// g.HTTP.RespondCreated(cx, data)
	cx.AbortWithStatusJSON(http.StatusCreated, data)
}

func RespondNotFound(cx *gin.Context, err error) {
	// g.HTTP.RespondNotFound(cx, err)
	status := http.StatusNotFound
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "not found", debugError(err)))
}

func RespondBadRequest(cx *gin.Context, err error) {
	// g.HTTP.RespondBadRequest(cx, err)
	status := http.StatusBadRequest
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "invalid request", debugError(err)))
}

func RespondValidationError(cx *gin.Context, err error) {
	// g.HTTP.RespondValidationError(cx, err)
	status := http.StatusBadRequest
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "invalid request", err))
}

func RespondUnauthorized(cx *gin.Context, err error) {
	// g.HTTP.RespondUnauthorized(cx, err)
	status := http.StatusUnauthorized
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "unauthorized", debugError(err)))
}

func RespondForbidden(cx *gin.Context, err error) {
	// g.HTTP.RespondForbidden(cx, err)
	status := http.StatusForbidden
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "authentication error", debugError(err)))
}

func RespondServerError(cx *gin.Context, err error) {
	// g.HTTP.RespondServerError(cx, err)
	status := http.StatusInternalServerError
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "internal server error", debugError(err)))
}

type MessageResponse struct {
	Message string `json:"message"`
	resource.Resource
}

func HealthHandler(path string, message ...string) (string, gin.HandlerFunc) {
	msg := "ok"
	if len(message) > 0 {
		msg = strings.Join(message, "")
	}
	return path, func(cx *gin.Context) {
		RespondOk(cx, MessageResponse{
			Message:  msg,
			Resource: resource.MakeResource(GetUrl().JoinPath(path).String()),
		})
	}
}

func VersionHandler(path string, version ...string) (string, gin.HandlerFunc) {
	msg := "dev"
	if len(version) > 0 {
		msg = strings.Join(version, "")
	}
	return path, func(cx *gin.Context) {
		RespondOk(cx, MessageResponse{
			Message:  msg,
			Resource: resource.MakeResource(GetUrl().JoinPath(path).String()),
		})
	}
}
