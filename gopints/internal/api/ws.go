package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

// Hub manages active WebSocket connections and broadcasts messages.
type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func newHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]struct{})}
}

// ServeWS upgrades the HTTP connection and registers the client with the hub.
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Warn("ws upgrade failed", "err", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = struct{}{}
	h.mu.Unlock()

	// Read loop to detect disconnection.
	go func() {
		defer func() {
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			conn.Close()
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}

// Broadcast sends v as JSON to all connected clients.
// Clients that fail to receive are removed.
func (h *Hub) Broadcast(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		slog.Error("ws marshal", "err", err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			delete(h.clients, conn)
			conn.Close()
		}
	}
}
