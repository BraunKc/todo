package redis

import (
	"os"
	"strconv"

	"github.com/BraunKc/todo/db-service/config"
	"github.com/BraunKc/todo/db-service/internal/repository/models"
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

	var allTasks []models.Task
	config.DB.Model(&models.Task{}).Find(&allTasks)

	for _, task := range allTasks {
		err := client.HMSet(strconv.FormatUint(uint64(task.ID), 10), map[string]interface{}{
			"name":        task.Name,
			"description": task.Description,
			"complete":    task.Complete,
		}).Err()
		if err != nil {
			config.Logger.Error("add task in redis error", zap.Error(err))
		}
	}

	config.Logger.Info("redis inited", zap.Any("redis", client))

	return client
}
