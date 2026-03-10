package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/dojo-harvester/backend/internal/config"
)

var log *zap.Logger

func Init(cfg *config.LoggerConfig) error {
	var zapCfg zap.Config

	if cfg.Development {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	level, err := zapcore.ParseLevel(cfg.Level)
	if err == nil {
		zapCfg.Level.SetLevel(level)
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return err
	}

	log = logger
	zap.ReplaceGlobals(logger)

	return nil
}

func Get() *zap.Logger {
	if log == nil {
		l, _ := zap.NewDevelopment()
		return l
	}
	return log
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
