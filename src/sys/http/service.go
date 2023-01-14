package http

import (
	"net/url"

	"gopkg.in/gin-contrib/cors.v1"
)

type Config struct {
	CORS    cors.Config
	BaseUrl *url.URL
	Debug   bool
	Host    string
	Port    int
}

type Service interface {
	InitRouter() Router
	GetUrl(key ...string) *url.URL
	DebugEnabled() bool
}

type service struct {
	config Config
	router Router
}

func NewService(c Config) Service {
	return &service{
		config: c,
		router: nil,
	}
}

func (s *service) initRouter() Router {
	if s.router == nil {
		s.router = NewRouter(s.config)
	}
	return s.router
}

func (s *service) InitRouter() Router {
	return s.initRouter()
}

func (s *service) GetUrl(key ...string) *url.URL {
	return s.initRouter().GetUrl(key...)
}

func (s *service) DebugEnabled() bool {
	return s.config.Debug
}
