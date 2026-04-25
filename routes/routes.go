package routes

import (
	"desafio-prefeitura-rio/database"
	"desafio-prefeitura-rio/logic"
	"desafio-prefeitura-rio/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Print("Database connected!")

	notificationRepo := &database.PostgresNotificationRepo{DB: db}

	dlq := initializeDLQ()

	handler := &logic.NotificationHandler{
		DB:    notificationRepo,
		Store: dlq,
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

func initializeDLQ() logic.DLQStore {
	var store logic.DLQStore
	storageType := os.Getenv("DLQ_TYPE")

	switch storageType {
	case "redis":
		rdb := database.ConnectRedis()
		store = &logic.RedisDLQRepo{Client: rdb, Key: "notifications_dlq"}
	case "postgres":
		db, _ := database.ConnectDB()
		store = &logic.PostgresDLQRepo{DB: db}
	default:
		log.Fatal("Unknown DLQ type")
	}

	return store
}
