package middlewares

import (
	"github.com/68696c6c/goat"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/68696c6c/web/app/lib/auth"
	"github.com/68696c6c/web/app/repos"
	"github.com/68696c6c/web/utils"
)

func ValidateJwt(authService auth.Service, usersRepo repos.UsersRepo) gin.HandlerFunc {
	return func(cx *gin.Context) {
		accessJwt, err := authService.ValidateToken(cx.Request)
		if err != nil {
			goat.RespondUnauthorized(cx, err)
			return
		}
		token := accessJwt.GetAccess()
		cx.Set(utils.ContextKeyCurrentUserToken, token)

		claims, err := authService.GetTokenClaims(token)
		if err != nil {
			goat.RespondServerError(cx, errors.Wrap(err, "failed to parse jwt claims"))
			return
		}

		user, err := usersRepo.GetByEmail(cx, claims.Subject)
		if err != nil {
			goat.RespondUnauthorized(cx, errors.Wrap(err, "failed to load jwt user"))
			return
		}
		cx.Set(utils.ContextKeyCurrentUser, user)
	}
}
