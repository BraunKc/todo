package redis

import (
	"os"

	"github.com/BraunKc/todo/db-service/config"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func InitRedis() *redis.Client {
	if err := godotenv.Load(); err != nil {
		config.Logger.Fatal("no .env file found", zap.Error(err))
	}
	addr := os.Getenv("REDIS_ADDR")

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := client.Ping().Result()
	if err != nil {
		config.Logger.Fatal("redis init error", zap.Error(err))
	}

	config.Logger.Info("redis inited", zap.Any("redis", client))

	return client
}
