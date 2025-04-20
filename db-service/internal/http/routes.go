package api

import (
	"github.com/BraunKc/todo/db-service/config"
	"github.com/gin-gonic/gin"
)

func InitRoutes() {
	router := gin.Default()

	router.POST("/create-task", createTask)
	router.GET("/set-status-complete/:id", setStatusComplete)
	router.GET("/get-all-tasks", getAllTasks)
	router.GET("/get-task/:id", getTask)
	router.DELETE("/delete-task/:id", deleteTask)
	router.DELETE("/clean-completed-tasks", cleanCompleteTasks)

	config.Logger.Info("routes inited")

	router.Run(":8282")
}
