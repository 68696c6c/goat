package sys

import (
	db "github.com/68696c6c/goat/sys/database"
	log "github.com/68696c6c/goat/sys/logging"
	"github.com/68696c6c/goat/sys/router"
)

type Goat struct {
	config    Config
	HttpDebug bool
	DB        db.Service
	Log       log.Service
	Router    router.Service
}

func Init() Goat {
	config := mustGetConfig()
	return Goat{
		config:    config,
		HttpDebug: config.HttpDebug,
		DB:        db.NewService(config.DB),
		Log:       log.NewService(config.Log),
		Router:    router.NewService(config.Router),
	}
}
