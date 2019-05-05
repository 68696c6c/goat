package http

import (
	"io"
	"os"

	"github.com/68696c6c/goat/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/gin-contrib/cors.v1"
)

const (
	htttpMode    = gin.DebugMode
	httpPort     = "80"
	httpAuthType = "basic"
)

// Goat writes all request logging to standard out and always enables CORS.
// Since Goat is intended to be used to build RESTful APIs, it is assumed that
// the router will run on port 80, that DELETE requests will need to be
// supported, and that Basic Auth will be used.
// By default, all origins are allows for CORS.
// @TODO add support for more auth types.
type Service interface {
	NewRouter(setRoutes func(Router)) Router
}

type Config struct {
	Mode                  string
	Port                  string
	AuthType              string
	DisableCORSAllOrigins bool
	DisableDeleteMethod   bool
}

type ServiceGin struct {
	mode              string
	port              string
	authType          string
	disableAllOrigins bool
	disableDelete     bool
}

func NewServiceGin(c Config) ServiceGin {
	return ServiceGin{
		mode:              utils.ArgStringD(c.Mode, htttpMode),
		port:              utils.ArgStringD(c.Port, httpPort),
		authType:          utils.ArgStringD(c.AuthType, httpAuthType),
		disableAllOrigins: c.DisableCORSAllOrigins,
		disableDelete:     c.DisableDeleteMethod,
	}
}

func (s ServiceGin) NewRouter(setRoutes func(Router)) Router {
	r := NewRouterGin(s.port)

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

	// Setup routes.
	setRoutes(&r)

	return r
}
