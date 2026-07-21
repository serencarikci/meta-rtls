package location

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]string
}

func NewHub() *Hub {
	return &Hub{clients: map[*websocket.Conn]string{}}
}

func (h *Hub) Add(conn *websocket.Conn, tenantID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = tenantID
}

func (h *Hub) Remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
	_ = conn.Close()
}

func (h *Hub) Broadcast(pos LivePosition) {
	payload, err := json.Marshal(pos)
	if err != nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	for conn, tenantID := range h.clients {
		if tenantID != pos.TenantID {
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			_ = conn.Close()
			delete(h.clients, conn)
		}
	}
}
