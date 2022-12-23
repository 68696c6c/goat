package http

import (
	"context"
	"io"
	"net/http"
	"os"
	"reflect"

	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/src/types"
	"github.com/68696c6c/goat/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/go-playground/validator.v8"
)

const (
	httpHost          = ""
	httpPort          = "80"
	httpAuthType      = "basic"
	contextKeyRequest = "goat_request"
	contextKeyQuery   = "goat_query"
)

// Goat writes all request logging to standard out and always enables CORS.
// Since Goat is intended to be used to build RESTful APIs, it is assumed that
// the router will run on port 80, that DELETE requests will need to be
// supported, and that Basic Auth will be used.
// By default, all origins are allows for CORS.
//
//
// Requests are bound to an application-level (i.e. non-Goat) struct provided by
// the caller.
//
//
// Request struct fields can be annotated with tags to mark them as required and
// to provide custom labels for use in validation error messages.
//
// For example:
// 	type userRequest struct {
// 		Name string `json:"name" binding:"required" label:"Username"`
// 	}
//
// Note that the `binding:"required"` annotation is not recursive.  For example,
// if an Order request struct has an Items field that is a slice of Item structs
// and the Item struct also has fields that are required, the binding service
// will not throw an error if one of the Items is missing a field.  To change
// this behavior, use 'dive,required' annotation to tell the binding service to
// check the field's own fields for required annotations as well.
//
// For example:
// 	type orderRequest struct {
// 		Items []*Item `binding:"required,dive,required"`
// 	}
//
//
// If a required field is not provided, or a request body is not sent, the
// binding service will automatically send a bad request response with a list of
// errors.  The service will attempt to use the JSON names of the fields in the
// error messages rather than the struct field names so that the user sees the
// field name as they sent it, e.g. "name" instead of "Name".  If a custom label
// annotation exists on the field, that will be used instead.
//
//
// When building a validation error message, the binding service will use
// the fully-qualified JSON names of missing fields.  For example, if a request
// struct has a struct field named "Item" that has a required field named "SKU",
// the error message will show "items.sku" instead of "SKU" to make it easier
// for the user to determine where the missing field should be.  If yous request
// structs have fields that are also structs that will never have validation
// errors on their own internal fields (e.g. goat.ID, time.Time, etc), you
// should exclude these types from the error message building to avoid
// unnecessary recursion.  For example:
//
// 	BINDING_EXCLUDED_STRUCTS="goat.ID,time.Time"
//
//
// Since the binding service will return an error response if an error occurs,
// binding is done in a middleware.  The middleware will set the  bound request
// in the Gin context where it can be accessed by subsequent handler functions.
//
// @TODO add support for more auth types.
type Service interface {
	DebugEnabled() bool
	NewRouter() *RouterGin
	GetHandlerContext(c *gin.Context) context.Context
	BindMiddleware(r interface{}) gin.HandlerFunc
	GetRequest(c *gin.Context) interface{}
	FilterMiddleware() gin.HandlerFunc
	GetFilter(c *gin.Context) *query.Query
}

type Config struct {
	Debug                 bool
	Host                  string
	Port                  string
	AuthType              string
	DisableCORSAllOrigins bool
	DisableDeleteMethod   bool
	ExcludedStructs       string
}

type ServiceGin struct {
	mode              string
	host              string
	port              string
	authType          string
	disableAllOrigins bool
	disableDelete     bool
	structMetaService *structMetaService
}

func NewServiceGin(c Config) ServiceGin {
	mode := gin.ReleaseMode
	if c.Debug {
		mode = gin.DebugMode
	}
	return ServiceGin{
		mode:              mode,
		host:              utils.ArgStringD(c.Host, httpHost),
		port:              utils.ArgStringD(c.Port, httpPort),
		authType:          utils.ArgStringD(c.AuthType, httpAuthType),
		disableAllOrigins: c.DisableCORSAllOrigins,
		disableDelete:     c.DisableDeleteMethod,
		structMetaService: newStructMetaService(c.ExcludedStructs),
	}
}

func (s ServiceGin) DebugEnabled() bool {
	return s.mode == gin.DebugMode
}

func (s ServiceGin) NewRouter() *RouterGin {
	r := NewRouterGin(s.host, s.port)

	gin.SetMode(s.mode)
	r.Engine = gin.New()
	r.Engine.Use(gin.Recovery())

	// Setup logging.
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	r.Engine.Use(gin.Logger())

	// Configure CORS.
	config := cors.DefaultConfig()
	if s.authType == "basic" {
		config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	}
	if !s.disableDelete {
		config.AllowMethods = append(config.AllowMethods, "DELETE")
	}
	if !s.disableAllOrigins {
		config.AllowAllOrigins = true
	}
	r.Engine.Use(cors.New(config), r.InitRegistry())

	return r
}

func (s ServiceGin) GetHandlerContext(c *gin.Context) context.Context {
	return c.Request.Context()
}

// Attempts to bind a JSON request body from the Gin Context to the provided
// struct. If any of the struct's required fields are missing from the request
// body, a 400 response is sent.
func (s ServiceGin) BindMiddleware(r interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := reflect.ValueOf(r)
		if value.Kind() == reflect.Ptr {
			panic("Bind struct can not be a pointer")
		}
		typ := value.Type()
		obj := reflect.New(typ).Interface()
		if err := c.ShouldBindWith(obj, binding.JSON); err != nil {
			s.respondRequestBindingError(c, err, typ)
			return
		}
		c.Set(contextKeyRequest, obj)
		return
	}
}

func (s ServiceGin) GetRequest(c *gin.Context) interface{} {
	return s.getByKey(c, contextKeyRequest)
}

func (s ServiceGin) FilterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := query.NewQueryBuilder(c)
		c.Set(contextKeyQuery, q)
		return
	}
}

func (s ServiceGin) GetFilter(c *gin.Context) *query.Query {
	r := s.getByKey(c, contextKeyQuery)
	return r.(*query.Query)
}

func (s ServiceGin) respondRequestBindingError(c *gin.Context, err error, t reflect.Type) {
	// If no request body was sent at all, Gin will try to return 'EOF'
	// Show the user something more helpful instead.
	if err.Error() == "EOF" {
		e := errors.New("a request body is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, types.Response{"Bad Request.", []error{e}, nil})
		return
	}

	// Cast the error to a ValidationErrors struct so that we can access detailed
	// information about the error.  If the error cannot be cast, a generic
	// Bad Request response will be sent instead of the detailed message.
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		e := errors.New("invalid request")
		c.AbortWithStatusJSON(http.StatusBadRequest, types.Response{"Bad Request.", []error{e}, nil})
		return
	}

	// Create an error message for each missing field.
	msgs := make(map[string]string)
	for _, e := range ve {
		meta, err := s.structMetaService.GetStructFieldMeta(t, e.Field)
		if err != nil {
			// If we couldn't find a JSON tag annotation for the field, fallback to
			// the struct field name.
			meta = &fieldMeta{
				Path:  e.Field,
				Label: e.Name,
			}
		}
		msgs[meta.Path] = meta.Label + " is required"
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, types.ValidationResponse{"Invalid Request.", msgs})
}

func (s ServiceGin) getByKey(c *gin.Context, key string) interface{} {
	r, exists := c.Get(key)
	if !exists {
		return nil
	}
	return r
}
