package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

type App interface {
	GetDB() *gorm.DB
	GetLogger() *logrus.Logger
}

type Initializer func() App
