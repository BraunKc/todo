package app

import (
	"fmt"

	"github.com/BraunKc/todo/api-service/config"
	api "github.com/BraunKc/todo/api-service/internal/http"
	"github.com/BraunKc/todo/api-service/internal/utils"
)

func Run() {
	config.Logger = utils.InitLogger()
	config.InitYaml()
	fmt.Println(config.Config.DBService.Endpoint)
	api.InitRoutes()
	defer config.Logger.Sync()
}
