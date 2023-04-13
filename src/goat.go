package goat

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/68696c6c/goat/hal"
	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/sys"
	"github.com/68696c6c/goat/sys/database"
	"github.com/68696c6c/goat/sys/http"
)

var g sys.Goat

// Init initializes the Goat runtime services.
// Goat has three primary concerns, each encapsulated by their own service:
// - logging
// - database connections
// - route-based response hypermedia (linking)
func Init() error {
	if g != (sys.Goat{}) {
		return nil
	}
	var err error
	g, err = sys.Init()
	if err != nil {
		return err
	}
	return nil
}

func MustInit() {
	err := Init()
	if err != nil {
		panic(err)
	}
}

type Router http.Router

func InitRouter() Router {
	return g.HTTP.InitRouter()
}

func GetLogger() *zap.SugaredLogger {
	return g.Log.Logger()
}

func GetStrictLogger() *zap.Logger {
	return g.Log.StrictLogger()
}

func GetUrl(key ...string) *url.URL {
	return g.HTTP.GetUrl(key...)
}

func NewResourceLinks(key, path string) *hal.ResourceLinks {
	return hal.NewResourceLinks(g.HTTP.GetUrl(key).JoinPath(path).String())
}

func GetMainDB() (*gorm.DB, error) {
	return g.DB.GetMainDB()
}

type DatabaseConfig database.Config

func GetDB(c DatabaseConfig) (*gorm.DB, error) {
	return g.DB.GetConnection(database.Config(c))
}

func ApplyQueryToGorm(db *gorm.DB, q query.Builder, paginate bool) {
	t := q.Build()
	if t.Where != "" {
		db = db.Where(t.Where, t.Params...)
	}
	for _, p := range t.Joins {
		db = db.Preload(p.Query, p.Args...)
	}
	if t.OrderBy != "" {
		db = db.Order(t.OrderBy)
	}
	if paginate {
		if t.Limit > 0 {
			db = db.Limit(t.Limit)
		}
		if t.Offset > 0 {
			db = db.Offset(t.Offset)
		}
	}
}

// BindRequest returns a T with values set by binding the request JSON from the provided Gin context.
func BindRequest[T any](cx *gin.Context) (T, error) {
	var result T
	err := cx.ShouldBindWith(&result, binding.JSON)
	if err != nil {
		return result, err
	}
	return result, nil
}

type ParamParser[T any] func(string) (T, error)

func ParseParam[T any](cx *gin.Context, key string, parser ParamParser[T]) (T, error) {
	param := cx.Param(key)
	result, err := parser(param)
	if err != nil {
		return result, errors.Wrapf(err, "failed to parse request param '%s' from value '%s'", key, param)
	}
	return result, nil
}

func GetIDParam(cx *gin.Context) (ID, error) {
	return ParseParam[ID](cx, "id", ParseID)
}
