package api

import (
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		config.Logger.Error("POST /add error:", zap.Error(err))
		return
	}
}

func complete(c *gin.Context) {

}

func tasks(c *gin.Context) {

}

func task(c *gin.Context) {

}

func delete(c *gin.Context) {

}

func clean(c *gin.Context) {

}
