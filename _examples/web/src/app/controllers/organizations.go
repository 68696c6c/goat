package controllers

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/controller"
	"github.com/68696c6c/goat/query"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/68696c6c/web/app/lib/auth"
	"github.com/68696c6c/web/app/models"
	"github.com/68696c6c/web/app/repos"
)

type OrganizationsController interface {
	controller.CRUD
}

type organizations struct {
	repo repos.OrganizationsRepo
	auth auth.Service
}

func NewOrganizationsController(repo repos.OrganizationsRepo, authService auth.Service) OrganizationsController {
	return organizations{
		repo: repo,
		auth: authService,
	}
}

func (c organizations) List(cx *gin.Context) {
	queryString := cx.Request.URL.Query()
	q := query.NewQueryFromUrl(queryString)

	currentUser, claims, err := c.auth.GetCurrentUserAndClaims(cx)
	if err != nil {
		goat.RespondBadRequest(cx, err)
		return
	}

	repos.FilterOrganizationsForUser(q, currentUser, claims)
	repos.FilterOrganizationsQuery(q, queryString)

	controller.HandleList[*models.Organization](cx, c.repo, goat.GetUrl(models.OrganizationLinkKey), q)
}

func (c organizations) View(cx *gin.Context) {
	currentUser, claims, err := c.auth.GetCurrentUserAndClaims(cx)
	if err != nil {
		goat.RespondBadRequest(cx, err)
		return
	}
	controller.HandleView[*models.Organization](cx, c.repo, func(m *models.Organization) error {
		if repos.UserHasOrganizationAccess(m, currentUser, claims, auth.Read) {
			return nil
		}
		return errors.New("access denied")
	})
}

func (c organizations) Create(cx *gin.Context) {
	currentUser, claims, err := c.auth.GetCurrentUserAndClaims(cx)
	if err != nil {
		goat.RespondBadRequest(cx, err)
		return
	}
	controller.HandleCreate[*models.Organization, models.OrganizationRequest](cx, c.repo, func(m *models.Organization) error {
		if repos.UserHasOrganizationAccess(m, currentUser, claims, auth.Create) {
			return nil
		}
		return errors.New("access denied")
	})
}

func (c organizations) Update(cx *gin.Context) {
	currentUser, claims, err := c.auth.GetCurrentUserAndClaims(cx)
	if err != nil {
		goat.RespondBadRequest(cx, err)
		return
	}
	controller.HandleUpdate[*models.Organization, models.OrganizationRequest](cx, c.repo, func(m *models.Organization) error {
		if repos.UserHasOrganizationAccess(m, currentUser, claims, auth.Update) {
			return nil
		}
		return errors.New("access denied")
	})
}

func (c organizations) Delete(cx *gin.Context) {
	currentUser, claims, err := c.auth.GetCurrentUserAndClaims(cx)
	if err != nil {
		goat.RespondBadRequest(cx, err)
		return
	}
	controller.HandleDelete[*models.Organization](cx, c.repo, func(m *models.Organization) error {
		if repos.UserHasOrganizationAccess(m, currentUser, claims, auth.Delete) {
			return nil
		}
		return errors.New("access denied")
	})
}
