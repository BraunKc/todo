package app

import (
	"fmt"

	"github.com/braunkc/todo/api-service/config"
	api "github.com/braunkc/todo/api-service/internal/http"
	"github.com/braunkc/todo/api-service/internal/utils"
)

func Run() {
	config.Logger = utils.InitLogger()
	config.InitYaml()
	fmt.Println(config.Config.DBService.Endpoint)
	api.InitRoutes()
	defer config.Logger.Sync()
}
