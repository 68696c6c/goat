package http

import (
	"github.com/gin-gonic/gin"

	"github.com/68696c6c/goat/sys/http/links"
)

type Group struct {
	*gin.RouterGroup
	links links.Service
}

func (g *Group) Group(key, relativePath string, handlers ...gin.HandlerFunc) *Group {
	result := g.RouterGroup.Group(relativePath, handlers...)
	g.links.AddBaseUrlPath(key, result.BasePath())
	return &Group{
		RouterGroup: result,
		links:       g.links,
	}
}
