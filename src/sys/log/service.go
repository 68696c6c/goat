package log

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	gormlog "gorm.io/gorm/logger"
	gormzap "moul.io/zapgorm2"
)

type Config struct {
	Level      zap.AtomicLevel
	Stacktrace bool
}

type Service interface {
	GetLogger() *zap.SugaredLogger
	GetStrictLogger() *zap.Logger
	GinLogger() gin.HandlerFunc
	GinRecovery() gin.HandlerFunc
	GormLogger() gormlog.Interface
}

func NewService(c Config) (Service, error) {
	logger, err := initLogger(c)
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

func (s service) GetLogger() *zap.SugaredLogger {
	return s.logger.Sugar()
}

func (s service) GetStrictLogger() *zap.Logger {
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

func initLogger(c Config) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = c.Level
	config.DisableStacktrace = !c.Stacktrace
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
