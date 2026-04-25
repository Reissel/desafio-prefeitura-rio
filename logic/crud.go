package logic

import (
	"context"
	"desafio-prefeitura-rio/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationProvider interface {
	GetNotifications(cpfEncrypted string) ([]model.Notification, error)
	CreateNotification(notification model.Notification, blindIndex string, encryptedBlob []byte) (int, error)
	SetIsReadNotification(id string, cpfEncrypted string) (string, error)
	CountUnreadNotifications(cpfEncrypted string) (int, error)
}

type NotificationHandler struct {
	DB    NotificationProvider
	Store DLQStore
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	cpfQuery, exists := c.Get("cpf")
	if !exists {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	searchHash := GenerateBlindIndex(cpfQuery.(string))

	notifications, err := h.DB.GetNotifications(searchHash)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notifications})
}

func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var notification model.Notification

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input. All fields are required."})
		return
	}

	blindIndex := GenerateBlindIndex(notification.Cpf)
	encryptedBlob, err := Encrypt(notification.Cpf)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption failed: " + err.Error()})
		return
	}

	createdId, err := h.DB.CreateNotification(notification, blindIndex, encryptedBlob)
	if err != nil {
		log.Print(err)

		ctx := context.Background()

		dlqId, dlqErr := h.Store.Save(ctx, notification)
		if dlqErr != nil {
			log.Print(dlqErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to dlq: " + dlqErr.Error()})
			return
		}

		log.Printf("Saved request to dlq with id %d", dlqId)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database: " + err.Error()})
		return
	}

	c.Set("response", gin.H{
		"message": "Notification created successfully",
		"data":    createdId,
	})

	c.Set("notification", notification)

	c.Next()
}

func (h *NotificationHandler) SetIsReadNotification(c *gin.Context) {
	id := c.Param("id")

	cpfQuery := c.MustGet("cpf").(string)

	searchHash := GenerateBlindIndex(cpfQuery)

	updatedId, err := h.DB.SetIsReadNotification(id, searchHash)

	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Notification updated successfully",
		"data":    updatedId,
	})
}

func (h *NotificationHandler) CountUnreadNotifications(c *gin.Context) {
	var count int

	cpfQuery := c.MustGet("cpf").(string)

	searchHash := GenerateBlindIndex(cpfQuery)

	count, err := h.DB.CountUnreadNotifications(searchHash)

	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read from database: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Notification count read successfully",
		"data":    count,
	})
}
