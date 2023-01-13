package controllers

import (
	"github.com/68696c6c/goat"
	"github.com/gin-gonic/gin"

	"github.com/68696c6c/example/app/enums"
	"github.com/68696c6c/example/app/models"
)

const userKey = "currentUser"

func SetCurrentUser() gin.HandlerFunc {
	return func(cx *gin.Context) {
		user := models.MakeUser()
		// user.Level = enums.UserLevelSuper
		orgId, _ := goat.ParseID("41be9141-3d11-4133-8464-1f4fcfe43225")
		user.OrganizationId = orgId
		// user.Level = enums.UserLevelAdmin
		user.Level = enums.UserLevelSuper
		cx.Set(userKey, user)
	}
}

func getCurrentUser(cx *gin.Context) (*models.User, bool) {
	value, exists := cx.Get(userKey)
	return value.(*models.User), exists
}
