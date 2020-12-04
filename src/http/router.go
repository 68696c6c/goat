package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type RouterGin struct {
	Engine  *gin.Engine
	host    string
	port    string
	address string
	config  map[string]interface{}
}

func NewRouterGin(host, port string) *RouterGin {
	return &RouterGin{
		host:    host,
		port:    port,
		address: fmt.Sprintf("%s:%s", host, port),
	}
}

func (r *RouterGin) GetEngine() *gin.Engine {
	return r.Engine
}

// Run the Gin engine.
func (r *RouterGin) Run() error {
	if err := validPort(r.port); err != nil {
		return errors.Wrap(err, "failed to start router")
	}
	return r.Engine.Run(r.address)
}

// Set a map of key-value pairs that will be added to the Gin registry when the
// router initializes.
func (r *RouterGin) SetRegistry(d map[string]interface{}) {
	r.config = d
}

// Copy router config values to the Gin registry where they can be accessed by
// the app.
func (r *RouterGin) InitRegistry() gin.HandlerFunc {
	return func(c *gin.Context) {
		for key, value := range r.config {
			c.Set(key, value)
		}
		c.Next()
	}
}
