package logic

import (
	"desafio-prefeitura-rio/database"
	"desafio-prefeitura-rio/model"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	cpfQuery := c.MustGet("cpf").(string)

	searchHash := GenerateBlindIndex(cpfQuery)

	notifications, err := database.GetNotifications(searchHash)
	if err != nil {
		log.Fatalf("Error fetching notifications: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"data": notifications})
}

func CreateNotification(c *gin.Context) {
	var notification model.Notification

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input. All fields are required."})
		return
	}

	blindIndex := GenerateBlindIndex(notification.Cpf)
	encryptedBlob, err := Encrypt(notification.Cpf)
	if err != nil {
		c.JSON(500, gin.H{"error": "Encryption failed"})
		return
	}

	query := `
        INSERT INTO notifications (chamado_id, tipo, cpf_encrypted, cpf_blind_index, status_anterior, status_novo, titulo, descricao, timestamp) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
        RETURNING id`

	db, err := database.ConnectDB()

	err = db.QueryRow(query, notification.Chamado_id, notification.Tipo, encryptedBlob, blindIndex,
		notification.Status_anterior, notification.Status_novo, notification.Titulo, notification.Descricao, notification.Timestamp).Scan(&notification.ID)

	if err != nil {
		// TODO: Criar lógica para impedir criar para o mesmo registro
		//if strings.Contains(err.Error(), "unique constraint") {
		//  c.JSON(409, gin.H{"error": "A notification with this CPF already exists"})
		//   return
		//}

		fmt.Print(err)

		c.JSON(500, gin.H{"error": "Failed to save to database"})
		return
	}

	c.JSON(201, gin.H{
		"message": "Notification created successfully",
		"data":    notification.ID,
	})
}

func SetIsReadNotification(c *gin.Context) {
	id := c.Param("id")

	cpfQuery := c.MustGet("cpf").(string)

	searchHash := GenerateBlindIndex(cpfQuery)

	updatedId, err := database.SetIsReadNotification(id, searchHash)

	if err != nil {
		fmt.Print(err)

		c.JSON(500, gin.H{"error": "Failed to save to database"})
		return
	}

	c.JSON(201, gin.H{
		"message": "Notification updated successfully",
		"data":    updatedId,
	})
}

func CountUnreadNotifications(c *gin.Context) {
	var count int

	cpfQuery := c.MustGet("cpf").(string)

	searchHash := GenerateBlindIndex(cpfQuery)

	count, err := database.CountUnreadNotifications(searchHash)

	if err != nil {
		fmt.Print(err)

		c.JSON(500, gin.H{"error": "Failed to read from database"})
		return
	}

	c.JSON(201, gin.H{
		"message": "Notification count read successfully",
		"data":    count,
	})
}
