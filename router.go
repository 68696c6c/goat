package goat

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"gopkg.in/gin-contrib/cors.v1"
	"github.com/jinzhu/gorm"
)

type Router struct {
	Mode       string
	Engine     *gin.Engine
	RoutesFunc func(*Router)
	DB         *gorm.DB
	config     map[string]interface{}
}

// Router constructor.  Will panic if goat has not been initialized.
func NewRouter(mode string, getRoutes func(*Router)) *Router {
	mustBeInitialized()
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

// Copy router config values to the Gun registry where they can be accessed by
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
	logger := NewCustomLogger("http")
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logger.Out)
	r.Engine.Use(gin.Logger())

	// Configure CORS.
	config := cors.DefaultConfig()
	config.AllowMethods = append(config.AllowMethods, "DELETE")
	config.AllowAllOrigins = true
	r.Engine.Use(cors.New(config), r.initRegistry())

	// Setup routes.
	r.RoutesFunc(r)

	return r
}
