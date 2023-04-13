package utils

import (
	"github.com/gin-gonic/gin"

	"github.com/68696c6c/web/app/models"
)

const (
	ContextKeyCurrentUser      = "currentUser"
	ContextKeyCurrentUserToken = "currentUserToken"
)

func GetCurrentUser(cx *gin.Context) (*models.User, bool) {
	value, exists := cx.Get(ContextKeyCurrentUser)
	return value.(*models.User), exists
}

func GetCurrentUserToken(cx *gin.Context) (string, bool) {
	value, exists := cx.Get(ContextKeyCurrentUserToken)
	return value.(string), exists
}
