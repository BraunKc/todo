package utils

import "go.uber.org/zap"

func InitLogger() *zap.Logger {
	logger, _ := zap.NewProduction()

	logger.Debug("logger inited:", zap.Any("logger", logger))

	return logger
}
