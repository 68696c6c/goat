package controllers

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/controller"
	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/resource"
	"github.com/gin-gonic/gin"

	"github.com/68696c6c/example/app/models"
	"github.com/68696c6c/example/app/repos"
)

type OrganizationsController interface {
	controller.CRUD
}

type organizations struct {
	repo repos.OrganizationsRepo
}

func NewOrganizationsController(repo repos.OrganizationsRepo) OrganizationsController {
	return organizations{
		repo: repo,
	}
}

func (c organizations) List(cx *gin.Context) {
	queryString := cx.Request.URL.Query()
	q := query.NewQueryFromUrl(queryString)
	p := resource.NewPaginationFromUrl(queryString)
	controller.HandleList[*models.Organization](cx, c.repo, goat.GetUrl(models.OrganizationLinkKey), q, p)
}

func (c organizations) View(cx *gin.Context) {
	controller.HandleView[*models.Organization](cx, c.repo)
}

func (c organizations) Create(cx *gin.Context) {
	controller.HandleCreate[*models.Organization, models.OrganizationRequest](cx, c.repo)
}

func (c organizations) Update(cx *gin.Context) {
	controller.HandleUpdate[*models.Organization, models.OrganizationRequest](cx, c.repo)
}

func (c organizations) Delete(cx *gin.Context) {
	controller.HandleDelete[*models.Organization](cx, c.repo)
}
