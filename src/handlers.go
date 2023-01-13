package goat

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/68696c6c/goat/resource"
)

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
