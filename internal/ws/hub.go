package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// Map of user_id to their websocket connection
	clients map[string]*websocket.Conn
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*websocket.Conn),
	}
}

func (h *Hub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// If there's an existing connection, close it
	if oldConn, ok := h.clients[userID]; ok {
		oldConn.Close()
	}
	h.clients[userID] = conn
}

func (h *Hub) Unregister(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conn, ok := h.clients[userID]; ok {
		conn.Close()
		delete(h.clients, userID)
	}
}

func (h *Hub) SendToUser(userID string, message interface{}) error {
	h.mu.RLock()
	conn, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return nil // User not connected
	}

	return conn.WriteJSON(message)
}
