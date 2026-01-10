package logger

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func Logger() *zap.SugaredLogger {
	return logger
}

func Init() *zap.SugaredLogger {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger = zapLogger.Sugar()
	return logger
}
