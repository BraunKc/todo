package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/BraunKc/todo/api-service/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func add(c *gin.Context) {
	var task struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		config.Logger.Error("POST /add error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	jsonTask, err := json.Marshal(&task)
	if err != nil {
		config.Logger.Error("task marshal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	url := fmt.Sprintf("%s/create-task", config.Config.DBService.Endpoint)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonTask))
	if err != nil {
		config.Logger.Error("http.Post error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		config.Logger.Error("io.ReadAll error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	var respBody map[string]any
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		config.Logger.Error("body unmarshal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(resp.StatusCode, respBody)
}

func complete(c *gin.Context) {
	id := c.Param("id")

	url := fmt.Sprintf("%s/set-status-complete/%s", config.Config.DBService.Endpoint, id)
	resp, err := http.Get(url)
	if err != nil {
		config.Logger.Error("http.Get error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		config.Logger.Error("io.ReadAll error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	var respBody map[string]any
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		config.Logger.Error("body unmarshal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(resp.StatusCode, respBody)
}

func tasks(c *gin.Context) {
	url := fmt.Sprintf("%s/get-all-tasks", config.Config.DBService.Endpoint)
	resp, err := http.Get(url)
	if err != nil {
		config.Logger.Error("http.Get error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		config.Logger.Error("io.ReadAll error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	type Task struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	var respBody map[string][]Task
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		config.Logger.Error("body unmarshal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(resp.StatusCode, respBody)
}

func task(c *gin.Context) {
	id := c.Param("id")

	url := fmt.Sprintf("%s/get-task/%s", config.Config.DBService.Endpoint, id)
	resp, err := http.Get(url)
	if err != nil {
		config.Logger.Error("http.Get error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		config.Logger.Error("io.ReadAll error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	type Task struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Complete    bool   `json:"complete"`
	}
	var respBody map[string]Task
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		config.Logger.Error("body unmarshal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(resp.StatusCode, respBody)
}

func delete(c *gin.Context) {
	id := c.Param("id")

	url := fmt.Sprintf("%s/delete-task/%s", config.Config.DBService.Endpoint, id)
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		config.Logger.Error("http.NewRequest error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		config.Logger.Error("client.Do error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		config.Logger.Error("io.ReadAll error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	var respBody map[string]any
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		config.Logger.Error("body unmarshal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(resp.StatusCode, respBody)
}

func clean(c *gin.Context) {
	url := fmt.Sprintf("%s/clean-completed-tasks", config.Config.DBService.Endpoint)
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		config.Logger.Error("http.NewRequest error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		config.Logger.Error("client.Do error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		config.Logger.Error("io.ReadAll error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	var respBody map[string]any
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		config.Logger.Error("body unmarshal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(resp.StatusCode, respBody)
}
