package controllers

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/controller"
	"github.com/68696c6c/goat/query2"
	"github.com/68696c6c/goat/resource"
	"github.com/gin-gonic/gin"

	"github.com/68696c6c/example/app/models"
	"github.com/68696c6c/example/app/repos"
)

type OrganizationsController interface {
	controller.CRUD
}

type organizations struct {
	// baseUrl string
	// baseUrl *url.URL
	// errors goat.ErrorHandler
	repo repos.OrganizationsRepo
}

func NewOrganizationsController(
	// baseUrl *url.URL,
	// baseUrl string,
	// errors goat.ErrorHandler,
	repo repos.OrganizationsRepo,
) OrganizationsController {
	return organizations{
		// baseUrl: baseUrl,
		// baseUrl: baseUrl,
		// errors: errors,
		repo: repo,
	}
}

// func (c organizations) GetErrorHandler() goat.ErrorHandler {
// 	return c.errors
// }

// func (c organizations) GetLink(path string) resource.Link {
// 	println("ORGANIZATIONS GET LINK")
// 	return goat.GetLink(models.OrganizationLinkKey, path)
// 	// return resource.MakeLink(c.baseUrl.JoinPath(path).String())
// }

// func (c organizations) GetBaseUrl() string {
// 	return c.baseUrl
// }

// func (c organizations) GetBaseUrl() *url.URL {
// 	return goat.GetUrl(models.OrganizationLinkKey)
// }

func (c organizations) List(cx *gin.Context) {
	// controller.HandleList[*models.Organization](cx, c, c.repo)
	queryString := cx.Request.URL.Query()
	q := query2.NewQueryFromUrl(queryString)
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
