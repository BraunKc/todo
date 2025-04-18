package config

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger      *zap.Logger
	DB          *gorm.DB
	RedisClient *redis.Client
)
