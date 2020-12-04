package goat

import (
	"context"
	"strings"

	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/src/database"
	"github.com/68696c6c/goat/src/sys"

	"github.com/68696c6c/goose"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var g sys.Goat

type Router interface {
	Run() error
	SetRegistry(d map[string]interface{})
	InitRegistry() gin.HandlerFunc
	GetEngine() *gin.Engine
}

// Goat has three primary concerns:
// - database connections and schema management,
// - request handling
// - logging
// These concerns are encapsulated inside of services that are bootstrapped when goat.Init() is called.
func Init() {
	if g != (sys.Goat{}) {
		return
	}
	g = sys.Init()
}

// Global functions for calling encapsulated services.

func GetMainDB() (*gorm.DB, error) {
	return g.DB.GetMainDB()
}

func GetMigrationDB() (*gorm.DB, error) {
	return g.DB.GetMigrationDB()
}

func GetDB(key string) (*gorm.DB, error) {
	return g.DB.GetCustomDB(key)
}

func GetCustomDB(c database.ConnectionConfig) (*gorm.DB, error) {
	return g.DB.GetConnection(c)
}

func GetSchema(connection *gorm.DB) (goose.SchemaInterface, error) {
	return g.DB.GetSchema(connection)
}

func ApplyPaginationToQuery(q *query.Query, baseGormQuery *gorm.DB) error {
	return g.DB.ApplyPaginationToQuery(q, baseGormQuery)
}

func GetRouter() Router {
	return g.HTTP.NewRouter()
}

func GetLogger() *logrus.Logger {
	return g.Log.NewLogger()
}

func GetFileLogger(name string) (*logrus.Logger, error) {
	return g.Log.NewFileLogger(name)
}

func ErrorIfProd() error {
	if g.Env == sys.EnvironmentProd {
		return errors.Errorf("app environment is set to '%s'", sys.EnvironmentProd.String())
	}
	return nil
}

func DebugEnabled() bool {
	return g.HTTP.DebugEnabled()
}

// Validates an incoming request and binds the request body to the provided
// struct if the validation passes.
//
// Returns a 400 error with validation errors if binding fails.
//
// Sets the bound request as an interface{} in the Gin registry if binding
// succeeds.  You can retrieve in your handlers it like this:
//
// r, ok := goat.GetRequest(c).(*yourRequestStruct)
//
// This middleware allows you to annotate your request struct fields with
// `binding:"required"` to make required fields.
//
// @TODO it seems that if a request struct has a field that is named the same as one of it's child struct's fields that the validation messages don't prefix the field name with child struct's name
func BindMiddleware(r interface{}) gin.HandlerFunc {
	return g.HTTP.BindMiddleware(r)
}

// Deprecated in favor of BindMiddleware; preserving for backwards-compatibility.
func BindRequestMiddleware(req interface{}) gin.HandlerFunc {
	return BindMiddleware(req)
}

// Returns the bound request struct from the provided Gin context or nil if a goat request has not been bound.
// After binding a request using BindMiddleware, call this function to retrieve it in your handler:
// 	req, ok := goat.GetRequest(c).(*MyRequestType)
//	if !ok {
//		h.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
//		return
//	}
func GetRequest(c *gin.Context) interface{} {
	return g.HTTP.GetRequest(c)
}

func FilterMiddleware() gin.HandlerFunc {
	return g.HTTP.FilterMiddleware()
}

func GetFilter(c *gin.Context) *query.Query {
	return g.HTTP.GetFilter(c)
}

func GetHandlerContext(c *gin.Context) context.Context {
	return g.HTTP.GetHandlerContext(c)
}

// Returns a random string that can be used as a Basic Auth token.
func GenerateToken() string {
	u := uuid.New().String()
	return strings.Replace(u, "-", "", -1)
}
