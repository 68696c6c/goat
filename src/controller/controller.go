package controller

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/repo"
	"github.com/68696c6c/goat/resource"
)

type CRUD interface {
	Creator
	Lister
	Viewer
	Updater
	Deleter
}

type Creator interface {
	Create(cx *gin.Context)
}

type Lister interface {
	List(cx *gin.Context)
}

type Viewer interface {
	View(cx *gin.Context)
}

type Updater interface {
	Update(cx *gin.Context)
}

type Deleter interface {
	Delete(cx *gin.Context)
}

func HandleList[M any](cx *gin.Context, r repo.Filterer[M], baseUrl *url.URL, q query.Builder, p resource.Pagination) {
	e := goat.InitErrorHandler()

	resources, pagination, err := r.Filter(cx.Request.Context(), q, p)
	if err != nil {
		e.HandleError(cx, errors.Wrap(err, "failed to list resources"), goat.RespondServerError)
		return
	}

	goat.RespondOk(cx, resource.MakeCollection[M](resources, pagination, baseUrl))
}

func HandleView[M any](cx *gin.Context, r repo.Identifier[M]) {
	e := goat.InitErrorHandler()

	id, err := goat.GetIdParam(cx)
	if err != nil {
		e.HandleError(cx, err, goat.RespondBadRequest)
		return
	}

	m, err := r.GetById(cx.Request.Context(), id, true)
	if err != nil {
		if goat.RecordNotFound(err) {
			e.HandleError(cx, errors.New("resource does not exist"), goat.RespondNotFound)
			return
		} else {
			e.HandleError(cx, errors.Wrap(err, "failed to get resource"), goat.RespondServerError)
			return
		}
	}

	goat.RespondOk(cx, m)
}

type repoCreate[Model, Request any] interface {
	repo.Identifier[Model]
	repo.Saver[Model]
	repo.Creator[Model, Request]
}

func HandleCreate[M any, U any](cx *gin.Context, r repoCreate[M, U]) {
	e := goat.InitErrorHandler()

	ctx := cx.Request.Context()

	req, err := goat.GetRequest[U](cx)
	if err != nil {
		e.HandleError(cx, err, goat.RespondBadRequest)
		return
	}

	m, err := r.Create(ctx, req)
	if err != nil {
		e.HandleError(cx, err, goat.RespondValidationError)
		return
	}

	err = r.Save(ctx, m)
	if err != nil {
		e.HandleError(cx, errors.Wrap(err, "failed to save resource"), goat.RespondServerError)
		return
	}

	goat.RespondCreated(cx, m)
}

type repoUpdate[M any, U any] interface {
	repo.Identifier[M]
	repo.Saver[M]
	repo.Updater[M, U]
}

func HandleUpdate[M any, U any](cx *gin.Context, r repoUpdate[M, U]) {
	e := goat.InitErrorHandler()
	ctx := cx.Request.Context()

	id, err := goat.GetIdParam(cx)
	if err != nil {
		e.HandleError(cx, err, goat.RespondBadRequest)
		return
	}

	req, err := goat.GetRequest[U](cx)
	if err != nil {
		e.HandleError(cx, err, goat.RespondBadRequest)
		return
	}

	m, err := r.Update(ctx, id, req)
	if err != nil {
		if goat.RecordNotFound(err) {
			e.HandleError(cx, errors.New("resource does not exist"), goat.RespondNotFound)
			return
		} else {
			e.HandleError(cx, err, goat.RespondValidationError)
			return
		}
	}

	err = r.Save(ctx, m)
	if err != nil {
		e.HandleError(cx, errors.Wrap(err, "failed to save resource"), goat.RespondServerError)
		return
	}

	goat.RespondOk(cx, m)
}

type repoDelete[M any] interface {
	repo.Identifier[M]
	repo.Deleter[M]
}

func HandleDelete[M any](cx *gin.Context, r repoDelete[M]) {
	e := goat.InitErrorHandler()
	ctx := cx.Request.Context()

	id, err := goat.GetIdParam(cx)
	if err != nil {
		e.HandleError(cx, err, goat.RespondBadRequest)
		return
	}

	m, err := r.GetById(ctx, id)
	if err != nil {
		if goat.RecordNotFound(err) {
			e.HandleError(cx, errors.New("resource does not exist"), goat.RespondNotFound)
			return
		} else {
			e.HandleError(cx, errors.Wrap(err, "failed to load resource"), goat.RespondServerError)
			return
		}
	}

	err = r.Delete(ctx, m)
	if err != nil {
		e.HandleError(cx, errors.Wrap(err, "failed to delete resource"), goat.RespondServerError)
		return
	}

	goat.RespondNoContent(cx)
}
