package utils

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	genericLogger *zap.SugaredLogger
	initOnce      = sync.Once{}
)

func InitializeLogger() {
	var config zap.Config

	if !IsDebugMode() {
		config = zap.NewProductionConfig()
		config.OutputPaths = []string{"logs/moncaveau.log"}
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.TimeKey = "TIME"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.LevelKey = "LEVEL"
		config.EncoderConfig.MessageKey = "MESSAGE"
		config.EncoderConfig.CallerKey = "CALLER"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	genericLogger = logger.Sugar()
	defer logger.Sync()
}

func CreateLogger(packageName string) *zap.SugaredLogger {
	initOnce.Do(InitializeLogger)
	return genericLogger.Named(packageName)
}
