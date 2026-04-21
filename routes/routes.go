package routes

import (
	"desafio-prefeitura-rio/logic"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	v1 := r.Group("/notifications")
	{
		v1.GET("/", logic.GetNotifications)

		v1.POST("/", logic.CreateNotification)
	}
}
