package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type App interface {
	GetDB() *gorm.DB
	GetLogger() *logrus.Logger
}

type Initializer func() App

type Router interface {
	Run() error
	SetRegistry(d map[string]interface{})
	InitRegistry() gin.HandlerFunc
	GetEngine() *gin.Engine
}
