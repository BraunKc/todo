package app

import (
	"github.com/BraunKc/todo/api-service/config"
	api "github.com/BraunKc/todo/api-service/internal/http"
	"github.com/BraunKc/todo/api-service/internal/utils"
)

func Run() {
	config.Logger = utils.InitLogger()
	api.InitRoutes()
	defer config.Logger.Sync()
}
