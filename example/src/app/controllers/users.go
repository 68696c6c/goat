package controllers

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/controller"
	"github.com/68696c6c/goat/query2"
	"github.com/68696c6c/goat/resource"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/68696c6c/example/app/models"
	"github.com/68696c6c/example/app/repos"
)

type UsersController interface {
	controller.CRUD
	// controller.CRUD[*models.User]
}

type users struct {
	// errors goat.ErrorHandler
	repo repos.UsersRepo
}

func NewUsersController(
	// errors goat.ErrorHandler,
	repo repos.UsersRepo,
) UsersController {
	return users{
		repo: repo,
		// errors: errors,
	}
}

// func (c users) GetErrorHandler() goat.ErrorHandler {
// 	return c.errors
// }
//
// func (c users) GetBaseUrl() *url.URL {
// 	return goat.GetUrl(models.UserLinkKey)
// }

// func (c users) List(cx *gin.Context) {
// 	f, err := goat.GetFilter(cx)
// 	if err != nil {
// 		c.errors.HandleError(cx, err, goat.RespondServerError)
// 		return
// 	}
// 	currentUser, ok := getCurrentUser(cx)
// 	if !ok {
// 		c.errors.HandleError(cx, errors.New("login required"), goat.RespondUnauthorized)
// 		return
// 	}
// 	err = c.repo.ApplyFilterForUser(f, currentUser)
// 	if err != nil {
// 		c.errors.HandleError(cx, err, goat.RespondServerError)
// 		return
// 	}
// 	controller.HandleList[*models.User](cx, c, c.repo, f)
// }
//
// func (c users) View(cx *gin.Context) {
// 	controller.HandleView[*models.User](cx, c, c.repo)
// }
//
// func (c users) Create(cx *gin.Context) {
// 	controller.HandleCreate[*models.User, models.UserRequest](cx, c, c.repo)
// }
//
// func (c users) Update(cx *gin.Context) {
// 	controller.HandleUpdate[*models.User, models.UserRequest](cx, c, c.repo)
// }
//
// func (c users) Delete(cx *gin.Context) {
// 	controller.HandleDelete[*models.User](cx, c, c.repo)
// }

func (c users) List(cx *gin.Context) {
	e := goat.InitErrorHandler()
	var err error
	queryString := cx.Request.URL.Query()
	q := query2.NewQueryFromUrl(queryString)
	p := resource.NewPaginationFromUrl(queryString)
	// f, err := goat.GetFilter(cx)
	// if err != nil {
	// 	return resource.Collection[[]*models.User]{}, err, goat.RespondServerError
	// 	// c.errors.HandleError(cx, err, goat.RespondServerError)
	// 	// return
	// }
	currentUser, ok := getCurrentUser(cx)
	if !ok || currentUser == nil {
		// return resource.Collection[[]*models.User]{}, errors.New("login required"), goat.RespondUnauthorized
		e.HandleError(cx, errors.New("login required"), goat.RespondUnauthorized)
		return
	}
	err = c.repo.ApplyFilterForUser(q, currentUser)
	if err != nil {
		// return resource.Collection[[]*models.User]{}, err, goat.RespondServerError
		e.HandleError(cx, err, goat.RespondServerError)
		return
	}
	err = c.repo.FilterStrings(q, queryString)
	if err != nil {
		// return resource.Collection[[]*models.User]{}, err, goat.RespondServerError
		e.HandleError(cx, err, goat.RespondServerError)
		return
	}
	controller.HandleList[*models.User](cx, c.repo, goat.GetUrl(models.UserLinkKey), q, p)
	// var err error
	//
	// currentUser, ok := getCurrentUser(cx)
	// if !ok || currentUser == nil {
	// 	return resource.Collection[[]*models.User]{}, errors.New("login required"), goat.RespondUnauthorized
	// }
	//
	// pagination := resource.NewPaginationFromQuery(cx.Request.URL.Query())
	//
	// resources, p, err := c.repo.Filter(cx.Request.Context(), pagination, c.repo.UserLevelFilter(*currentUser))
	// if err != nil {
	// 	return resource.Collection[[]*models.User]{}, errors.Wrap(err, "failed to list resources"), goat.RespondServerError
	// 	// errorHandler.HandleError(cx, errors.Wrap(err, "failed to list resources"), goat.RespondServerError)
	// 	// return
	// }
	//
	// // goat.RespondOk(cx, resource.MakeCollection[M](resources, q.GetPagination(), c.GetBaseUrl()))
	// collection := resource.MakeCollection[*models.User](resources, p, goat.GetUrl(models.UserLinkKey))
	// return collection, nil, nil
	//
	// return controller.DoList[*models.User](cx, c.repo, goat.GetUrl(models.UserLinkKey), c.repo.UserLevelFilter(*currentUser), c.repo.SearchFilter(cx.Request.URL.Query()))
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
