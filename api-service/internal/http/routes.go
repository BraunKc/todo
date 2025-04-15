package api

import (
	"github.com/BraunKc/todo/api-service/config"
	"github.com/gin-gonic/gin"
)

func InitRoutes() {
	router := gin.Default()

	router.POST("/add", add)
	router.GET("/complete/:id", complete)
	router.GET("/tasks", tasks)
	router.GET("/task/:id", task)
	router.GET("/delete/:id", delete)
	router.GET("/clean", clean)

	config.Logger.Info("routes inited")

	router.Run(":8080")
}
