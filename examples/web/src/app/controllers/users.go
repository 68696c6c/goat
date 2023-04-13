package controllers

import (
	"errors"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/controller"
	"github.com/68696c6c/goat/query"
	"github.com/gin-gonic/gin"

	"github.com/68696c6c/web/app/lib/auth"
	"github.com/68696c6c/web/app/models"
	"github.com/68696c6c/web/app/repos"
)

type UsersController interface {
	controller.CRUD
}

type users struct {
	repo repos.UsersRepo
	auth auth.Service
}

func NewUsersController(repo repos.UsersRepo, authService auth.Service) UsersController {
	return users{
		repo: repo,
		auth: authService,
	}
}

func (c users) List(cx *gin.Context) {
	queryString := cx.Request.URL.Query()
	q := query.NewQueryFromUrl(queryString)

	currentUser, claims, err := c.auth.GetCurrentUserAndClaims(cx)
	if err != nil {
		goat.RespondBadRequest(cx, err)
		return
	}

	repos.FilterUsersForUser(q, currentUser, claims)
	repos.FilterUsersQuery(q, queryString)

	controller.HandleList[*models.User](cx, c.repo, goat.GetUrl(models.UserLinkKey), q)
}

func (c users) View(cx *gin.Context) {
	currentUser, claims, err := c.auth.GetCurrentUserAndClaims(cx)
	if err != nil {
		goat.RespondBadRequest(cx, err)
		return
	}
	controller.HandleView[*models.User](cx, c.repo, func(m *models.User) error {
		if repos.UserHasUserAccess(m, currentUser, claims, auth.Read) {
			return nil
		}
		return errors.New("access denied")
	})
}

func (c users) Create(cx *gin.Context) {
	currentUser, claims, err := c.auth.GetCurrentUserAndClaims(cx)
	if err != nil {
		goat.RespondBadRequest(cx, err)
		return
	}
	controller.HandleCreate[*models.User, models.UserRequest](cx, c.repo, func(m *models.User) error {
		if repos.UserHasUserAccess(m, currentUser, claims, auth.Create) && repos.CanUserWriteUserLevel(m.Level, currentUser) {
			return nil
		}
		return errors.New("access denied")
	})
}

func (c users) Update(cx *gin.Context) {
	currentUser, claims, err := c.auth.GetCurrentUserAndClaims(cx)
	if err != nil {
		goat.RespondBadRequest(cx, err)
		return
	}
	controller.HandleUpdate[*models.User, models.UserRequest](cx, c.repo, func(m *models.User) error {
		if repos.UserHasUserAccess(m, currentUser, claims, auth.Update) && repos.CanUserWriteUserLevel(m.Level, currentUser) {
			return nil
		}
		return errors.New("access denied")
	})
}

func (c users) Delete(cx *gin.Context) {
	currentUser, claims, err := c.auth.GetCurrentUserAndClaims(cx)
	if err != nil {
		goat.RespondBadRequest(cx, err)
		return
	}
	controller.HandleDelete[*models.User](cx, c.repo, func(m *models.User) error {
		if repos.UserHasUserAccess(m, currentUser, claims, auth.Delete) {
			return nil
		}
		return errors.New("access denied")
	})
}
