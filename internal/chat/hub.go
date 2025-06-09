package chat

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Hub struct {
	clients map[uint]*websocket.Conn
	mu      sync.RWMutex
}

var HubInstance = &Hub{
	clients: make(map[uint]*websocket.Conn),
}
