package http

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gopkg.in/gin-contrib/cors.v1"

	"github.com/68696c6c/goat/sys/http/links"
	"github.com/68696c6c/goat/sys/log"
)

type Router interface {
	gin.IRoutes
	http.Handler

	Group(path string, handlers ...gin.HandlerFunc) *Group
	GetUrl(key ...string) *url.URL
	SetUrl(key string, value *url.URL)
	Run() error
}

func NewRouter(config Config, log log.Service) Router {
	host :=  config.BaseUrl.Hostname()
	port :=  config.BaseUrl.Port()
	return &router{
		Engine:  newEngine(config, log),
		host:   host,
		port:    port,
		address: fmt.Sprintf("%s:%s", host, port),
		links:   links.NewService(config.BaseUrl),
	}
}

type router struct {
	*gin.Engine
	host    string
	port    string
	address string
	links   links.Service
}

func (r *router) Run() error {
	err := validPort(r.port)
	if err != nil {
		return errors.Wrap(err, "failed to start router")
	}
	return r.Engine.Run(r.address)
}

func (r *router) GetUrl(key ...string) *url.URL {
	return r.links.GetUrl(key...)
}

func (r *router) SetUrl(key string, value *url.URL) {
	r.links.SetUrl(key, value)
}

func (r *router) Group(relativePath string, handlers ...gin.HandlerFunc) *Group {
	result := r.Engine.Group(relativePath, handlers...)
	return &Group{
		RouterGroup: result,
		links:       r.links,
	}
}

func newEngine(c Config, l log.Service) *gin.Engine {
	engine := gin.New()

	// Set Gin debug mode
	mode := gin.ReleaseMode
	if c.Debug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	// Setup logging.
	engine.Use(l.GinLogger())
	engine.Use(l.GinRecovery())

	// Configure CORS.
	engine.Use(cors.New(c.GetCors()))

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

// validPort returns an error if the port is already in use.
func validPort(port string) error {
	_, err := strconv.Atoi(port)
	if err != nil {
		return errors.Errorf("port %v is not an integer", port)
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return errors.Errorf("port %v is already in use", port)
	}
	_ = ln.Close()
	return nil
}
