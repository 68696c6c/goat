package http

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/resource"
	"github.com/68696c6c/goat/sys/utils"
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
// Requests are bound to an application-level (i.e. non-Goat) struct provided by
// the caller.
//
// Request struct fields can be annotated with tags to mark them as required and
// to provide custom labels for use in validation error messages.
//
// For example:
//
//	type userRequest struct {
//		Name string `json:"name" binding:"required" label:"Username"`
//	}
//
// Note that the `binding:"required"` annotation is not recursive.  For example,
// if an Order request struct has an Items field that is a slice of Item structs
// and the Item struct also has fields that are required, the binding service
// will not throw an error if one of the Items is missing a field.  To change
// this behavior, use 'dive,required' annotation to tell the binding service to
// check the field's own fields for required annotations as well.
//
// For example:
//
//	type orderRequest struct {
//		Items []*Item `binding:"required,dive,required"`
//	}
//
// If a required field is not provided, or a request body is not sent, the
// binding service will automatically send a bad request response with a list of
// errors.  The service will attempt to use the JSON names of the fields in the
// error messages rather than the struct field names so that the user sees the
// field name as they sent it, e.g. "name" instead of "Name".  If a custom label
// annotation exists on the field, that will be used instead.
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
//	BINDING_EXCLUDED_STRUCTS="goat.ID,time.Time"
//
// Since the binding service will return an error response if an error occurs,
// binding is done in a middleware.  The middleware will set the  bound request
// in the Gin context where it can be accessed by subsequent handler functions.
//
// @TODO add support for more auth types.
// TODO: add base url
type Service interface {
	DebugEnabled() bool
	// NewRouter() *RouterGin
	// NewRouter() router.Router
	GetHandlerContext(cx *gin.Context) context.Context
	BindMiddleware(r any) gin.HandlerFunc
	GetRequest(cx *gin.Context) any
	FilterMiddleware() gin.HandlerFunc
	GetFilter(cx *gin.Context) *query.Query
	RespondValid(cx *gin.Context)
	RespondOk(cx *gin.Context, data any)
	RespondUpdated(cx *gin.Context, data any)
	RespondAccepted(cx *gin.Context, data any)
	RespondCreated(cx *gin.Context, data any)
	RespondNotFound(cx *gin.Context, err error)
	RespondBadRequest(cx *gin.Context, err error)
	RespondValidationError(cx *gin.Context, err error)
	// RespondValidationErrors(cx *gin.Context, errs map[string]error)
	RespondUnauthorized(cx *gin.Context, err error)
	RespondForbidden(cx *gin.Context, err error)
	RespondServerError(cx *gin.Context, err error)
	// router.Service
}

// type RouterConfig router.Config

type Config struct {
	// BaseUrl               *url.URL
	Debug bool
	// Host                  string
	// Port                  string
	// AuthType              string
	// DisableCORSAllOrigins bool
	// DisableDeleteMethod   bool
	// RouterConfig
	ExcludedStructs string
}

type service struct {
	config Config
	// routerConfig      router.Config
	// baseUrl           *url.URL
	mode string
	// host              string
	// port              string
	// authType          string
	// disableAllOrigins bool
	// disableDelete     bool
	// structMetaService meta.Service
}

func NewServiceGin(c Config) Service {
	mode := gin.ReleaseMode
	if c.Debug {
		mode = gin.DebugMode
	}
	return service{
		config: c,
		// routerConfig:      router.Config(c.RouterConfig),
		// baseUrl:           c.BaseUrl,
		mode: mode,
		// host:              c.Host,
		// port:              c.Port,
		// authType:          c.AuthType,
		// disableAllOrigins: c.DisableCORSAllOrigins,
		// disableDelete:     c.DisableDeleteMethod,
		// structMetaService: meta.NewService(c.ExcludedStructs),
	}
}

func (s service) DebugEnabled() bool {
	return s.mode == gin.DebugMode
}

// func (s service) InitRouter() router.Router {
// 	println("ROUTER HOST: " + s.host)
// 	println("ROUTER PORT: " + s.port)
// 	println("ROUTER BASE URL: " + s.baseUrl.String())
// 	return router.NewRouter(s.routerConfig)
// 	// return router.NewRouter(router.Config{
// 	// 	BaseUrl:               s.config.BaseUrl,
// 	// 	Debug:                 s.config.Debug,
// 	// 	Host:                  s.config.Host,
// 	// 	Port:                  s.config.Port,
// 	// 	AuthType:              s.config.AuthType,
// 	// 	DisableCORSAllOrigins: s.config.DisableCORSAllOrigins,
// 	// 	DisableDeleteMethod:   s.config.DisableDeleteMethod,
// 	// 	ExcludedStructs:       s.config.ExcludedStructs,
// 	// })
// }

// func (s service) NewRouter() *RouterGin {
// 	println("ROUTER HOST: " + s.host)
// 	println("ROUTER PORT: " + s.port)
// 	println("ROUTER BASE URL: " + s.baseUrl.String())
// 	r := NewRouterGin(s.host, s.port, s.baseUrl)
//
// 	gin.SetMode(s.mode)
// 	r.Engine = gin.New()
// 	r.Engine.Use(gin.Recovery())
//
// 	// Setup logging.
// 	gin.DefaultWriter = io.MultiWriter(os.Stdout)
// 	r.Engine.Use(gin.Logger())
//
// 	// Configure CORS.
// 	config := cors.DefaultConfig()
// 	if s.authType == "basic" {
// 		config.AllowHeaders = append(config.AllowHeaders, "Authorization")
// 	}
// 	if !s.disableDelete {
// 		config.AllowMethods = append(config.AllowMethods, "DELETE")
// 	}
// 	if !s.disableAllOrigins {
// 		config.AllowAllOrigins = true
// 	}
// 	r.Engine.Use(cors.New(config), r.InitRegistry())
//
// 	return r
// }

func (s service) GetHandlerContext(c *gin.Context) context.Context {
	return c.Request.Context()
}

// BindMiddleware returns a gim.HandlerFunc that attempts to bind a JSON request body from the current Gin Context to
// the provided struct. If any of the struct's required fields are missing from the request body, a 400 response is sent.
func (s service) BindMiddleware(r any) gin.HandlerFunc {
	return func(cx *gin.Context) {
		value := reflect.ValueOf(r)
		if value.Kind() == reflect.Ptr {
			panic("bind struct must not be a pointer")
		}
		typ := value.Type()
		obj := reflect.New(typ).Interface()
		if err := cx.ShouldBindWith(obj, binding.JSON); err != nil {
			s.respondRequestBindingError(cx, err, typ)
			// s.RespondBadRequest(cx, err)
			return
		}
		cx.Set(contextKeyRequest, obj)
		return
	}
}

func (s service) GetRequest(cx *gin.Context) any {
	return s.getByKey(cx, contextKeyRequest)
}

func (s service) FilterMiddleware() gin.HandlerFunc {
	return func(cx *gin.Context) {
		q := query.NewQueryBuilder(cx)
		cx.Set(contextKeyQuery, q)
		return
	}
}

func (s service) GetFilter(c *gin.Context) *query.Query {
	r := s.getByKey(c, contextKeyQuery)
	return r.(*query.Query)
}

func (s service) respondRequestBindingError(c *gin.Context, err error, t reflect.Type) {
	// TODO: this provides error messages than the default gin binding/validation errors but it's a lot of complexity...
	// If no request body was sent at all, Gin will try to return 'EOF'
	// Show the user something more helpful instead.
	if err.Error() == "EOF" {
		s.RespondBadRequest(c, errors.New("a request body is required"))
		return
	}

	// Cast the error to a ValidationErrors struct so that we can access detailed
	// information about the error.  If the error cannot be cast, a generic
	// Bad Request response will be sent instead of the detailed message.
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		s.RespondBadRequest(c, errors.New("invalid request"))
		return
	}

	var errs []string
	for _, e := range ve {
		// errs = append(errs, e.Error())
		// errs = append(errs, fmt.Sprintf("%s failed on validation '%s'", e.Field(), e.Tag()))
		errs = append(errs, utils.MakeValidationError(e.Field(), e.Tag()).Error())
	}
	s.RespondValidationError(c, errors.New(strings.Join(errs, "; ")))

	// // Create an error message for each missing field.
	// msgs := make(map[string]error)
	// for _, e := range ve {
	// 	println(e.Error())
	// 	// m, err := s.structMetaService.GetStructFieldMeta(t, e.Field())
	// 	// if err != nil {
	// 	// 	// If we couldn't find a JSON tag annotation for the field, fallback to
	// 	// 	// the struct field name.
	// 	// 	m = &meta.FieldMeta{
	// 	// 		Path:  e.StructNamespace(),
	// 	// 		Label: e.Field(),
	// 	// 	}
	// 	// }
	// 	// msgs[m.Path] = errors.New(m.Label + " is required")
	// 	field := e.Field()
	// 	msgs[field] = errors.Errorf("%s failed on validation '%s'", field, e.Tag())
	// }
	// s.RespondValidationErrors(c, msgs)
	return
}

func (s service) getByKey(cx *gin.Context, key string) any {
	r, exists := cx.Get(key)
	if !exists {
		return nil
	}
	return r
}

func (s service) debugError(err error) error {
	// Only show errors to the user if we are in debug mode.
	if s.DebugEnabled() {
		return err
	}
	return nil
}

// func (s service) debugValidationErrors(errs map[string]error) error {
// 	var result []string
// 	for k, v := range errs {
// 		err := s.debugError(errors.Errorf("%s: %s", k, v.Error()))
// 		if err != nil {
// 			result = append(result, err.Error())
// 		}
// 	}
// 	return errors.New(strings.Join(result, ";"))
// }

func (s service) RespondValid(cx *gin.Context) {
	cx.AbortWithStatusJSON(http.StatusOK, true)
}

func (s service) RespondOk(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusOK, data)
}

func (s service) RespondUpdated(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusIMUsed, data)
}

func (s service) RespondAccepted(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusAccepted, data)
}

func (s service) RespondCreated(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusCreated, data)
}

func (s service) RespondNotFound(cx *gin.Context, err error) {
	status := http.StatusNotFound
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "not found", s.debugError(err)))
}

func (s service) RespondBadRequest(cx *gin.Context, err error) {
	status := http.StatusBadRequest
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "invalid request", s.debugError(err)))
}

// func (s service) RespondValidationErrors(cx *gin.Context, errs map[string]error) {
// 	status := http.StatusBadRequest
// 	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "invalid request", s.debugValidationErrors(errs)))
// }

func (s service) RespondValidationError(cx *gin.Context, err error) {
	status := http.StatusBadRequest
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "invalid request", err))
}

func (s service) RespondUnauthorized(cx *gin.Context, err error) {
	status := http.StatusUnauthorized
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "unauthorized", s.debugError(err)))
}

func (s service) RespondForbidden(cx *gin.Context, err error) {
	status := http.StatusForbidden
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "authentication error", s.debugError(err)))
}

func (s service) RespondServerError(cx *gin.Context, err error) {
	status := http.StatusInternalServerError
	cx.AbortWithStatusJSON(status, resource.NewApiProblem(status, "internal server error", s.debugError(err)))
}
