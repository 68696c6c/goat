package log

import "go.uber.org/zap"

type Config struct {
	Level      zap.AtomicLevel
	Stacktrace bool
}

func InitLogger(c Config) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.Level = c.Level
	config.DisableStacktrace = !c.Stacktrace
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
