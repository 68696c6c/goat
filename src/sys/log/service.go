package log

import (
	"time"

	gormzap "github.com/68696c6c/zapgorm2"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	gormlog "gorm.io/gorm/logger"
)

type Config struct {
	Level      zap.AtomicLevel
	Stacktrace bool
}

type Service interface {
	Logger() *zap.SugaredLogger
	StrictLogger() *zap.Logger
	GinLogger() gin.HandlerFunc
	GinRecovery() gin.HandlerFunc
	GormLogger() gormlog.Interface
}

func NewService(c Config) (Service, error) {
	config := zap.NewProductionConfig()
	config.Level = c.Level
	config.DisableStacktrace = !c.Stacktrace
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return service{
		config: c,
		logger: logger,
	}, nil
}

type service struct {
	config Config
	logger *zap.Logger
}

func (s service) Logger() *zap.SugaredLogger {
	return s.logger.Sugar()
}

func (s service) StrictLogger() *zap.Logger {
	return s.logger
}

func (s service) GinLogger() gin.HandlerFunc {
	return ginzap.Ginzap(s.logger, time.RFC3339, true)
}

func (s service) GinRecovery() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(s.logger, s.config.Stacktrace)
}

func (s service) GormLogger() gormlog.Interface {
	logger := gormzap.New(s.logger)
	logger.SetAsDefault()
	return logger
}
