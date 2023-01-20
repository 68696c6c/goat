package controllers

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/controller"
	"github.com/68696c6c/goat/query"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/68696c6c/example/app/models"
	"github.com/68696c6c/example/app/repos"
)

type UsersController interface {
	controller.CRUD
}

type users struct {
	repo repos.UsersRepo
}

func NewUsersController(repo repos.UsersRepo) UsersController {
	return users{
		repo: repo,
	}
}

func (c users) List(cx *gin.Context) {
	queryString := cx.Request.URL.Query()
	q := query.NewQueryFromUrl(queryString)

	currentUser, ok := getCurrentUser(cx)
	if !ok || currentUser == nil {
		goat.RespondUnauthorized(cx, errors.New("login required"))
		return
	}
	err := c.repo.ApplyFilterForUser(q, currentUser)
	if err != nil {
		goat.RespondServerError(cx, err)
		return
	}
	err = c.repo.FilterStrings(q, queryString)
	if err != nil {
		goat.RespondServerError(cx, err)
		return
	}
	controller.HandleList[*models.User](cx, c.repo, goat.GetUrl(models.UserLinkKey), q)
}

func (c users) View(cx *gin.Context) {
	controller.HandleView[*models.User](cx, c.repo)
}

func (c users) Create(cx *gin.Context) {
	controller.HandleCreate[*models.User, models.UserRequest](cx, c.repo)
}

func (c users) Update(cx *gin.Context) {
	controller.HandleUpdate[*models.User, models.UserRequest](cx, c.repo)
}

func (c users) Delete(cx *gin.Context) {
	controller.HandleDelete[*models.User](cx, c.repo)
}
