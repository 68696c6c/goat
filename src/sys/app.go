package sys

import (
	db "github.com/68696c6c/goat/sys/database"
	"github.com/68696c6c/goat/sys/http/router"
	log "github.com/68696c6c/goat/sys/logging"
)

type Goat struct {
	config    Config
	HttpDebug bool
	DB        db.Service
	// HTTP      http.Service
	Log    log.Service
	Router router.Service
}

func Init() Goat {
	config := mustGetConfig()
	return Goat{
		config:    config,
		HttpDebug: config.HttpDebug,
		DB:        db.NewService(config.DB),
		// HTTP:      http.NewServiceGin(config.HTTP),
		Log:    log.NewServiceLogrus(config.Log),
		Router: router.NewService(config.Router),
	}
}
