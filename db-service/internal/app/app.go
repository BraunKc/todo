package app

import (
	"github.com/BraunKc/todo/db-service/config"
	api "github.com/BraunKc/todo/db-service/internal/http"
	database "github.com/BraunKc/todo/db-service/internal/repository/gorm"
	"github.com/BraunKc/todo/db-service/internal/utils"
)

func Run() {
	config.Logger = utils.InitLogger()
	config.DB = database.InitDatabase()
	api.InitRoutes()
	defer config.Logger.Sync()
}
