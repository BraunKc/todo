package utils

import (
	"log"

	"go.uber.org/zap"
)

func InitLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("logger inited")

	return logger
}
