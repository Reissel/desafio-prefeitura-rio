package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub manages active connections
type Hub struct {
	sync.RWMutex
	Clients map[string]*websocket.Conn
}

var WebHub = Hub{
	Clients: make(map[string]*websocket.Conn),
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
