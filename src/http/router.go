package http

import (
	"io"
	"os"

	"github.com/68696c6c/goat/src/logging"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gopkg.in/gin-contrib/cors.v1"
)

type Router struct {
	Mode       string
	Engine     *gin.Engine
	RoutesFunc func(*Router)
	DB         *gorm.DB
	config     map[string]interface{}
}

// Router constructor.
func NewRouter(mode string, getRoutes func(*Router)) *Router {
	r := &Router{
		Mode:       mode,
		RoutesFunc: getRoutes,
	}
	return r.initRouter()
}

// Run the Gin engine.
func (r *Router) Run(addr ...string) error {
	return r.Engine.Run(addr...)
}

// Set a map of key-value pairs that will be added to the Gin registry when the
// router initializes.
func (r *Router) SetConfig(d map[string]interface{}) {
	r.config = d
}

// Copy router config values to the Gin registry where they can be accessed by
// the app.
func (r *Router) initRegistry() gin.HandlerFunc {
	return func(c *gin.Context) {
		for key, value := range r.config {
			c.Set(key, value)
		}
		c.Next()
	}
}

// Initialize Gin and call the routes callback.
func (r *Router) initRouter() *Router {
	gin.SetMode(r.Mode)
	r.Engine = gin.New()
	r.Engine.Use(gin.Recovery())

	// Setup logging.
	logger := logging.NewCustomLogger("http")
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logger.Out)
	r.Engine.Use(gin.Logger())

	// Configure CORS.
	config := cors.DefaultConfig()
	config.AllowMethods = append(config.AllowMethods, "DELETE")
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowAllOrigins = true
	r.Engine.Use(cors.New(config), r.initRegistry())

	// Setup routes.
	r.RoutesFunc(r)

	return r
}
