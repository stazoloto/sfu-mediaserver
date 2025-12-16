package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	upgrader websocket.Upgrader

	mu      sync.RWMutex
	clients map[string]*websocket.Conn // clientID -> conn

	controller   Controller
	onDisconnect func(clientID string)
}

// NewHub создаёт WebSocket hub
func NewHub() *Hub {
	return &Hub{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients: make(map[string]*websocket.Conn),
	}
}

// SetController подключает inbound adapter (WSController)
func (h *Hub) SetController(c Controller) {
	h.controller = c
}

// ServeHTTP — entrypoint для net/http
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// client_id передаётся как query param
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		_ = conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(4001, "client_id is required"),
		)
		return
	}

	h.register(clientID, conn)
	defer h.unregister(clientID)

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("ws closed:", clientID)
			} else {
				log.Println("wa read error:", err)
			}

			if h.onDisconnect != nil {
				h.onDisconnect(clientID)
			}
			return
		}

		if h.controller != nil {
			if err := h.controller.Handle(data); err != nil {
				log.Println("controller error:", err)
			}
		}
	}
}

func (h *Hub) register(clientID string, conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[clientID] = conn
	log.Println("WS CONNECT:", clientID)
	h.mu.Unlock()
}

func (h *Hub) unregister(clientID string) {
	h.mu.Lock()
	delete(h.clients, clientID)
	log.Println("WS DISCONNECT:", clientID)
	h.mu.Unlock()
}

// Send отправляет raw JSON конкретному клиенту
func (h *Hub) Send(clientID string, data []byte) error {
	h.mu.RLock()
	conn := h.clients[clientID]
	h.mu.RUnlock()

	if conn == nil {
		return nil
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

func (h *Hub) SetOnDisconnect(fn func(string)) {
	h.onDisconnect = fn
}
