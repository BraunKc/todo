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
			"err": "invalid request",
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

	config.RedisClient.HMSet(strconv.FormatUint(uint64(newTask.ID), 10),
		map[string]interface{}{
			"name":        newTask.Name,
			"description": newTask.Description,
			"complete":    newTask.Complete,
		})

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
	})
}

func setStatusComplete(c *gin.Context) {
	id := c.Param("id")

	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "task not found",
		})
		return
	}

	if task.Complete {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "task alredy completed",
		})
		return
	}

	config.DB.Model(&task).Update("complete", true)
	if err := config.RedisClient.HSet(id, "complete", true).Err(); err != nil {
		config.Logger.Error("redis update task error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
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

			var description = task[1].(string)
			if len(description) > 10 {
				description = task[1].(string)[:10] + "..."
			}
			if task[2] == "1" {
				completed = append(completed, Task{
					ID:          uint(id),
					Name:        task[0].(string),
					Description: description,
				})
			} else {
				uncompleted = append(uncompleted, Task{
					ID:          uint(id),
					Name:        task[0].(string),
					Description: description,
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
					Description: task.Description[:10] + "...",
				})
			} else {
				uncompleted = append(uncompleted, Task{
					ID:          task.ID,
					Name:        task.Name,
					Description: task.Description[:10] + "...",
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
	id := c.Param("id")
	task, err := config.RedisClient.HMGet(id, "name", "description", "complete").Result()
	if err != nil {
		config.Logger.Error("redis HMGet error", zap.Error(err))
	}

	if task[0] == nil {
		var dbTask models.Task
		if err := config.DB.First(&dbTask, id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "task not found",
			})
			return
		}
		task = []interface{}{dbTask.Name, dbTask.Description, bool(dbTask.Complete)}

		config.RedisClient.HMSet(id, map[string]interface{}{
			"name":        task[0].(string),
			"description": task[1].(string),
			"complete":    task[2].(bool),
		})
	}

	respID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		config.Logger.Error("ParseUint error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	var complete bool
	if task[2] != nil {
		switch task[2].(type) {
		case string:
			complete = task[2].(string) == "1"
		case bool:
			complete = task[2].(bool)
		}
	}

	responseTask := struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Complete    bool   `json:"complete"`
	}{
		ID:          uint(respID),
		Name:        task[0].(string),
		Description: task[1].(string),
		Complete:    complete,
	}

	c.JSON(http.StatusOK, gin.H{
		"task": responseTask,
	})
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")
	config.RedisClient.Del(id)
	if err := config.DB.Delete(&models.Task{}, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func cleanCompleteTasks(c *gin.Context) {
	keys, err := config.RedisClient.Keys("*").Result()
	if err != nil {
		config.Logger.Error("get redis keys error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
	}

	if len(keys) != 0 {
		for _, key := range keys {
			task, err := config.RedisClient.HGet(key, "complete").Result()
			if err != nil {
				config.Logger.Error("get task from redis error", zap.Error(err))
				continue
			}

			if task == "1" {
				config.RedisClient.Del(key)
			}
		}
	}

	var allTasks []models.Task
	config.DB.Model(&models.Task{}).Find(&allTasks)

	for _, task := range allTasks {
		if task.Complete {
			if err := config.DB.Delete(&models.Task{}, task.ID).Error; err != nil {
				config.Logger.Error("delete task from DB error", zap.Error(err))
				continue
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
