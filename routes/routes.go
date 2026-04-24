package routes

import (
	"desafio-prefeitura-rio/database"
	"desafio-prefeitura-rio/logic"
	"desafio-prefeitura-rio/middleware"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	db, err := database.ConnectDB()
	if err != nil {
		fmt.Print(err)
	}

	notificationRepo := &database.PostgresNotificationRepo{DB: db}

	handler := &logic.NotificationHandler{
		DB: notificationRepo,
	}

	r.POST("/", middleware.SignatureMiddleware(), middleware.IdempotencyMiddleware(), handler.CreateNotification, logic.UpdateClient)

	r.GET("/ws/:cpf", logic.SetupClient)

	v1 := r.Group("/notifications")
	v1.Use(middleware.AuthMiddleware())
	{
		v1.GET("/", handler.GetNotifications)
		v1.PATCH("/:id/read", handler.SetIsReadNotification)
		v1.GET("/unread-count", handler.CountUnreadNotifications)
	}
}
