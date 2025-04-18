package app

import (
	"github.com/BraunKc/todo/db-service/config"
	"github.com/BraunKc/todo/db-service/internal/cache/redis"
	api "github.com/BraunKc/todo/db-service/internal/http"
	database "github.com/BraunKc/todo/db-service/internal/repository/gorm"
	"github.com/BraunKc/todo/db-service/internal/utils"
)

func Run() {
	config.Logger = utils.InitLogger()
	config.DB = database.InitDatabase()
	config.RedisClient = redis.InitRedis()
	api.InitRoutes()
	defer config.Logger.Sync()
}
