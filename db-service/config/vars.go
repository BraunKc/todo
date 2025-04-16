package config

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger *zap.Logger
	DB     *gorm.DB
)
