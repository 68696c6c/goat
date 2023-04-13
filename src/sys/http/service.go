package http

import (
	"net/url"

	"gopkg.in/gin-contrib/cors.v1"

	"github.com/68696c6c/goat/sys/log"
)

type Config struct {
	CORS    cors.Config
	BaseUrl *url.URL
	Debug   bool
	Host    string
	Port    int
}

func (c Config) GetCors() cors.Config {
	result := c.CORS
	if len(result.AllowOrigins) == 1 && result.AllowOrigins[0] == "*" {
		result.AllowAllOrigins = true
		result.AllowOrigins = []string{}
	}
	if len(result.AllowMethods) == 1 && result.AllowMethods[0] == "*" {
		result.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}
	}
	return result
}

type Service interface {
	SetBaseUrl(baseUrl string) error
	GetBaseUrl() *url.URL
	InitRouter() Router
	GetUrl(key ...string) *url.URL
	SetUrl(key string, value *url.URL)
	DebugEnabled() bool
}

type service struct {
	config Config
	log    log.Service
	router Router
}

func NewService(c Config, l log.Service) Service {
	return &service{
		config: c,
		log:    l,
		router: nil,
	}
}

func (s *service) initRouter() Router {
	if s.router == nil {
		s.router = NewRouter(s.config, s.log)
	}
	return s.router
}

func (s *service) InitRouter() Router {
	return s.initRouter()
}

func (s *service) SetBaseUrl(baseUrl string) error {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return err
	}
	s.config.BaseUrl = u
	return nil
}

func (s *service) GetBaseUrl() *url.URL {
	return s.config.BaseUrl
}

func (s *service) GetUrl(key ...string) *url.URL {
	return s.initRouter().GetUrl(key...)
}

func (s *service) SetUrl(key string, value *url.URL) {
	s.initRouter().SetUrl(key, value)
}

func (s *service) DebugEnabled() bool {
	return s.config.Debug
}
