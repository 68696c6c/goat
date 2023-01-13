package controller

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query2"
	"github.com/68696c6c/goat/repo"
	"github.com/68696c6c/goat/resource"
)

type CRUD interface {
	// Controller
	Creator
	Lister
	Viewer
	Updater
	Deleter
}

// type CRUD[T any] interface {
// 	Create(cx *gin.Context) (T, error, goat.ErrorResponder)
// 	List(cx *gin.Context) (resource.Collection[[]T], error, goat.ErrorResponder)
// 	View(cx *gin.Context) (T, error, goat.ErrorResponder)
// 	Update(cx *gin.Context) (T, error, goat.ErrorResponder)
// 	Delete(cx *gin.Context) (error, goat.ErrorResponder)
// }

// type Controller interface {
// 	GetErrorHandler() goat.ErrorHandler
// }

type Creator interface {
	Create(cx *gin.Context)
}

type Lister interface {
	// GetBaseUrl() *url.URL
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

// type controllerList interface {
// 	Controller
// 	Lister
// }

func HandleList[M any](cx *gin.Context, r repo.Filterer[M], baseUrl *url.URL, q query2.Builder, p resource.Pagination) {
	e := goat.InitErrorHandler()
	// q := query.NewQueryBuilder(cx)
	// if len(baseQuery) > 0 {
	// 	q = baseQuery[0]
	// }

	resources, pagination, err := r.Filter(cx.Request.Context(), q, p)
	if err != nil {
		e.HandleError(cx, errors.Wrap(err, "failed to list resources"), goat.RespondServerError)
		return
	}

	goat.RespondOk(cx, resource.MakeCollection[M](resources, pagination, baseUrl))
}

// func List[M any](cx *gin.Context, r repo.Filterer[M], baseUrl *url.URL, q query2.Builder, p resource.Pagination) (resource.Collection[[]M], error, goat.ErrorResponder) {
// 	// errorHandler := c.GetErrorHandler()
// 	// q := query.NewQueryBuilder(cx)
// 	// if len(baseQuery) > 0 {
// 	// 	q = baseQuery[0]
// 	// }
//
// 	resources, pagination, err := r.Filter(cx.Request.Context(), q, p)
// 	if err != nil {
// 		return resource.Collection[[]M]{}, errors.Wrap(err, "failed to list resources"), goat.RespondServerError
// 		// errorHandler.HandleError(cx, errors.Wrap(err, "failed to list resources"), goat.RespondServerError)
// 		// return
// 	}
// 	// println("PAGINATION: " + p.String())
//
// 	// goat.RespondOk(cx, resource.MakeCollection[M](resources, q.GetPagination(), c.GetBaseUrl()))
// 	return resource.MakeCollection[M](resources, pagination, baseUrl), nil, nil
// }

// func DoList[M any](cx *gin.Context, r repo.Filterer[M], baseUrl *url.URL, filters ...repo.QueryFilter) (resource.Collection[[]M], error, goat.ErrorResponder) {
// 	pagination := resource.NewPaginationFromQuery(cx.Request.URL.Query())
//
// 	resources, p, err := r.Filter(cx.Request.Context(), pagination, filters...)
// 	if err != nil {
// 		return resource.Collection[[]M]{}, errors.Wrap(err, "failed to list resources"), goat.RespondServerError
// 		// errorHandler.HandleError(cx, errors.Wrap(err, "failed to list resources"), goat.RespondServerError)
// 		// return
// 	}
//
// 	// goat.RespondOk(cx, resource.MakeCollection[M](resources, q.GetPagination(), c.GetBaseUrl()))
// 	return resource.MakeCollection[M](resources, p, baseUrl), nil, nil
// }

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

// func View[M any](cx *gin.Context, r repo.Identifier[M]) (M, error, goat.ErrorResponder) {
// 	// errorHandler := c.GetErrorHandler()
// 	var temp M
//
// 	id, err := goat.GetIdParam(cx)
// 	if err != nil {
// 		return temp, err, goat.RespondBadRequest
// 		// errorHandler.HandleError(cx, err, goat.RespondBadRequest)
// 		// return
// 	}
//
// 	m, err := r.GetById(cx.Request.Context(), id, true)
// 	if err != nil {
// 		if goat.RecordNotFound(err) {
// 			return temp, errors.New("resource does not exist"), goat.RespondNotFound
// 			// errorHandler.HandleError(cx, errors.New("resource does not exist"), goat.RespondNotFound)
// 			// return
// 		} else {
// 			return temp, errors.New("failed to get resource"), goat.RespondServerError
// 			// errorHandler.HandleError(cx, errors.Wrap(err, "failed to get resource"), goat.RespondServerError)
// 			// return
// 		}
// 	}
//
// 	return m, nil, nil
// 	// goat.RespondOk(cx, m)
// }

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

// func Create[M any, U any](cx *gin.Context, r repoCreate[M, U]) (M, error, goat.ErrorResponder) {
// 	// errorHandler := c.GetErrorHandler()
// 	var result M
// 	ctx := cx.Request.Context()
//
// 	// req, err := goat.GetRequest[U](cx)
// 	// if err != nil {
// 	// 	return result, err, goat.RespondBadRequest
// 	// 	// errorHandler.HandleError(cx, err, goat.RespondBadRequest)
// 	// 	// return
// 	// }
// 	req, err := goat.BindRequest[U](cx)
// 	if err != nil {
// 		return result, err, goat.RespondBadRequest
// 	}
//
// 	m, err := r.Create(ctx, req)
// 	if err != nil {
// 		return result, err, goat.RespondValidationError
// 		// errorHandler.HandleError(cx, err, goat.RespondValidationError)
// 		// return
// 	}
//
// 	err = r.Save(ctx, m)
// 	if err != nil {
// 		return result, errors.Wrap(err, "failed to save resource"), goat.RespondServerError
// 		// errorHandler.HandleError(cx, errors.Wrap(err, "failed to save resource"), goat.RespondServerError)
// 		// return
// 	}
//
// 	// goat.RespondCreated(cx, m)
// 	return m, nil, nil
// }

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

// func Update[M any, U any](cx *gin.Context, r repoUpdate[M, U]) (M, error, goat.ErrorResponder) {
// 	// errorHandler := c.GetErrorHandler()
//
// 	var temp M
// 	ctx := cx.Request.Context()
//
// 	id, err := goat.GetIdParam(cx)
// 	if err != nil {
// 		return temp, err, goat.RespondBadRequest
// 		// errorHandler.HandleError(cx, err, goat.RespondBadRequest)
// 		// return
// 	}
//
// 	req, err := goat.GetRequest[U](cx)
// 	if err != nil {
// 		return temp, err, goat.RespondBadRequest
// 		// errorHandler.HandleError(cx, err, goat.RespondBadRequest)
// 		// return
// 	}
//
// 	m, err := r.Update(ctx, id, req)
// 	if err != nil {
// 		if goat.RecordNotFound(err) {
// 			return temp, errors.New("resource does not exist"), goat.RespondNotFound
// 			// errorHandler.HandleError(cx, errors.New("resource does not exist"), goat.RespondNotFound)
// 			// return
// 		} else {
// 			// TODO: this could also be a server error; we need a way to differentiate it
// 			return temp, err, goat.RespondValidationError
// 			// errorHandler.HandleError(cx, err, goat.RespondValidationError)
// 			// return
// 		}
// 	}
//
// 	err = r.Save(ctx, m)
// 	if err != nil {
// 		return temp, errors.Wrap(err, "failed to save resource"), goat.RespondServerError
// 		// errorHandler.HandleError(cx, errors.Wrap(err, "failed to save resource"), goat.RespondServerError)
// 		// return
// 	}
//
// 	return m, nil, nil
// 	// goat.RespondUsed(cx, m)
// }

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

// func Delete[M any](cx *gin.Context, r repoDelete[M]) (error, goat.ErrorResponder) {
// 	// errorHandler := c.GetErrorHandler()
//
// 	ctx := cx.Request.Context()
//
// 	id, err := goat.GetIdParam(cx)
// 	if err != nil {
// 		return err, goat.RespondBadRequest
// 		// errorHandler.HandleError(cx, err, goat.RespondBadRequest)
// 		// return
// 	}
//
// 	m, err := r.GetById(ctx, id)
// 	if err != nil {
// 		if goat.RecordNotFound(err) {
// 			return errors.New("resource does not exist"), goat.RespondNotFound
// 			// errorHandler.HandleError(cx, errors.New("resource does not exist"), goat.RespondNotFound)
// 			// return
// 		} else {
// 			return errors.Wrap(err, "failed to load resource"), goat.RespondServerError
// 			// errorHandler.HandleError(cx, errors.Wrap(err, "failed to load resource"), goat.RespondServerError)
// 			// return
// 		}
// 	}
//
// 	err = r.Delete(ctx, m)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to delete resource"), goat.RespondServerError
// 		// errorHandler.HandleError(cx, errors.Wrap(err, "failed to delete resource"), goat.RespondServerError)
// 		// return
// 	}
//
// 	// goat.RespondValid(cx)
// 	return nil, nil
// }
