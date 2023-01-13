package controllers

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/controller"
	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/resource"
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
	e := goat.InitErrorHandler()
	queryString := cx.Request.URL.Query()
	q := query.NewQueryFromUrl(queryString)
	p := resource.NewPaginationFromUrl(queryString)

	currentUser, ok := getCurrentUser(cx)
	if !ok || currentUser == nil {
		e.HandleError(cx, errors.New("login required"), goat.RespondUnauthorized)
		return
	}
	err := c.repo.ApplyFilterForUser(q, currentUser)
	if err != nil {
		e.HandleError(cx, err, goat.RespondServerError)
		return
	}
	err = c.repo.FilterStrings(q, queryString)
	if err != nil {
		e.HandleError(cx, err, goat.RespondServerError)
		return
	}
	controller.HandleList[*models.User](cx, c.repo, goat.GetUrl(models.UserLinkKey), q, p)
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
