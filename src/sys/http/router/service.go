package router

import (
	"net/url"
)

type Config struct {
	BaseUrl               *url.URL
	Debug                 bool
	Host                  string
	Port                  string
	AuthType              string
	DisableCORSAllOrigins bool
	DisableDeleteMethod   bool
}

type Service interface {
	InitRouter() Router
	GetUrl(key ...string) *url.URL
	// GetValidator() (*validator.Validate, error)
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

// func (s *service) GetValidator() (*validator.Validate, error) {
// 	return s.initRouter().GetValidator()
// }
