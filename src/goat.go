package goat

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/resource"
	"github.com/68696c6c/goat/sys"
	"github.com/68696c6c/goat/sys/http"
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

type Router http.Router

func InitRouter() Router {
	return g.HTTP.InitRouter()
}

func GetLogger() *zap.SugaredLogger {
	return g.Log.GetLogger()
}

func GetStrictLogger() *zap.Logger {
	return g.Log.GetStrictLogger()
}

// GenerateToken returns a random string that can be used as a Basic Auth token.
// TODO: does it really?
func GenerateToken() string {
	u := uuid.New().String()
	return strings.Replace(u, "-", "", -1)
}

func GetUrl(key ...string) *url.URL {
	return g.HTTP.GetUrl(key...)
}

func MakeResourceLinks(key, path string) *resource.Links {
	return resource.MakeResourceLinks(g.HTTP.GetUrl(key).JoinPath(path).String())
}

// TODO: find a final resting place for these

func QueryFromGin(cx *gin.Context) query.Builder {
	if cx == nil {
		return query.NewQuery()
	}
	return query.NewQueryFromUrl(cx.Request.URL.Query())
}

func PaginationFromGin(cx *gin.Context) resource.Pagination {
	if cx == nil {
		return resource.NewPagination()
	}
	return resource.NewPaginationFromUrl(cx.Request.URL.Query())
}

// includes limit and offset, use for general purpose querying
func ApplyQueryToGorm(g *gorm.DB, q query.Builder) (*gorm.DB, error) {
	where, params, err := q.GetWhere()
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply filter")
	}

	if where != "" {
		g = g.Where(where, params...)
	}

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

	for _, p := range q.GetPreload() {
		g = g.Preload(p)
	}

	return g, nil
}

// does not set limit or offset, use for filtering
func ApplyQueryToGormNoLimitOffset(g *gorm.DB, q query.Builder) (*gorm.DB, error) {
	where, params, err := q.GetWhere()
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply filter")
	}

	if where != "" {
		g = g.Where(where, params...)
	}

	order := q.GetOrder()
	if order != "" {
		g = g.Order(order)
	}

	for _, p := range q.GetPreload() {
		g = g.Preload(p)
	}

	return g, nil
}
