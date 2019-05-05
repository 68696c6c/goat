package sys

import (
	"github.com/68696c6c/goat/src/cmd"
	db "github.com/68696c6c/goat/src/database"
	"github.com/68696c6c/goat/src/http"
	log "github.com/68696c6c/goat/src/logging"
)

type Goat struct {
	config Config
	Env    Environment
	CMD    cmd.Service
	DB     db.Service
	HTTP   http.Service
	Log    log.Service
}

func Init() Goat {
	mustReadConfig()
	return Goat{
		config: config,
		Env:    config.Env,
		CMD:    cmd.NewServiceCobra(config.CMD),
		DB:     db.NewServiceGORM(config.DB),
		HTTP:   http.NewServiceGin(config.HTTP),
		Log:    log.NewServiceLogrus(config.Log),
	}
}
