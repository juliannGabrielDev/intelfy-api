package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/juliannGabrielDev/intelfy-api/internal/middleware"
	"github.com/juliannGabrielDev/intelfy-api/internal/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust for production
	},
}

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

func (h *WSHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	h.hub.Register(userID, conn)

	// Keep connection alive and wait for close
	defer h.hub.Unregister(userID)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
