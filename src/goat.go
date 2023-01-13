package goat

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/68696c6c/goat/query2"
	"github.com/68696c6c/goat/resource"
	"github.com/68696c6c/goat/sys"
	"github.com/68696c6c/goat/sys/http/router"
)

// type ResourceLinks map[string]*url.URL

var g sys.Goat

// var _links = make(ResourceLinks)

// Init initializes the Goat runtime services.
// Goat has four primary concerns:
// - logging
// - database connections
// - request handling
// - response hypermedia (linking)
// These concerns are encapsulated inside of services that are bootstrapped when goat.Init() is called.
func Init() {
	if g != (sys.Goat{}) {
		return
	}
	g = sys.Init()
}

// Global functions for calling encapsulated services.

type Router router.Router

func InitRouter() Router {
	return g.Router.InitRouter()
}

func GetLogger() *logrus.Logger {
	return g.Log.NewLogger()
}

func GetFileLogger(name string) (*logrus.Logger, error) {
	return g.Log.NewFileLogger(name)
}

func DebugEnabled() bool {
	// return g.HTTP.DebugEnabled()
	return g.HttpDebug
}

// func FilterMiddleware() gin.HandlerFunc {
// 	return g.HTTP.FilterMiddleware()
// }
//
// func GetFilter(cx *gin.Context) *query.Query {
// 	return g.HTTP.GetFilter(cx)
// }

// This doesn't appear to be used?
// func GetHandlerContext(cx *gin.Context) context.Context {
// 	return g.HTTP.GetHandlerContext(cx)
// }

// Returns a random string that can be used as a Basic Auth token.
func GenerateToken() string {
	u := uuid.New().String()
	return strings.Replace(u, "-", "", -1)
}

func GetUrl(key ...string) *url.URL {
	return g.Router.GetUrl(key...)
}

func MakeResourceLinks(key, path string) *resource.Links {
	return resource.MakeResourceLinks(g.Router.GetUrl(key).JoinPath(path).String())
}

// TODO: find a final resting place for these

func QueryFromGin(cx *gin.Context) query2.Builder {
	if cx == nil {
		return query2.NewQuery()
	}
	return query2.NewQueryFromUrl(cx.Request.URL.Query())
}

func PaginationFromGin(cx *gin.Context) resource.Pagination {
	if cx == nil {
		return resource.NewPagination()
	}
	return resource.NewPaginationFromUrl(cx.Request.URL.Query())
}

// includes limit and offset, use for general purpose querying
func ApplyQueryToGorm(g *gorm.DB, q query2.Builder) (*gorm.DB, error) {
	// if q.Filter != nil {
	where, params, err := q.GetWhere()
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply filter")
	}

	if where != "" {
		g = g.Where(where, params...)
	}
	// }

	order := q.GetOrder()
	if order != "" {
		g = g.Order(order)
	}

	limit := q.GetLimit()
	if limit > 0 {
		g = g.Limit(limit)
	}

	offset := q.GetOffset()
	if offset > 0 {
		g = g.Offset(offset)
	}
	// page := q.Pagination.Page
	// size := q.Pagination.PageSize
	//
	// if size > 0 {
	// 	g = g.Limit(int(size)).Offset(int((page - 1) * size))
	// }

	for _, p := range q.GetPreload() {
		g = g.Preload(p)
	}

	return g, nil
}

// does not set limit or offset, use for filtering
func ApplyQueryToGormNoLimitOffset(g *gorm.DB, q query2.Builder) (*gorm.DB, error) {
	// if q.Filter != nil {
	where, params, err := q.GetWhere()
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply filter")
	}

	if where != "" {
		g = g.Where(where, params...)
	}
	// }

	order := q.GetOrder()
	if order != "" {
		g = g.Order(order)
	}

	// limit := q.GetLimit()
	// if limit > 0 {
	// 	g = g.Limit(limit)
	// }
	//
	// offset := q.GetOffset()
	// if offset > 0 {
	// 	g = g.Offset(offset)
	// }
	// // page := q.Pagination.Page
	// // size := q.Pagination.PageSize
	// //
	// // if size > 0 {
	// // 	g = g.Limit(int(size)).Offset(int((page - 1) * size))
	// // }

	for _, p := range q.GetPreload() {
		g = g.Preload(p)
	}

	return g, nil
}
