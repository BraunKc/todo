package api

import (
	"net/http"
	"strconv"

	"github.com/BraunKc/todo/db-service/config"
	"github.com/BraunKc/todo/db-service/internal/repository/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func createTask(c *gin.Context) {
	createTask := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{}
	if err := c.ShouldBindJSON(&createTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	if createTask.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "name must be not empty",
		})
		return
	}

	newTask := models.Task{
		Name:        createTask.Name,
		Description: createTask.Description,
		Complete:    false,
	}

	task := config.DB.First(&models.Task{}, "name =?", newTask.Name)
	if task.RowsAffected != 0 {
		c.JSON(http.StatusConflict, gin.H{
			"err": "task with same name already created",
		})
		return
	}

	result := config.DB.Create(&newTask)
	if result.Error != nil {
		config.Logger.Error("add task to DB error", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": result.Error,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
	})
}

func setStatusComplete(c *gin.Context) {

}

func getAllTasks(c *gin.Context) {
	type Task struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	var completed, uncompleted []Task
	var allTasks []models.Task

	keys, err := config.RedisClient.Keys("*").Result()
	if err != nil {
		config.Logger.Error("get redis keys error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
	}

	if len(keys) != 0 {
		for _, key := range keys {
			task, err := config.RedisClient.HMGet(key, "name", "description", "complete").Result()
			if err != nil {
				config.Logger.Error("get task from redis error", zap.Error(err))
				continue
			}

			id, err := strconv.ParseUint(key, 10, 64)
			if err != nil {
				config.Logger.Error("ParseUint error", zap.Error(err))
				continue
			}

			if task[2] == "1" {
				completed = append(completed, Task{
					ID:          uint(id),
					Name:        task[0].(string),
					Description: task[1].(string),
				})
			} else {
				uncompleted = append(uncompleted, Task{
					ID:          uint(id),
					Name:        task[0].(string),
					Description: task[1].(string),
				})
			}
		}
	} else {
		config.DB.Model(&models.Task{}).Find(&allTasks)

		for _, task := range allTasks {
			if task.Complete {
				completed = append(completed, Task{
					ID:          task.ID,
					Name:        task.Name,
					Description: task.Description,
				})
			} else {
				uncompleted = append(uncompleted, Task{
					ID:          task.ID,
					Name:        task.Name,
					Description: task.Description,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"uncompleted": uncompleted,
		"completed":   completed,
	})
}

func getTask(c *gin.Context) {

}

func deleteTask(c *gin.Context) {

}

func cleanCompleteTasks(c *gin.Context) {

}
