package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Router interface {
	Run() error
	SetRegistry(d map[string]interface{})
	InitRegistry() gin.HandlerFunc
}

type RouterGin struct {
	Engine *gin.Engine
	port   string
	config map[string]interface{}
}

func NewRouterGin(port string) RouterGin {
	return RouterGin{
		port: port,
	}
}

// Run the Gin engine.
func (r RouterGin) Run() error {
	if err := validPort(r.port); err != nil {
		return errors.Wrap(err, "failed to start router")
	}
	return r.Engine.Run(r.port)
}

// Set a map of key-value pairs that will be added to the Gin registry when the
// router initializes.
func (r RouterGin) SetRegistry(d map[string]interface{}) {
	r.config = d
}

// Copy router config values to the Gin registry where they can be accessed by
// the app.
func (r RouterGin) InitRegistry() gin.HandlerFunc {
	return func(c *gin.Context) {
		for key, value := range r.config {
			c.Set(key, value)
		}
		c.Next()
	}
}
