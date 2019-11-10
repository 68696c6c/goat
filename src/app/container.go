package app

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type App interface {
	GetDB() *gorm.DB
	GetLogger() *logrus.Logger
}

type Initializer func() App
