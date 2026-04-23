package logic

import (
	"desafio-prefeitura-rio/model"
	"desafio-prefeitura-rio/websocket"

	"github.com/gin-gonic/gin"
)

func SetupClient(c *gin.Context) {
	cpf := c.Param("cpf")
	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	websocket.WebHub.Lock()
	websocket.WebHub.Clients[cpf] = conn
	websocket.WebHub.Unlock()

	// Keep connection alive/clean up on disconnect
	defer func() {
		websocket.WebHub.Lock()
		delete(websocket.WebHub.Clients, cpf)
		websocket.WebHub.Unlock()
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func UpdateClient(c *gin.Context) {
	notification := c.MustGet("notification").(model.Notification)

	var payload gin.H

	val := c.MustGet("response")
	payload = val.(gin.H)

	websocket.WebHub.RLock()
	conn, exists := websocket.WebHub.Clients[notification.Cpf]
	websocket.WebHub.RUnlock()

	if exists {
		conn.WriteJSON(gin.H{"payload": notification})
		payload["status"] = "user notified"
		c.JSON(201, payload)
	} else {
		payload["status"] = "user not connected"
		c.JSON(201, payload)
	}
}
