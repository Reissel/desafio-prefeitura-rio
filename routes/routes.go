package routes

import (
	"desafio-prefeitura-rio/logic"
	"desafio-prefeitura-rio/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	r.POST("/", middleware.SignatureMiddleware(), logic.CreateNotification)

	v1 := r.Group("/notifications")
	v1.Use(middleware.AuthMiddleware())
	{
		v1.GET("/", logic.GetNotifications)
	}
}
