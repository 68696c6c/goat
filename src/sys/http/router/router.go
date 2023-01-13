package router

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gopkg.in/gin-contrib/cors.v1"

	"github.com/68696c6c/goat/sys/http/router/links"
)

// type Routes interface {
// 	Use(...gin.HandlerFunc) Routes
//
// 	Handle(key, method, relativePath string, handlers ...gin.HandlerFunc) Routes
// 	Any(key, relativePath string, handlers ...gin.HandlerFunc) Routes
// 	GET(key, relativePath string, handlers ...gin.HandlerFunc) Routes
// 	POST(key, relativePath string, handlers ...gin.HandlerFunc) Routes
// 	DELETE(key, relativePath string, handlers ...gin.HandlerFunc) Routes
// 	PATCH(key, relativePath string, handlers ...gin.HandlerFunc) Routes
// 	PUT(key, relativePath string, handlers ...gin.HandlerFunc) Routes
// 	OPTIONS(key, relativePath string, handlers ...gin.HandlerFunc) Routes
// 	HEAD(key, relativePath string, handlers ...gin.HandlerFunc) Routes
//
// 	StaticFile(key, relativePath, filePath string) Routes
// 	StaticFileFS(key, relativePath, filePath string, fs http.FileSystem) Routes
// 	Static(key, relativePath, root string) Routes
// 	StaticFS(key, relativePath string, fs http.FileSystem) Routes
// }

type Router interface {
	gin.IRoutes
	http.Handler

	Group(path string, handlers ...gin.HandlerFunc) *Group
	GetUrl(key ...string) *url.URL
	GetValidator() (*validator.Validate, error)
	Run() error
}

func NewRouter(c Config) Router {
	mode := gin.ReleaseMode
	if c.Debug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	// Setup logging.
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	engine.Use(gin.Logger())

	return &router{
		Engine:  newEngine(c),
		host:    c.Host,
		port:    c.Port,
		address: fmt.Sprintf("%s:%s", c.Host, c.Port),
		links:   links.NewService(c.BaseUrl),
	}
}

func newEngine(c Config) *gin.Engine {
	mode := gin.ReleaseMode
	if c.Debug {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	// Setup logging.
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	engine.Use(gin.Logger())

	// Configure CORS.
	config := cors.DefaultConfig()
	if c.AuthType == "basic" {
		config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	}
	if !c.DisableCORSAllOrigins {
		config.AllowAllOrigins = true
	}
	if !c.DisableDeleteMethod {
		config.AllowMethods = append(config.AllowMethods, "DELETE")
	}
	engine.Use(cors.New(config))

	// Use json tag names in request binding validation errors.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	return engine
}

type router struct {
	*gin.Engine
	baseUrl string
	host    string
	port    string
	address string
	links   links.Service
}

func (r *router) Run() error {
	if err := validPort(r.port); err != nil {
		return errors.Wrap(err, "failed to start router")
	}
	return r.Engine.Run(r.address)
}

func (r *router) GetUrl(key ...string) *url.URL {
	return r.links.GetUrl(key...)
}

func (r *router) GetValidator() (*validator.Validate, error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		return v, nil
	}
	return nil, errors.New("failed to get validator")
}

func (r *router) Group(relativePath string, handlers ...gin.HandlerFunc) *Group {
	result := r.Engine.Group(relativePath, handlers...)
	return &Group{
		RouterGroup: result,
		links:       r.links,
	}
}

// Determines whether the provided value is a valid port that can be listened on.
func validPort(port string) error {

	// Must be numeric.
	if _, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("%s is not a valid port", port)
	}

	// Try and listen to see if the port is available.
	if ln, err := net.Listen("tcp", ":"+port); err == nil {
		_ = ln.Close()
		return nil
	}

	return fmt.Errorf("port %s is already in use", port)
}
